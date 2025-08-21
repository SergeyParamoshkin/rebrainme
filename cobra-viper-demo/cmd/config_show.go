package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Управление конфигурацией",
	Long:  `Команды для просмотра и управления конфигурацией приложения`,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Показать текущую конфигурацию",
	Run: func(cmd *cobra.Command, args []string) {
		format, _ := cmd.Flags().GetString("format")
		
		switch format {
		case "json":
			showConfigJSON()
		case "yaml":
			showConfigYAML()
		default:
			showConfigTable()
		}
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get KEY",
	Short: "Получить значение конкретного параметра",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := viper.Get(key)
		
		if value == nil {
			fmt.Printf("Ключ '%s' не найден\n", key)
			return
		}
		
		fmt.Printf("%s = %v\n", key, value)
		
		if viper.GetBool("verbose") {
			fmt.Printf("\nТип: %T\n", value)
			if viper.IsSet(key) {
				fmt.Println("Источник: установлен явно")
			} else {
				fmt.Println("Источник: значение по умолчанию")
			}
		}
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список всех доступных ключей конфигурации",
	Run: func(cmd *cobra.Command, args []string) {
		allKeys := viper.AllKeys()
		sort.Strings(allKeys)
		
		fmt.Println("Доступные ключи конфигурации:")
		fmt.Println(strings.Repeat("-", 40))
		
		currentPrefix := ""
		for _, key := range allKeys {
			parts := strings.Split(key, ".")
			if len(parts) > 0 && parts[0] != currentPrefix {
				currentPrefix = parts[0]
				fmt.Printf("\n[%s]\n", currentPrefix)
			}
			fmt.Printf("  %s\n", key)
		}
	},
}

func showConfigTable() {
	fmt.Println("Текущая конфигурация:")
	fmt.Println(strings.Repeat("=", 60))
	
	fmt.Println("\nСЕРВЕР:")
	fmt.Printf("  Host: %s\n", viper.GetString("server.host"))
	fmt.Printf("  Port: %d\n", viper.GetInt("server.port"))
	fmt.Printf("  Graceful Shutdown: %v\n", viper.GetBool("server.graceful"))
	
	fmt.Println("\nБАЗА ДАННЫХ:")
	fmt.Printf("  Host: %s\n", viper.GetString("database.host"))
	fmt.Printf("  Port: %d\n", viper.GetInt("database.port"))
	fmt.Printf("  Database: %s\n", viper.GetString("database.name"))
	fmt.Printf("  User: %s\n", viper.GetString("database.user"))
	
	password := viper.GetString("database.password")
	if password != "" {
		fmt.Printf("  Password: %s\n", strings.Repeat("*", len(password)))
	}
	
	fmt.Println("\nЛОГИРОВАНИЕ:")
	fmt.Printf("  Level: %s\n", viper.GetString("logging.level"))
	fmt.Printf("  Format: %s\n", viper.GetString("logging.format"))
	fmt.Printf("  Output: %s\n", viper.GetString("logging.output"))
	
	fmt.Println("\nПРИЛОЖЕНИЕ:")
	fmt.Printf("  Name: %s\n", viper.GetString("app.name"))
	fmt.Printf("  Version: %s\n", viper.GetString("app.version"))
	fmt.Printf("  Environment: %s\n", viper.GetString("app.environment"))
	fmt.Printf("  Debug: %v\n", viper.GetBool("app.debug"))
	
	if viper.ConfigFileUsed() != "" {
		fmt.Printf("\nФайл конфигурации: %s\n", viper.ConfigFileUsed())
	}
	
	fmt.Println(strings.Repeat("=", 60))
}

func showConfigJSON() {
	settings := viper.AllSettings()
	jsonBytes, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		fmt.Printf("Ошибка при формировании JSON: %v\n", err)
		return
	}
	fmt.Println(string(jsonBytes))
}

func showConfigYAML() {
	settings := viper.AllSettings()
	printYAML(settings, 0)
}

func printYAML(data interface{}, indent int) {
	indentStr := strings.Repeat("  ", indent)
	
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			fmt.Printf("%s%s:", indentStr, key)
			if _, ok := value.(map[string]interface{}); ok {
				fmt.Println()
				printYAML(value, indent+1)
			} else {
				fmt.Printf(" %v\n", value)
			}
		}
	default:
		fmt.Printf("%s%v\n", indentStr, v)
	}
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configListCmd)
	
	configShowCmd.Flags().String("format", "table", "формат вывода: table, json, yaml")
}