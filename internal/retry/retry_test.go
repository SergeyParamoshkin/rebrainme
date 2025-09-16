package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetrySuccess(t *testing.T) {
	config := &Config{
		MaxRetries:   3,
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     100 * time.Millisecond,
		Multiplier:   2.0,
		Jitter:       false,
		RetryableError: func(err error) bool {
			return true
		},
	}

	retrier := New(config)

	attempts := 0
	err := retrier.Do(context.Background(), func() error {
		attempts++
		if attempts < 3 {
			return errors.New("temporary error")
		}
		return nil
	})

	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}

	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestRetryContextTimeout(t *testing.T) {
	config := &Config{
		MaxRetries:   10,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     1 * time.Second,
		Multiplier:   2.0,
		Jitter:       false,
		RetryableError: func(err error) bool {
			return true
		},
	}

	retrier := New(config)

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	err := retrier.Do(ctx, func() error {
		return errors.New("always fails")
	})

	if err == nil {
		t.Error("Expected context timeout error")
	}

	if !errors.Is(err, ErrContextCancelled) {
		t.Errorf("Expected ErrContextCancelled, got %v", err)
	}
}

func TestRetryBudget(t *testing.T) {
	budget := NewBudget(0.5, 1*time.Second, 2)

	budget.RecordSuccess()
	budget.RecordSuccess()

	if !budget.CanRetry() {
		t.Error("Expected retry to be allowed with 100% success rate")
	}

	budget.RecordFailure()
	budget.RecordFailure()
	budget.RecordFailure()

	if budget.CanRetry() {
		t.Error("Expected retry to be denied with less than 50% success rate")
	}
}

func TestCircuitBreaker(t *testing.T) {
	cb := NewCircuitBreaker(2, 100*time.Millisecond)

	_, err := cb.Execute(func() (interface{}, error) {
		return nil, errors.New("error")
	})
	if err == nil {
		t.Error("Expected error")
	}

	_, err = cb.Execute(func() (interface{}, error) {
		return nil, errors.New("error")
	})
	if err == nil {
		t.Error("Expected error")
	}

	if !cb.IsOpen() {
		t.Error("Expected circuit to be open after 2 failures")
	}

	_, err = cb.Execute(func() (interface{}, error) {
		return "success", nil
	})
	if !errors.Is(err, ErrCircuitOpen) {
		t.Error("Expected ErrCircuitOpen")
	}

	time.Sleep(110 * time.Millisecond)

	result, err := cb.Execute(func() (interface{}, error) {
		return "success", nil
	})
	if err != nil {
		t.Errorf("Expected success in half-open state, got %v", err)
	}

	if result != "success" {
		t.Errorf("Expected 'success', got %v", result)
	}

	if cb.IsOpen() {
		t.Error("Expected circuit to be closed after success")
	}
}

func BenchmarkRetry(b *testing.B) {
	config := DefaultConfig()
	retrier := New(config)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			retrier.Do(context.Background(), func() error {
				return nil
			})
		}
	})
}