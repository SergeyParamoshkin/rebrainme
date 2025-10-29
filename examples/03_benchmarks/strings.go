package benchmarks

import (
	"fmt"
	"strings"
)

func ConcatWithPlus(strs []string) string {
	result := ""
	for _, s := range strs {
		result += s
	}
	return result
}

func ConcatWithBuilder(strs []string) string {
	var builder strings.Builder
	for _, s := range strs {
		builder.WriteString(s)
	}
	return builder.String()
}

func ConcatWithJoin(strs []string) string {
	return strings.Join(strs, "")
}

func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}

var fibCache = make(map[int]int)

func FibonacciMemo(n int) int {
	if n <= 1 {
		return n
	}

	if val, ok := fibCache[n]; ok {
		return val
	}

	result := FibonacciMemo(n-1) + FibonacciMemo(n-2)
	fibCache[n] = result
	return result
}

func FibonacciIterative(n int) int {
	if n <= 1 {
		return n
	}

	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

func FormatWithSprintf(name string, age int, city string) string {
	return fmt.Sprintf("Name: %s, Age: %d, City: %s", name, age, city)
}

func FormatWithBuilder(name string, age int, city string) string {
	var builder strings.Builder
	builder.WriteString("Name: ")
	builder.WriteString(name)
	builder.WriteString(", Age: ")
	builder.WriteString(fmt.Sprintf("%d", age))
	builder.WriteString(", City: ")
	builder.WriteString(city)
	return builder.String()
}

func SearchLinear(slice []int, target int) int {
	for i, v := range slice {
		if v == target {
			return i
		}
	}
	return -1
}

func SearchMap(m map[int]bool, target int) bool {
	return m[target]
}
