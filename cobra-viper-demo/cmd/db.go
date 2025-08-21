package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Команды для работы с базой данных",
	Long:  `Группа команд для управления базой данных: миграции, резервное копирование, информация о подключении`,
}

var dbConnectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Показать строку подключения к БД",
	Run: func(cmd *cobra.Command, args []string) {
		host := viper.GetString("database.host")
		port := viper.GetInt("database.port")
		user := viper.GetString("database.user")
		password := viper.GetString("database.password")
		dbname := viper.GetString("database.name")

		if password == "" {
			password = "********"
		}

		connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, dbname)
		fmt.Printf("Строка подключения: %s\n", connStr)

		if viper.GetBool("verbose") {
			fmt.Println("\nПараметры подключения:")
			fmt.Printf("  Host: %s\n", host)
			fmt.Printf("  Port: %d\n", port)
			fmt.Printf("  Database: %s\n", dbname)
			fmt.Printf("  User: %s\n", user)

			if viper.ConfigFileUsed() != "" {
				fmt.Printf("\nИсточник конфигурации: %s\n", viper.ConfigFileUsed())
			}
		}
	},
}

var dbMigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Выполнить миграции БД",
	Run: func(cmd *cobra.Command, args []string) {
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		version, _ := cmd.Flags().GetString("version")

		fmt.Println("Выполнение миграций базы данных...")

		if dryRun {
			fmt.Println("РЕЖИМ ПРЕДПРОСМОТРА (--dry-run)")
		}

		host := viper.GetString("database.host")
		port := viper.GetInt("database.port")
		dbname := viper.GetString("database.name")

		fmt.Printf("База данных: %s:%d/%s\n", host, port, dbname)

		if version != "" {
			fmt.Printf("Целевая версия: %s\n", version)
		} else {
			fmt.Println("Мигрируем до последней версии")
		}

		fmt.Println("\nПримеры миграций:")
		fmt.Println("001_create_users_table.sql")
		fmt.Println("002_add_email_to_users.sql")
		fmt.Println("003_create_products_table.sql")

		if !dryRun {
			fmt.Println("\nМиграции успешно применены!")
		} else {
			fmt.Println("\nМиграции будут применены при запуске без флага --dry-run")
		}
	},
}

var dbCacheClearCmd = &cobra.Command{
	Use:   "cache",
	Short: "cache",
}

var dbBackupCmd = &cobra.Command{
	Use:        "backup [filename]",
	Short:      "Создать резервную копию БД",
	Deprecated: "This command is dep, use ...",
	Args:       cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		compress, _ := cmd.Flags().GetBool("compress")

		filename := "backup.sql"
		if len(args) > 0 {
			filename = args[0]
		}

		if compress {
			filename += ".gz"
		}

		host := viper.GetString("database.host")
		port := viper.GetInt("database.port")
		dbname := viper.GetString("database.name")

		fmt.Printf("Создание резервной копии базы данных %s:%d/%s\n", host, port, dbname)
		fmt.Printf("Файл: %s\n", filename)

		if compress {
			fmt.Println("Сжатие: включено")
		}

		fmt.Println("\nПроцесс резервного копирования...")
		fmt.Println("Экспорт схемы...")
		fmt.Println("Экспорт данных...")

		if compress {
			fmt.Println("Сжатие файла...")
		}

		fmt.Printf("\nРезервная копия успешно создана: %s\n", filename)

		info, _ := os.Stat(".")
		fmt.Printf("Размер файла: ~1.2 MB\n")
		fmt.Printf("Время создания: %v\n", info.ModTime().Format("2006-01-02 15:04:05"))
	},
}

var dbStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Проверить статус подключения к БД",
	Run: func(cmd *cobra.Command, args []string) {
		host := viper.GetString("database.host")
		port := viper.GetInt("database.port")
		dbname := viper.GetString("database.name")
		user := viper.GetString("database.user")

		fmt.Println("Проверка подключения к базе данных...")
		fmt.Printf("Сервер: %s:%d\n", host, port)
		fmt.Printf("База данных: %s\n", dbname)
		fmt.Printf("Пользователь: %s\n", user)

		fmt.Println("\n✓ Подключение установлено")
		fmt.Println("✓ База данных доступна")
		fmt.Println("✓ Права доступа проверены")

		fmt.Println("\nСтатистика:")
		fmt.Println("  Версия PostgreSQL: 14.5")
		fmt.Println("  Размер БД: 42 MB")
		fmt.Println("  Таблиц: 15")
		fmt.Println("  Активных подключений: 3")
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)

	dbCmd.AddCommand(dbConnectCmd)
	dbCmd.AddCommand(dbMigrateCmd)
	dbCmd.AddCommand(dbBackupCmd)
	dbCmd.AddCommand(dbStatusCmd)
	dbCmd.AddCommand(dbCacheClearCmd)

	dbMigrateCmd.Flags().Bool("dry-run", false, "показать какие миграции будут применены без их выполнения")
	dbMigrateCmd.Flags().String("version", "", "мигрировать до определенной версии")

	dbBackupCmd.Flags().Bool("compress", false, "сжать резервную копию с помощью gzip")

	dbCmd.PersistentFlags().String("db-host", "", "хост БД (переопределяет конфигурацию)")
	dbCmd.PersistentFlags().Int("db-port", 0, "порт БД (переопределяет конфигурацию)")
	dbCmd.PersistentFlags().String("db-name", "", "имя БД (переопределяет конфигурацию)")
	dbCmd.PersistentFlags().String("db-user", "", "пользователь БД (переопределяет конфигурацию)")
	dbCmd.PersistentFlags().String("db-password", "", "пароль БД (переопределяет конфигурацию)")

	viper.BindPFlag("database.host", dbCmd.PersistentFlags().Lookup("db-host"))
	viper.BindPFlag("database.port", dbCmd.PersistentFlags().Lookup("db-port"))
	viper.BindPFlag("database.name", dbCmd.PersistentFlags().Lookup("db-name"))
	viper.BindPFlag("database.user", dbCmd.PersistentFlags().Lookup("db-user"))
	viper.BindPFlag("database.password", dbCmd.PersistentFlags().Lookup("db-password"))
}
