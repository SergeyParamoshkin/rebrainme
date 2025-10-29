package benchmarks

import (
	"testing"
)

func BenchmarkStringConcat(b *testing.B) {
	strs := []string{"hello", "world", "foo", "bar", "baz"}

	b.Run("Plus", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = ConcatWithPlus(strs)
		}
	})

	b.Run("Builder", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = ConcatWithBuilder(strs)
		}
	})

	b.Run("Join", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = ConcatWithJoin(strs)
		}
	})
}

func BenchmarkStringConcatSizes(b *testing.B) {
	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		strs := make([]string, size)
		for i := 0; i < size; i++ {
			strs[i] = "a"
		}

		b.Run("Plus_"+string(rune(size)), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = ConcatWithPlus(strs)
			}
		})

		b.Run("Builder_"+string(rune(size)), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = ConcatWithBuilder(strs)
			}
		})
	}
}

func BenchmarkFibonacci(b *testing.B) {
	benchmarks := []struct {
		name string
		n    int
	}{
		{"Fib10", 10},
		{"Fib15", 15},
		{"Fib20", 20},
	}

	for _, bm := range benchmarks {
		b.Run("Recursive_"+bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = Fibonacci(bm.n)
			}
		})

		b.Run("Iterative_"+bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = FibonacciIterative(bm.n)
			}
		})

		b.Run("Memoized_"+bm.name, func(b *testing.B) {
			fibCache = make(map[int]int)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = FibonacciMemo(bm.n)
			}
		})
	}
}

func BenchmarkFormatting(b *testing.B) {
	name := "John Doe"
	age := 30
	city := "New York"

	b.Run("Sprintf", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = FormatWithSprintf(name, age, city)
		}
	})

	b.Run("Builder", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = FormatWithBuilder(name, age, city)
		}
	})
}

func BenchmarkSearch(b *testing.B) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		slice := make([]int, size)
		m := make(map[int]bool, size)

		for i := 0; i < size; i++ {
			slice[i] = i
			m[i] = true
		}

		target := size - 1

		b.Run("Linear_"+string(rune(size)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = SearchLinear(slice, target)
			}
		})

		b.Run("Map_"+string(rune(size)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = SearchMap(m, target)
			}
		})
	}
}

func BenchmarkParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = Fibonacci(15)
		}
	})
}

func BenchmarkWithSetup(b *testing.B) {
	data := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = "test"
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = ConcatWithBuilder(data)
	}
}

func BenchmarkAllocs(b *testing.B) {
	b.ReportAllocs()

	strs := []string{"a", "b", "c", "d", "e"}

	for i := 0; i < b.N; i++ {
		_ = ConcatWithPlus(strs)
	}
}
