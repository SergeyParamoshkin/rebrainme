package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "cobra-viper-demo",
	Short: "Демонстрационное приложение для вебинара по Cobra и Viper",
	Long: `Это приложение демонстрирует основные возможности библиотек Cobra и Viper:
- Создание CLI приложений с подкомандами
- Работа с флагами и аргументами
- Конфигурирование через файлы, переменные окружения и флаги
- Приоритеты конфигурации`,
	Version: "1.0.0",
}

func Execute() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "путь к файлу конфигурации (по умолчанию ./config.yaml)")
	rootCmd.PersistentFlags().Bool("verbose", false, "включить подробный вывод")

	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
