package cmd

import (
	"strings"
	"testing"
	"os/exec"
)

// Простые тесты через exec для всех команд
func TestCommandHelp(t *testing.T) {
	tests := []struct {
		name     string
		cmd      []string
		contains []string
	}{
		{
			name: "Root help",
			cmd:  []string{"--help"},
			contains: []string{
				"основные возможности библиотек Cobra и Viper",
				"Available Commands:",
				"process",
				"advanced",
				"config",
				"db",
				"serve",
			},
		},
		{
			name: "Config help",
			cmd:  []string{"config", "--help"},
			contains: []string{
				"просмотра и управления конфигурацией",
				"Available Commands:",
				"show",
				"get",
				"list",
			},
		},
		{
			name: "Config show help",
			cmd:  []string{"config", "show", "--help"},
			contains: []string{
				"Показать текущую конфигурацию",
				"--format string",
				"table, json, yaml",
			},
		},
		{
			name: "DB help",
			cmd:  []string{"db", "--help"},
			contains: []string{
				"управления базой данных",
				"Available Commands:",
				"connect",
				"migrate",
				"status",
			},
		},
		{
			name: "DB migrate help",
			cmd:  []string{"db", "migrate", "--help"},
			contains: []string{
				"Выполнить миграции БД",
				"--dry-run",
				"--version",
			},
		},
		{
			name: "Serve help",
			cmd:  []string{"serve", "--help"},
			contains: []string{
				"HTTP сервер с настройками",
				"--host",
				"--port",
				"--graceful",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("../test-cli", tt.cmd...)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Errorf("Command failed: %v", err)
			}

			outputStr := string(output)
			for _, expected := range tt.contains {
				if !strings.Contains(outputStr, expected) {
					t.Errorf("Expected help to contain '%s', but got: %s", expected, outputStr)
				}
			}
		})
	}
}

func TestConfigCommands(t *testing.T) {
	tests := []struct {
		name     string
		cmd      []string
		contains []string
	}{
		{
			name: "Config show table format",
			cmd:  []string{"config", "show"},
			contains: []string{
				"Текущая конфигурация:",
				"СЕРВЕР:",
				"БАЗА ДАННЫХ:",
				"Host: localhost",
				"Port: 8080",
			},
		},
		{
			name: "Config show JSON format",
			cmd:  []string{"config", "show", "--format", "json"},
			contains: []string{
				`"server"`,
				`"host"`,
				`"port"`,
				`"database"`,
			},
		},
		{
			name: "Config show YAML format",
			cmd:  []string{"config", "show", "--format", "yaml"},
			contains: []string{
				"server:",
				"host:",
				"port:",
				"database:",
			},
		},
		{
			name: "Config list",
			cmd:  []string{"config", "list"},
			contains: []string{
				"Доступные ключи конфигурации:",
				"server.host",
				"server.port",
				"database.host",
			},
		},
		{
			name: "Config get existing key",
			cmd:  []string{"config", "get", "server.host"},
			contains: []string{
				"server.host = localhost",
			},
		},
		{
			name: "Config get non-existing key",
			cmd:  []string{"config", "get", "non.existing.key"},
			contains: []string{
				"Ключ 'non.existing.key' не найден",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("../test-cli", tt.cmd...)
			output, err := cmd.CombinedOutput()
			if err != nil && !strings.Contains(string(output), "не найден") {
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

func TestDBCommands(t *testing.T) {
	tests := []struct {
		name     string
		cmd      []string
		contains []string
	}{
		{
			name: "DB connect",
			cmd:  []string{"db", "connect"},
			contains: []string{
				"Строка подключения: postgres://",
				"localhost:5432",
			},
		},
		{
			name: "DB status",
			cmd:  []string{"db", "status"},
			contains: []string{
				"Проверка подключения к базе данных",
				"✓ Подключение установлено",
				"✓ База данных доступна",
				"Статистика:",
			},
		},
		{
			name: "DB migrate dry run",
			cmd:  []string{"db", "migrate", "--dry-run"},
			contains: []string{
				"Выполнение миграций базы данных",
				"РЕЖИМ ПРЕДПРОСМОТРА (--dry-run)",
				"001_create_users_table.sql",
				"Миграции будут применены при запуске без флага --dry-run",
			},
		},
		{
			name: "DB backup default",
			cmd:  []string{"db", "backup"},
			contains: []string{
				"Создание резервной копии базы данных",
				"Файл: backup.sql",
				"Экспорт схемы...",
				"Резервная копия успешно создана",
			},
		},
		{
			name: "DB backup with compression",
			cmd:  []string{"db", "backup", "--compress"},
			contains: []string{
				"Файл: backup.sql.gz",
				"Сжатие: включено",
				"Сжатие файла...",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("../test-cli", tt.cmd...)
			output, err := cmd.CombinedOutput()
			if err != nil {
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

func TestVersionAndFlags(t *testing.T) {
	// Тест версии
	cmd := exec.Command("../test-cli", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Version command failed: %v", err)
	}

	if !strings.Contains(string(output), "1.0.0") {
		t.Errorf("Expected version to contain '1.0.0', got: %s", output)
	}

	// Тест verbose флага
	cmd = exec.Command("../test-cli", "--verbose", "config", "get", "server.host")
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Verbose command failed: %v", err)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Тип:") || !strings.Contains(outputStr, "Источник:") {
		t.Errorf("Expected verbose output to contain type and source info, got: %s", outputStr)
	}
}