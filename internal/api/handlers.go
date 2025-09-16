package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/SergeyParamoshkin/rebrainme/internal/client"
	"github.com/SergeyParamoshkin/rebrainme/internal/logger"
	"github.com/SergeyParamoshkin/rebrainme/internal/monitoring"
	"github.com/SergeyParamoshkin/rebrainme/internal/retry"
	"github.com/go-chi/chi/v5"
)

func RegisterHandlers(router chi.Router) {
	router.Route("/api", func(r chi.Router) {
		r.Get("/simple", SimpleHandler)
		r.Get("/flaky", FlakyHandler)
		r.Get("/slow", SlowHandler)
		r.Get("/retry-demo", RetryDemoHandler)
		r.Get("/context-demo", ContextDemoHandler)
		r.Get("/circuit-breaker", CircuitBreakerHandler)
		r.Get("/budget-demo", BudgetDemoHandler)
	})
	router.Get("/health", HealthHandler)
}

func SimpleHandler(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Simple handler called", "path", r.URL.Path)

	response := map[string]interface{}{
		"message":   "Simple response",
		"timestamp": time.Now().Unix(),
		"path":      r.URL.Path,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func FlakyHandler(w http.ResponseWriter, r *http.Request) {
	if rand.Float32() < 0.5 {
		logger.Warn("Flaky handler simulating failure")
		http.Error(w, "Random failure", http.StatusServiceUnavailable)
		return
	}

	response := map[string]interface{}{
		"message": "Success after potential failure",
		"lucky":   true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func SlowHandler(w http.ResponseWriter, r *http.Request) {
	delay := time.Duration(rand.Intn(3000)) * time.Millisecond
	time.Sleep(delay)

	response := map[string]interface{}{
		"message":  "Slow response",
		"delay_ms": delay.Milliseconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func RetryDemoHandler(w http.ResponseWriter, r *http.Request) {
	config := &retry.Config{
		MaxRetries:   5,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     5 * time.Second,
		Multiplier:   2.0,
		Jitter:       true,
		RetryableError: func(err error) bool {
			return true
		},
	}

	retrier := retry.New(config)

	attempts := 0
	err := retrier.Do(r.Context(), func() error {
		attempts++
		logger.Debug("Retry attempt", "attempt", attempts)
		if attempts < 3 {
			monitoring.RetryAttempts.WithLabelValues("demo", "failure").Inc()
			return fmt.Errorf("simulated failure %d", attempts)
		}
		monitoring.RetryAttempts.WithLabelValues("demo", "success").Inc()
		return nil
	})

	if err != nil {
		logger.Error("Retry failed after all attempts", err, "attempts", attempts)
		http.Error(w, fmt.Sprintf("Retry failed: %v", err), http.StatusServiceUnavailable)
		return
	}

	response := map[string]interface{}{
		"message":  "Success after retries",
		"attempts": attempts,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func ContextDemoHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	resultChan := make(chan string)
	errorChan := make(chan error)

	go func() {
		delay := time.Duration(rand.Intn(4000)) * time.Millisecond
		select {
		case <-time.After(delay):
			resultChan <- fmt.Sprintf("Operation completed after %v", delay)
		case <-ctx.Done():
			errorChan <- ctx.Err()
		}
	}()

	select {
	case result := <-resultChan:
		response := map[string]interface{}{
			"message": result,
			"status":  "completed",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	case err := <-errorChan:
		http.Error(w, fmt.Sprintf("Context error: %v", err), http.StatusRequestTimeout)

	case <-ctx.Done():
		http.Error(w, "Request timeout", http.StatusRequestTimeout)
	}
}

func CircuitBreakerHandler(w http.ResponseWriter, r *http.Request) {
	cb := client.GetCircuitBreaker()

	result, err := cb.Execute(func() (interface{}, error) {
		if rand.Float32() < 0.3 {
			return nil, fmt.Errorf("simulated failure")
		}
		return "Success", nil
	})

	state := "closed"
	stateValue := 0.0
	if cb.IsOpen() {
		state = "open"
		stateValue = 1.0
	}

	monitoring.CircuitBreakerState.WithLabelValues("demo").Set(stateValue)

	if err != nil {
		response := map[string]interface{}{
			"error":  err.Error(),
			"state":  state,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"message": result,
		"state":   state,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func BudgetDemoHandler(w http.ResponseWriter, r *http.Request) {
	budget := retry.NewBudget(0.5, 10*time.Second, 10)

	config := &retry.Config{
		MaxRetries:   3,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     2 * time.Second,
		Multiplier:   2.0,
		Jitter:       true,
		RetryableError: func(err error) bool {
			return true
		},
	}

	retrier := retry.NewWithBudget(config, budget)

	err := retrier.Do(r.Context(), func() error {
		if rand.Float32() < 0.7 {
			return fmt.Errorf("simulated failure")
		}
		return nil
	})

	successRate, totalRequests := budget.GetStats()

	response := map[string]interface{}{
		"success":        err == nil,
		"budget_stats": map[string]interface{}{
			"success_rate":    successRate,
			"total_requests": totalRequests,
		},
	}

	if err != nil {
		monitoring.RetryBudgetExhausted.WithLabelValues("demo").Inc()
		response["error"] = err.Error()
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status": "healthy",
		"timestamp": time.Now().Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func ExternalAPICall(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

