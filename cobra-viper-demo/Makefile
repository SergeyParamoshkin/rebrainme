.PHONY: help build run clean test deps serve

help: ## Показать помощь
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Собрать приложение
	go build -o app main.go

run: ## Запустить приложение
	go run main.go

serve: ## Запустить HTTP сервер
	go run main.go serve --verbose

clean: ## Очистить сборку
	rm -f app
	go clean

test: ## Запустить тесты
	go test -v ./...

deps: ## Установить зависимости
	go mod download
	go mod tidy

demo-config: ## Демонстрация работы с конфигурацией
	@echo "=== Показать конфигурацию ==="
	go run main.go config show
	@echo "\n=== Получить значение server.port ==="
	go run main.go config get server.port
	@echo "\n=== Список всех ключей ==="
	go run main.go config list

demo-db: ## Демонстрация команд БД
	@echo "=== Статус БД ==="
	go run main.go db status
	@echo "\n=== Строка подключения ==="
	go run main.go db connect
	@echo "\n=== Миграции (dry-run) ==="
	go run main.go db migrate --dry-run

demo-env: ## Демонстрация с переменными окружения
	DEMO_SERVER_PORT=9999 DEMO_DATABASE_HOST=remote-db.com go run main.go config show

demo-priority: ## Демонстрация приоритетов конфигурации
	@echo "=== По умолчанию (из config.yaml) ==="
	go run main.go config get server.port
	@echo "\n=== С переменной окружения ==="
	DEMO_SERVER_PORT=9000 go run main.go config get server.port
	@echo "\n=== С флагом командной строки (высший приоритет) ==="
	DEMO_SERVER_PORT=9000 go run main.go serve --port 7000 --verbose | head -n 5