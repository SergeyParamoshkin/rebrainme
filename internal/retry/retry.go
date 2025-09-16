package retry

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

var (
	ErrMaxRetriesExceeded = errors.New("maximum retries exceeded")
	ErrContextCancelled   = errors.New("context cancelled")
)

type Config struct {
	MaxRetries     int
	InitialDelay   time.Duration
	MaxDelay       time.Duration
	Multiplier     float64
	Jitter         bool
	RetryableError func(error) bool
}

func DefaultConfig() *Config {
	return &Config{
		MaxRetries:   3,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     10 * time.Second,
		Multiplier:   2.0,
		Jitter:       true,
		RetryableError: func(err error) bool {
			return true
		},
	}
}

type Retrier struct {
	config *Config
	budget *Budget
}

func New(config *Config) *Retrier {
	if config == nil {
		config = DefaultConfig()
	}
	return &Retrier{
		config: config,
	}
}

func NewWithBudget(config *Config, budget *Budget) *Retrier {
	if config == nil {
		config = DefaultConfig()
	}
	return &Retrier{
		config: config,
		budget: budget,
	}
}

func (r *Retrier) Do(ctx context.Context, fn func() error) error {
	_, err := r.DoWithData(ctx, func() (interface{}, error) {
		return nil, fn()
	})
	return err
}

func (r *Retrier) DoWithData(ctx context.Context, fn func() (interface{}, error)) (interface{}, error) {
	var lastErr error
	delay := r.config.InitialDelay

	for attempt := 0; attempt <= r.config.MaxRetries; attempt++ {
		if attempt > 0 {
			if r.budget != nil && !r.budget.CanRetry() {
				return nil, fmt.Errorf("retry budget exhausted: %w", lastErr)
			}

			select {
			case <-ctx.Done():
				return nil, fmt.Errorf("%w: %w", ErrContextCancelled, ctx.Err())
			case <-time.After(delay):
			}

			delay = r.calculateNextDelay(delay)
		}

		result, err := fn()
		if err == nil {
			if r.budget != nil && attempt > 0 {
				r.budget.RecordSuccess()
			}
			return result, nil
		}

		lastErr = err

		if !r.config.RetryableError(err) {
			return nil, err
		}

		if r.budget != nil && attempt > 0 {
			r.budget.RecordFailure()
		}

		if attempt == r.config.MaxRetries {
			return nil, fmt.Errorf("%w: %w", ErrMaxRetriesExceeded, lastErr)
		}
	}

	return nil, lastErr
}

func (r *Retrier) calculateNextDelay(currentDelay time.Duration) time.Duration {
	nextDelay := time.Duration(float64(currentDelay) * r.config.Multiplier)

	if nextDelay > r.config.MaxDelay {
		nextDelay = r.config.MaxDelay
	}

	if r.config.Jitter {
		jitter := rand.Float64() * 0.3 * float64(nextDelay)
		nextDelay = time.Duration(float64(nextDelay) + jitter - jitter/2)
	}

	return nextDelay
}

func ExponentialBackoff(attempt int, baseDelay time.Duration) time.Duration {
	return time.Duration(math.Pow(2, float64(attempt))) * baseDelay
}

func LinearBackoff(attempt int, baseDelay time.Duration) time.Duration {
	return time.Duration(attempt+1) * baseDelay
}

func ConstantBackoff(baseDelay time.Duration) time.Duration {
	return baseDelay
}