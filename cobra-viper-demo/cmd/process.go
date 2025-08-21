package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	mode int
	f    int
)

var processCmd = &cobra.Command{
	Use:                "process",
	Short:              "Обработка данных с параметрами",
	Long:               `Команда для обработки данных с поддержкой флагов в формате -mode:4 -f:2`,
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, arg := range args {
			if strings.HasPrefix(arg, "-mode:") {
				val := strings.TrimPrefix(arg, "-mode:")
				m, err := strconv.Atoi(val)
				if err != nil {
					return fmt.Errorf("неверное значение для mode: %s", val)
				}
				mode = m
			} else if strings.HasPrefix(arg, "-f:") {
				val := strings.TrimPrefix(arg, "-f:")
				fVal, err := strconv.Atoi(val)
				if err != nil {
					return fmt.Errorf("неверное значение для f: %s", val)
				}
				f = fVal
			} else if arg == "-h" || arg == "--help" {
				cmd.Help()
				return nil
			}
		}
		
		fmt.Printf("Запуск обработки:\n")
		fmt.Printf("  Mode: %d\n", mode)
		fmt.Printf("  F: %d\n", f)
		
		switch mode {
		case 1:
			fmt.Println("Режим 1: Базовая обработка")
		case 2:
			fmt.Println("Режим 2: Расширенная обработка")
		case 3:
			fmt.Println("Режим 3: Полная обработка")
		case 4:
			fmt.Println("Режим 4: Экспериментальная обработка")
		default:
			fmt.Printf("Режим %d: Пользовательский режим\n", mode)
		}
		
		if f > 0 {
			fmt.Printf("Применение фильтра уровня %d\n", f)
		}
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(processCmd)
}