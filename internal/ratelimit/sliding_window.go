package ratelimit

import (
	"sync"
	"time"
)

type SlidingWindow struct {
	windowSize time.Duration
	limit      int
	requests   []time.Time
	mu         sync.Mutex
}

func NewSlidingWindow(windowSize time.Duration, limit int) *SlidingWindow {
	return &SlidingWindow{
		windowSize: windowSize,
		limit:      limit,
		requests:   make([]time.Time, 0),
	}
}

func (sw *SlidingWindow) Allow() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-sw.windowSize)

	newRequests := make([]time.Time, 0, len(sw.requests))
	for _, req := range sw.requests {
		if req.After(windowStart) {
			newRequests = append(newRequests, req)
		}
	}
	sw.requests = newRequests

	if len(sw.requests) < sw.limit {
		sw.requests = append(sw.requests, now)
		return true
	}

	return false
}

func (sw *SlidingWindow) Reset() {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.requests = make([]time.Time, 0)
}

func (sw *SlidingWindow) Count() int {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-sw.windowSize)

	count := 0
	for _, req := range sw.requests {
		if req.After(windowStart) {
			count++
		}
	}
	return count
}