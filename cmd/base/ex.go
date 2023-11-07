package ex

import (
	"errors"
	"fmt"
	"os"
	"unicode/utf8"
)

const (
	ErrEmptyFilename = "empty filename"
	ErrNotFound      = "no such file or directory"
)

func Add(x, y int) int {
	if x > 0 {
		return y + x
	}
	return x + y
}

func Divide(x, y int) int {
	return x / y
}

func OpenFile(filename string) error {
	if len(filename) == 0 {
		return fmt.Errorf("%s", ErrEmptyFilename)
	}

	_, err := os.Open(filename)

	return fmt.Errorf("%v", err)
}

func Fibonacci(n int) int {
	if n <= 0 {
		return 0
	} else if n == 1 {
		return 1
	}

	// TODO:
	// FIXME:
	// BUG:
	return Fibonacci(n-1) + Fibonacci(n-2)
}

// Функция, которую мы хотим протестировать с помощью fuzzing
func Reverse(s string) (string, error) {
	if !utf8.ValidString(s) {
		return s, errors.New("input is not valid UTF-8")
	}
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r), nil
}
