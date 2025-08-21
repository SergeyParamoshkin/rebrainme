package cmd

import (
	"os/exec"
	"strings"
	"testing"
)

func TestIntegrationProcessCommand(t *testing.T) {
	// Сначала собираем бинарник
	buildCmd := exec.Command("go", "build", "-o", "../test-cli", "../main.go")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build: %v", err)
	}

	tests := []struct {
		name     string
		args     []string
		contains []string
	}{
		{
			name: "Process with mode and filter",
			args: []string{"process", "-mode:4", "-f:2"},
			contains: []string{
				"Mode: 4",
				"F: 2",
				"Режим 4: Экспериментальная обработка",
				"Применение фильтра уровня 2",
			},
		},
		{
			name: "Process with only mode",
			args: []string{"process", "-mode:3"},
			contains: []string{
				"Mode: 3",
				"F: 0",
				"Режим 3: Полная обработка",
			},
		},
		{
			name: "Advanced with all flags",
			args: []string{"advanced", "-mode:5", "-f:3", "-debug:1", "-threads:4"},
			contains: []string{
				"Режим: 5 (производительный)",
				"Фильтр: 3 (уровень 3)",
				"Отладка: true",
				"Потоки: 4",
				"[DEBUG]",
				"Используется параллельная обработка (4 потоков)",
			},
		},
		{
			name: "Advanced with custom output",
			args: []string{"advanced", "-output:/tmp/test.txt"},
			contains: []string{
				"Вывод: /tmp/test.txt",
				"Результаты будут сохранены в: /tmp/test.txt",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("../test-cli", tt.args...)
			output, err := cmd.CombinedOutput()
			if err != nil {
				// Проверяем, что это не ошибка валидации
				outputStr := string(output)
				if strings.Contains(outputStr, "Error:") {
					// Это ожидаемая ошибка валидации
					return
				}
				t.Errorf("Command failed: %v\nOutput: %s", err, output)
			}

			outputStr := string(output)
			for _, expected := range tt.contains {
				if !strings.Contains(outputStr, expected) {
					t.Errorf("Expected output to contain '%s', but got: %s", expected, outputStr)
				}
			}
		})
	}
}

func TestIntegrationErrorCases(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		shouldError bool
		contains    string
	}{
		{
			name:        "Process invalid mode",
			args:        []string{"process", "-mode:abc"},
			shouldError: true,
			contains:    "неверное значение для mode",
		},
		{
			name:        "Process invalid filter",
			args:        []string{"process", "-f:xyz"},
			shouldError: true,
			contains:    "неверное значение для f",
		},
		{
			name:        "Advanced mode out of range",
			args:        []string{"advanced", "-mode:11"},
			shouldError: true,
			contains:    "mode должен быть от 1 до 10",
		},
		{
			name:        "Advanced threads out of range",
			args:        []string{"advanced", "-threads:20"},
			shouldError: true,
			contains:    "threads должен быть от 1 до 16",
		},
		{
			name:        "Advanced unknown flag",
			args:        []string{"advanced", "-unknown:value"},
			shouldError: true,
			contains:    "неизвестный флаг",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("../test-cli", tt.args...)
			output, err := cmd.CombinedOutput()
			
			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error but command succeeded")
				}
				outputStr := string(output)
				if !strings.Contains(outputStr, tt.contains) {
					t.Errorf("Expected error output to contain '%s', but got: %s", tt.contains, outputStr)
				}
			} else {
				if err != nil {
					t.Errorf("Expected success but got error: %v\nOutput: %s", err, output)
				}
			}
		})
	}
}