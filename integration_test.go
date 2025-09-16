package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/SergeyParamoshkin/rebrainme/internal/api"
	"github.com/SergeyParamoshkin/rebrainme/internal/logger"
	customMiddleware "github.com/SergeyParamoshkin/rebrainme/internal/middleware"
	"github.com/SergeyParamoshkin/rebrainme/internal/monitoring"
	"github.com/go-chi/chi/v5"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
	colorBold   = "\033[1m"
)

type TestResult struct {
	RequestNum   int
	StatusCode   int
	Duration     time.Duration
	Error        error
	ResponseBody string
}

func printTestHeader(title string) {
	fmt.Printf("\n%s%s╔══════════════════════════════════════════════════════════════╗%s\n", colorBold, colorCyan, colorReset)
	fmt.Printf("%s%s║  %-60s║%s\n", colorBold, colorCyan, title, colorReset)
	fmt.Printf("%s%s╚══════════════════════════════════════════════════════════════╝%s\n\n", colorBold, colorCyan, colorReset)
}

func printRequestResult(result TestResult) {
	statusColor := colorGreen
	statusIcon := "✓"

	switch {
	case result.StatusCode == 429:
		statusColor = colorYellow
		statusIcon = "⚠"
	case result.StatusCode >= 500:
		statusColor = colorRed
		statusIcon = "✗"
	case result.StatusCode >= 400:
		statusColor = colorYellow
		statusIcon = "⚠"
	}

	fmt.Printf("  %sRequest #%02d:%s %s%s %d%s (%.2fms)",
		colorBold, result.RequestNum, colorReset,
		statusColor, statusIcon, result.StatusCode, colorReset,
		float64(result.Duration.Microseconds())/1000)

	if result.ResponseBody != "" && result.StatusCode != 200 {
		fmt.Printf(" - %s%s%s", colorRed, result.ResponseBody, colorReset)
	}
	fmt.Println()
}

func printSummary(results []TestResult, testName string) {
	var success, rateLimited, serverError, other int
	var totalDuration time.Duration

	for _, r := range results {
		totalDuration += r.Duration
		switch r.StatusCode {
		case 200:
			success++
		case 429:
			rateLimited++
		case 503:
			serverError++
		default:
			other++
		}
	}

	avgDuration := totalDuration / time.Duration(len(results))

	fmt.Printf("\n%s📊 Summary for %s:%s\n", colorBold, testName, colorReset)
	fmt.Printf("  %s✓ Success:%s %s%d%s requests\n", colorGreen, colorReset, colorBold, success, colorReset)

	if rateLimited > 0 {
		fmt.Printf("  %s⚠ Rate Limited (429):%s %s%d%s requests\n", colorYellow, colorReset, colorBold, rateLimited, colorReset)
	}

	if serverError > 0 {
		fmt.Printf("  %s✗ Service Unavailable (503):%s %s%d%s requests\n", colorRed, colorReset, colorBold, serverError, colorReset)
	}

	if other > 0 {
		fmt.Printf("  %s? Other:%s %s%d%s requests\n", colorPurple, colorReset, colorBold, other, colorReset)
	}

	fmt.Printf("  %s⏱ Average Duration:%s %s%.2fms%s\n", colorBlue, colorReset, colorBold, float64(avgDuration.Microseconds())/1000, colorReset)
	fmt.Printf("  %s📈 Success Rate:%s %s%.1f%%%s\n", colorCyan, colorReset, colorBold, float64(success)/float64(len(results))*100, colorReset)
}

var testServerInitialized bool
var testServerMutex sync.Mutex

func init() {
	// Инициализируем logger для тестов
	logger.Init(false)
}

func setupTestServer() *httptest.Server {
	testServerMutex.Lock()
	if !testServerInitialized {
		monitoring.InitMetrics()
		testServerInitialized = true
	}
	testServerMutex.Unlock()

	router := chi.NewRouter()

	// Rate limiter с маленькими лимитами для демонстрации
	rateLimiter := customMiddleware.NewRateLimiter(5, 20) // 5 req/sec, 20 req/min
	router.Use(customMiddleware.RateLimitMiddleware(rateLimiter))
	router.Use(customMiddleware.MetricsMiddleware)

	api.RegisterHandlers(router)

	return httptest.NewServer(router)
}

func TestRateLimiting_429_Demo(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	printTestHeader("🚦 Rate Limiting Test - Demonstrating 429 Errors")

	const numRequests = 15
	results := make([]TestResult, numRequests)

	fmt.Println("  📤 Sending 15 rapid requests (limit: 5 req/sec)...")
	fmt.Println()

	// Отправляем запросы последовательно быстро
	for i := 0; i < numRequests; i++ {
		start := time.Now()

		resp, err := http.Get(server.URL + "/api/simple")
		duration := time.Since(start)

		result := TestResult{
			RequestNum: i + 1,
			Duration:   duration,
			Error:      err,
		}

		if err == nil {
			result.StatusCode = resp.StatusCode
			if resp.StatusCode == 429 {
				body, _ := io.ReadAll(resp.Body)
				result.ResponseBody = string(body)
			}
			resp.Body.Close()
		}

		results[i] = result
		printRequestResult(result)

		// Небольшая задержка между запросами
		if i < numRequests-1 {
			time.Sleep(10 * time.Millisecond)
		}
	}

	printSummary(results, "Rate Limiting")
}

func TestRateLimiting_Parallel_Demo(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	printTestHeader("🔥 Parallel Rate Limiting Test - Burst Traffic")

	const numRequests = 20
	results := make([]TestResult, numRequests)
	var wg sync.WaitGroup
	var mu sync.Mutex

	fmt.Printf("  🚀 Launching %d parallel requests...\n\n", numRequests)

	startTime := time.Now()

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			start := time.Now()
			resp, err := http.Get(server.URL + "/api/simple")
			duration := time.Since(start)

			result := TestResult{
				RequestNum: idx + 1,
				Duration:   duration,
				Error:      err,
			}

			if err == nil {
				result.StatusCode = resp.StatusCode
				if resp.StatusCode == 429 {
					body, _ := io.ReadAll(resp.Body)
					result.ResponseBody = string(body)
				}
				resp.Body.Close()
			}

			mu.Lock()
			results[idx] = result
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	totalDuration := time.Since(startTime)

	// Сортируем результаты по номеру запроса для красивого вывода
	for _, result := range results {
		printRequestResult(result)
	}

	printSummary(results, "Parallel Burst")
	fmt.Printf("  %s⏱ Total Time:%s %s%.2fms%s\n", colorPurple, colorReset, colorBold, float64(totalDuration.Microseconds())/1000, colorReset)
}

func TestRetryAndCircuitBreaker_503_Demo(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	printTestHeader("🔄 Circuit Breaker Test - Demonstrating 503 Errors")

	const numRequests = 15
	results := make([]TestResult, numRequests)

	fmt.Println("  🔮 Testing Circuit Breaker with flaky endpoint...")
	fmt.Println("  📊 Circuit opens after 5 failures, resets after 10s")
	fmt.Println()

	for i := 0; i < numRequests; i++ {
		start := time.Now()

		resp, err := http.Get(server.URL + "/api/circuit-breaker")
		duration := time.Since(start)

		result := TestResult{
			RequestNum: i + 1,
			Duration:   duration,
			Error:      err,
		}

		if err == nil {
			result.StatusCode = resp.StatusCode
			body, _ := io.ReadAll(resp.Body)

			// Parse response for circuit state
			var jsonResp map[string]interface{}
			if json.Unmarshal(body, &jsonResp) == nil {
				if state, ok := jsonResp["state"].(string); ok {
					if state == "open" {
						result.ResponseBody = "Circuit OPEN - fast fail"
					} else if state == "closed" {
						if jsonResp["error"] != nil {
							result.ResponseBody = "Circuit CLOSED - backend error"
						} else {
							result.ResponseBody = "Circuit CLOSED - success"
						}
					}
				}
			}
			resp.Body.Close()
		}

		results[i] = result
		printRequestResult(result)

		// Задержка между запросами
		time.Sleep(200 * time.Millisecond)
	}

	printSummary(results, "Circuit Breaker")
}

func TestFlakyService_503_Demo(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	printTestHeader("💥 Flaky Service Test - Random 503 Errors")

	const numRequests = 20
	results := make([]TestResult, numRequests)
	var success, failures atomic.Int32

	fmt.Println("  🎲 Testing endpoint with 50% failure rate...")
	fmt.Println()

	for i := 0; i < numRequests; i++ {
		start := time.Now()

		resp, err := http.Get(server.URL + "/api/flaky")
		duration := time.Since(start)

		result := TestResult{
			RequestNum: i + 1,
			Duration:   duration,
			Error:      err,
		}

		if err == nil {
			result.StatusCode = resp.StatusCode
			if resp.StatusCode == 200 {
				success.Add(1)
				result.ResponseBody = "Lucky! Got through"
			} else {
				failures.Add(1)
				body, _ := io.ReadAll(resp.Body)
				result.ResponseBody = string(body)
			}
			resp.Body.Close()
		}

		results[i] = result
		printRequestResult(result)

		time.Sleep(50 * time.Millisecond)
	}

	printSummary(results, "Flaky Service")

	// Дополнительная статистика для flaky service
	fmt.Printf("\n  %s🎰 Statistical Analysis:%s\n", colorBold, colorReset)
	fmt.Printf("    Expected: ~50%% success rate\n")
	fmt.Printf("    Actual: %.1f%% success rate\n", float64(success.Load())/float64(numRequests)*100)
}

func TestSlidingWindow_Demo(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	printTestHeader("⏱ Sliding Window Test - Time-based Rate Limiting")

	fmt.Println("  📋 Testing sliding window (limit: 20 req/min)...")
	fmt.Println("  💡 Sending 25 requests to exceed the minute limit")
	fmt.Println()

	results := make([]TestResult, 0, 25)

	// Отправляем первые 20 запросов быстро
	fmt.Printf("\n  %s--- Phase 1: Rapid burst (20 requests) ---%s\n\n", colorYellow, colorReset)
	for i := 0; i < 20; i++ {
		start := time.Now()
		resp, err := http.Get(server.URL + "/api/simple")
		duration := time.Since(start)

		result := TestResult{
			RequestNum:   i + 1,
			Duration:     duration,
			Error:        err,
		}

		if err == nil {
			result.StatusCode = resp.StatusCode
			resp.Body.Close()
		}

		results = append(results, result)
		printRequestResult(result)
		time.Sleep(20 * time.Millisecond)
	}

	// Пытаемся отправить ещё 5 запросов
	fmt.Printf("\n  %s--- Phase 2: Exceeding limit (5 more requests) ---%s\n\n", colorYellow, colorReset)
	for i := 20; i < 25; i++ {
		start := time.Now()
		resp, err := http.Get(server.URL + "/api/simple")
		duration := time.Since(start)

		result := TestResult{
			RequestNum:   i + 1,
			Duration:     duration,
			Error:        err,
		}

		if err == nil {
			result.StatusCode = resp.StatusCode
			if resp.StatusCode == 429 {
				body, _ := io.ReadAll(resp.Body)
				result.ResponseBody = string(body)
			}
			resp.Body.Close()
		}

		results = append(results, result)
		printRequestResult(result)
		time.Sleep(20 * time.Millisecond)
	}

	printSummary(results, "Sliding Window")
}