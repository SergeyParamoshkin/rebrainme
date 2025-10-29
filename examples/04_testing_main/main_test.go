package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestGreet(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{"empty name", "", "Hello, World!"},
		{"with name", "Alice", "Hello, Alice!"},
		{"special chars", "Алиса", "Hello, Алиса!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := greet(tt.arg)
			if got != tt.want {
				t.Errorf("greet(%q) = %q; want %q", tt.arg, got, tt.want)
			}
		})
	}
}

func TestMainWithExitCode(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectedExit   int
	}{
		{
			name:           "default name",
			args:           []string{},
			expectedOutput: "Hello, World!\n",
			expectedExit:   0,
		},
		{
			name:           "custom name",
			args:           []string{"-name", "Alice"},
			expectedOutput: "Hello, Alice!\n",
			expectedExit:   0,
		},
		{
			name:           "verbose mode",
			args:           []string{"-v"},
			expectedOutput: "[VERBOSE] Starting application...\nHello, World!\n[VERBOSE] Application finished successfully\n",
			expectedExit:   0,
		},
		{
			name:           "repeat",
			args:           []string{"-repeat", "3"},
			expectedOutput: "Hello, World!\nHello, World!\nHello, World!\n",
			expectedExit:   0,
		},
		{
			name:           "invalid flag",
			args:           []string{"-invalid"},
			expectedOutput: "",
			expectedExit:   2,
		},
		{
			name:           "combined flags",
			args:           []string{"-name", "Bob", "-repeat", "2"},
			expectedOutput: "Hello, Bob!\nHello, Bob!\n",
			expectedExit:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var outBuf, errBuf bytes.Buffer

			oldStdout := stdout
			oldStderr := stderr
			oldArgs := args

			stdout = &outBuf
			stderr = &errBuf
			args = tt.args

			defer func() {
				stdout = oldStdout
				stderr = oldStderr
				args = oldArgs
			}()

			exitCode := mainWithExitCode()

			if exitCode != tt.expectedExit {
				t.Errorf("exit code = %d; want %d", exitCode, tt.expectedExit)
			}

			if tt.expectedOutput != "" {
				got := outBuf.String()
				if got != tt.expectedOutput {
					t.Errorf("output = %q; want %q", got, tt.expectedOutput)
				}
			}

			if tt.expectedExit != 0 && errBuf.Len() == 0 {
				t.Error("expected error output, got none")
			}
		})
	}
}

func TestMainVerboseOutput(t *testing.T) {
	var outBuf bytes.Buffer

	oldStdout := stdout
	oldArgs := args

	stdout = &outBuf
	args = []string{"-v", "-name", "Test"}

	defer func() {
		stdout = oldStdout
		args = oldArgs
	}()

	exitCode := mainWithExitCode()

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	output := outBuf.String()

	expectedParts := []string{
		"[VERBOSE] Starting application...",
		"Hello, Test!",
		"[VERBOSE] Application finished successfully",
	}

	for _, part := range expectedParts {
		if !strings.Contains(output, part) {
			t.Errorf("output missing expected part: %q\nfull output: %q", part, output)
		}
	}
}

func TestMainHelperFunction(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want int
	}{
		{"success", []string{"-name", "Alice"}, 0},
		{"invalid flag", []string{"-invalid"}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args = tt.args
			defer func() { args = os.Args[1:] }()

			got := mainWithExitCode()
			if got != tt.want {
				t.Errorf("mainWithExitCode() = %d; want %d", got, tt.want)
			}
		})
	}
}

// Example-тест для внутренней функции
func Example() {
	result := greet("World")
	fmt.Println(result)
	// Output: Hello, World!
}

func Example_emptyName() {
	result := greet("")
	fmt.Println(result)
	// Output: Hello, World!
}
