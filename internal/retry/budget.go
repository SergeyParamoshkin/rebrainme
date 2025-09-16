package retry

import (
	"sync"
	"time"
)

type Budget struct {
	successThreshold float64
	windowSize       time.Duration
	minRequests      int

	requests []request
	mu       sync.RWMutex
}

type request struct {
	timestamp time.Time
	success   bool
}

func NewBudget(successThreshold float64, windowSize time.Duration, minRequests int) *Budget {
	return &Budget{
		successThreshold: successThreshold,
		windowSize:       windowSize,
		minRequests:      minRequests,
		requests:         make([]request, 0),
	}
}

func (b *Budget) CanRetry() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	b.cleanOldRequests()

	if len(b.requests) < b.minRequests {
		return true
	}

	successCount := 0
	for _, req := range b.requests {
		if req.success {
			successCount++
		}
	}

	successRate := float64(successCount) / float64(len(b.requests))
	return successRate >= b.successThreshold
}

func (b *Budget) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.requests = append(b.requests, request{
		timestamp: time.Now(),
		success:   true,
	})
	b.cleanOldRequests()
}

func (b *Budget) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.requests = append(b.requests, request{
		timestamp: time.Now(),
		success:   false,
	})
	b.cleanOldRequests()
}

func (b *Budget) cleanOldRequests() {
	now := time.Now()
	windowStart := now.Add(-b.windowSize)

	newRequests := make([]request, 0, len(b.requests))
	for _, req := range b.requests {
		if req.timestamp.After(windowStart) {
			newRequests = append(newRequests, req)
		}
	}
	b.requests = newRequests
}

func (b *Budget) GetStats() (successRate float64, totalRequests int) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	b.cleanOldRequests()

	if len(b.requests) == 0 {
		return 1.0, 0
	}

	successCount := 0
	for _, req := range b.requests {
		if req.success {
			successCount++
		}
	}

	return float64(successCount) / float64(len(b.requests)), len(b.requests)
}