package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type AdvancedConfig struct {
	Mode    int
	Filter  int
	Debug   bool
	Output  string
	Threads int
}

var advancedCmd = &cobra.Command{
	Use:   "advanced [аргументы]",
	Short: "Продвинутая команда с кастомными флагами",
	Long: `Поддерживает флаги в формате:
  -mode:N     режим работы (1-10)
  -f:N        уровень фильтрации (0-5)
  -debug:1    включить отладку (0/1)
  -output:path путь для вывода
  -threads:N   количество потоков (1-16)
  
Примеры:
  advanced -mode:4 -f:2
  advanced -mode:5 -f:3 -debug:1 -output:/tmp/result.txt
  advanced -threads:8 -mode:2`,
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		config := &AdvancedConfig{
			Mode:    1,
			Filter:  0,
			Debug:   false,
			Output:  "stdout",
			Threads: 1,
		}
		
		for _, arg := range args {
			if arg == "-h" || arg == "--help" {
				cmd.Help()
				return nil
			}
			
			if !strings.HasPrefix(arg, "-") {
				continue
			}
			
			parts := strings.SplitN(arg[1:], ":", 2)
			if len(parts) != 2 {
				return fmt.Errorf("неверный формат флага: %s (ожидается -name:value)", arg)
			}
			
			flag := parts[0]
			value := parts[1]
			
			switch flag {
			case "mode":
				val, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("mode должен быть числом: %s", value)
				}
				if val < 1 || val > 10 {
					return fmt.Errorf("mode должен быть от 1 до 10, получено: %d", val)
				}
				config.Mode = val
				
			case "f":
				val, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("f должен быть числом: %s", value)
				}
				if val < 0 || val > 5 {
					return fmt.Errorf("f должен быть от 0 до 5, получено: %d", val)
				}
				config.Filter = val
				
			case "debug":
				val, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("debug должен быть 0 или 1: %s", value)
				}
				config.Debug = val == 1
				
			case "output":
				config.Output = value
				
			case "threads":
				val, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("threads должен быть числом: %s", value)
				}
				if val < 1 || val > 16 {
					return fmt.Errorf("threads должен быть от 1 до 16, получено: %d", val)
				}
				config.Threads = val
				
			default:
				return fmt.Errorf("неизвестный флаг: %s", flag)
			}
		}
		
		return runAdvanced(config)
	},
}

func runAdvanced(config *AdvancedConfig) error {
	fmt.Println("╔════════════════════════════════════╗")
	fmt.Println("║     Продвинутая обработка         ║")
	fmt.Println("╚════════════════════════════════════╝")
	fmt.Println()
	
	fmt.Println("Конфигурация:")
	fmt.Printf("  • Режим: %d", config.Mode)
	switch config.Mode {
	case 1:
		fmt.Println(" (базовый)")
	case 2:
		fmt.Println(" (стандартный)")
	case 3:
		fmt.Println(" (расширенный)")
	case 4:
		fmt.Println(" (экспериментальный)")
	case 5:
		fmt.Println(" (производительный)")
	default:
		fmt.Println(" (пользовательский)")
	}
	
	fmt.Printf("  • Фильтр: %d", config.Filter)
	if config.Filter == 0 {
		fmt.Println(" (отключен)")
	} else {
		fmt.Printf(" (уровень %d)\n", config.Filter)
	}
	
	fmt.Printf("  • Отладка: %v\n", config.Debug)
	fmt.Printf("  • Вывод: %s\n", config.Output)
	fmt.Printf("  • Потоки: %d\n", config.Threads)
	
	if config.Debug {
		fmt.Println("\n[DEBUG] Дополнительная информация:")
		fmt.Printf("[DEBUG] PID: %d\n", os.Getpid())
		fmt.Printf("[DEBUG] Аргументы: %v\n", os.Args)
	}
	
	fmt.Println("\n▸ Начало обработки...")
	
	if config.Output != "stdout" {
		fmt.Printf("▸ Результаты будут сохранены в: %s\n", config.Output)
	}
	
	if config.Threads > 1 {
		fmt.Printf("▸ Используется параллельная обработка (%d потоков)\n", config.Threads)
	}
	
	fmt.Println("✓ Обработка завершена успешно")
	
	return nil
}

func init() {
	rootCmd.AddCommand(advancedCmd)
}