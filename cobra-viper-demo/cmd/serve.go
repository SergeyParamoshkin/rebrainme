package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Запуск HTTP сервера",
	Long:  `Запускает HTTP сервер с настройками из конфигурации`,
	Run:   runServer,
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().String("host", "", "хост сервера (переопределяет конфигурацию)")
	serveCmd.Flags().Int("port", 0, "порт сервера (переопределяет конфигурацию)")
	serveCmd.Flags().Bool("graceful", true, "включить graceful shutdown")

	viper.BindPFlag("server.host", serveCmd.Flags().Lookup("host"))
	viper.BindPFlag("server.port", serveCmd.Flags().Lookup("port"))
	viper.BindPFlag("server.graceful", serveCmd.Flags().Lookup("graceful"))
}

func runServer(cmd *cobra.Command, args []string) {
	host := viper.GetString("server.host")
	port := viper.GetInt("server.port")
	graceful := viper.GetBool("server.graceful")
	verbose := viper.GetBool("verbose")

	if verbose {
		fmt.Println("Конфигурация сервера:")
		fmt.Printf("  Host: %s\n", host)
		fmt.Printf("  Port: %d\n", port)
		fmt.Printf("  Graceful shutdown: %v\n", graceful)
		if viper.ConfigFileUsed() != "" {
			fmt.Printf("  Конфигурационный файл: %s\n", viper.ConfigFileUsed())
		}
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Привет от Cobra и Viper!\n\n")
		fmt.Fprintf(w, "Текущая конфигурация:\n")
		fmt.Fprintf(w, "Server: %s:%d\n", host, port)
		fmt.Fprintf(w, "Database: %s:%d/%s\n",
			viper.GetString("database.host"),
			viper.GetInt("database.port"),
			viper.GetString("database.name"))
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	mux.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		settings := viper.AllSettings()
		fmt.Fprintf(w, "%v\n", settings)
	})

	addr := fmt.Sprintf("%s:%d", host, port)
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	if graceful {
		go func() {
			sigint := make(chan os.Signal, 1)
			signal.Notify(sigint, os.Interrupt)
			<-sigint

			fmt.Println("\nПолучен сигнал прерывания, завершаем работу...")

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			if err := srv.Shutdown(ctx); err != nil {
				log.Printf("Ошибка при завершении сервера: %v", err)
			}
		}()
	}

	fmt.Printf("Сервер запущен на http://%s\n", addr)
	fmt.Println("Нажмите Ctrl+C для остановки")

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}

	fmt.Println("Сервер остановлен")
}
