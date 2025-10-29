//go:build integration
//+build integration

package basic

import "testing"

func TestAdd(t *testing.T) {
	result := Add(2, 3)
	expected := 5

	if result != expected {
		t.Errorf("Add(2, 3) = %d; want %d", result, expected)
	}
}

func TestSubtract(t *testing.T) {
	result := Subtract(5, 3)
	expected := 2

	if result != expected {
		t.Errorf("Subtract(5, 3) = %d; want %d", result, expected)
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"positive", 2, 3, 6},
		{"negative", -2, 3, -6},
		{"zero", 0, 5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Multiply(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Multiply(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	t.Run("normal division", func(t *testing.T) {
		result, err := Divide(10, 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := 5.0
		if result != expected {
			t.Errorf("Divide(10, 2) = %f; want %f", result, expected)
		}
	})

	t.Run("division by zero", func(t *testing.T) {
		_, err := Divide(10, 0)
		if err == nil {
			t.Error("expected error for division by zero, got nil")
		}
	})
}
