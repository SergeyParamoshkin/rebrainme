package ratelimit

import (
	"testing"
	"time"
)

func TestTokenBucket(t *testing.T) {
	tb := NewTokenBucket(10, 10)

	for i := 0; i < 10; i++ {
		if !tb.Allow() {
			t.Errorf("Expected request %d to be allowed", i)
		}
	}

	if tb.Allow() {
		t.Error("Expected 11th request to be denied")
	}

	time.Sleep(100 * time.Millisecond)

	if !tb.Allow() {
		t.Error("Expected request after refill to be allowed")
	}
}

func TestSlidingWindow(t *testing.T) {
	sw := NewSlidingWindow(100*time.Millisecond, 5)

	for i := 0; i < 5; i++ {
		if !sw.Allow() {
			t.Errorf("Expected request %d to be allowed", i)
		}
	}

	if sw.Allow() {
		t.Error("Expected 6th request to be denied")
	}

	time.Sleep(110 * time.Millisecond)

	if !sw.Allow() {
		t.Error("Expected request after window slide to be allowed")
	}
}

func BenchmarkTokenBucket(b *testing.B) {
	tb := NewTokenBucket(1000, 1000)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tb.Allow()
		}
	})
}

func BenchmarkSlidingWindow(b *testing.B) {
	sw := NewSlidingWindow(time.Second, 1000)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sw.Allow()
		}
	})
}