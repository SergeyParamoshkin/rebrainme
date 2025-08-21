# Cobra + Viper Demo Application

Демонстрационное приложение для вебинара по использованию библиотек Cobra и Viper в Go.

## Описание

Это приложение демонстрирует основные возможности двух популярных библиотек экосистемы Go:

- **[Cobra](https://github.com/spf13/cobra)** - мощная библиотека для создания CLI-приложений
- **[Viper](https://github.com/spf13/viper)** - универсальная система управления конфигурацией

## Возможности приложения

### 1. CLI с подкомандами (Cobra)
- Структурированная организация команд
- Автоматическая генерация help
- Поддержка флагов и аргументов
- Вложенные подкоманды

### 2. Управление конфигурацией (Viper)
- Чтение из YAML файлов
- Переменные окружения
- Флаги командной строки
- Значения по умолчанию
- Приоритеты источников конфигурации

## Установка

```bash
# Клонировать репозиторий
git clone <repository-url>
cd cobra-viper-demo

# Установить зависимости
go mod download

# Собрать приложение
go build -o app
```

## Структура проекта

```
cobra-viper-demo/
├── main.go              # Точка входа
├── cmd/
│   ├── root.go         # Корневая команда
│   ├── config.go       # Инициализация конфигурации
│   ├── config_show.go  # Команды для работы с конфигурацией
│   ├── serve.go        # HTTP сервер
│   └── db.go          # Команды базы данных
├── config.yaml         # Файл конфигурации
├── .env.example        # Пример переменных окружения
├── go.mod             # Зависимости
└── README.md          # Документация
```

## Использование

### Основные команды

```bash
# Показать помощь
./app --help

# Показать версию
./app --version

# Запустить с подробным выводом
./app --verbose <command>
```

### Работа с HTTP сервером

```bash
# Запустить сервер с настройками по умолчанию
./app serve

# Запустить на другом порту
./app serve --port 3000

# Запустить с кастомным хостом и портом
./app serve --host 0.0.0.0 --port 8090

# Отключить graceful shutdown
./app serve --graceful=false

# Использовать другой конфигурационный файл
./app --config custom-config.yaml serve
```

Endpoints сервера:
- `GET /` - главная страница с информацией о конфигурации
- `GET /health` - проверка состояния сервера
- `GET /config` - текущая конфигурация в JSON

### Работа с базой данных

```bash
# Показать строку подключения
./app db connect

# Проверить статус подключения
./app db status

# Выполнить миграции
./app db migrate

# Предпросмотр миграций без выполнения
./app db migrate --dry-run

# Мигрировать до конкретной версии
./app db migrate --version 003

# Создать резервную копию
./app db backup

# Создать сжатую резервную копию
./app db backup --compress backup_2024.sql

# Переопределить параметры подключения
./app db connect --db-host=192.168.1.100 --db-port=5433
```

### Управление конфигурацией

```bash
# Показать текущую конфигурацию
./app config show

# Показать в формате JSON
./app config show --format json

# Показать в формате YAML
./app config show --format yaml

# Получить конкретное значение
./app config get server.port
./app config get database.host

# Показать все доступные ключи
./app config list
```

## Приоритеты конфигурации

Viper применяет следующий порядок приоритетов (от высшего к низшему):

1. **Флаги командной строки**
   ```bash
   ./app serve --port 9000
   ```

2. **Переменные окружения**
   ```bash
   export DEMO_SERVER_PORT=8080
   ./app serve
   ```

3. **Конфигурационный файл**
   ```yaml
   server:
     port: 8080
   ```

4. **Значения по умолчанию**
   ```go
   viper.SetDefault("server.port", 8080)
   ```

## Переменные окружения

Приложение автоматически читает переменные с префиксом `DEMO_`:

```bash
# Настройки сервера
export DEMO_SERVER_HOST=0.0.0.0
export DEMO_SERVER_PORT=3000

# Настройки БД
export DEMO_DATABASE_HOST=postgres.example.com
export DEMO_DATABASE_PORT=5432
export DEMO_DATABASE_NAME=production
export DEMO_DATABASE_USER=admin
export DEMO_DATABASE_PASSWORD=secret

# Настройки логирования
export DEMO_LOGGING_LEVEL=debug

# Запустить с переменными окружения
./app serve --verbose
```

## Примеры использования для вебинара

### Пример 1: Демонстрация приоритетов

```bash
# 1. Запустить с конфигом по умолчанию
./app serve --verbose

# 2. Переопределить через переменную окружения
DEMO_SERVER_PORT=9000 ./app serve --verbose

# 3. Переопределить через флаг (высший приоритет)
DEMO_SERVER_PORT=9000 ./app serve --port 7000 --verbose
```

### Пример 2: Работа с несколькими конфигами

```bash
# Создать конфиг для разработки
cp config.yaml config.dev.yaml
# Отредактировать config.dev.yaml

# Создать конфиг для продакшена
cp config.yaml config.prod.yaml
# Отредактировать config.prod.yaml

# Использовать разные конфиги
./app --config config.dev.yaml serve
./app --config config.prod.yaml serve
```

### Пример 3: Комплексная демонстрация

```bash
# Проверить конфигурацию
./app config show

# Запустить сервер
./app serve --verbose &

# Проверить здоровье сервера
curl http://localhost:8080/health

# Получить конфигурацию через API
curl http://localhost:8080/config

# Проверить БД
./app db status

# Выполнить миграции
./app db migrate --dry-run
```

## Расширение приложения

### Добавление новой команды

1. Создайте файл `cmd/mycommand.go`:
```go
package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Описание команды",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Выполнение команды")
    },
}

func init() {
    rootCmd.AddCommand(myCmd)
    myCmd.Flags().String("param", "", "описание параметра")
}
```

### Добавление новых параметров конфигурации

1. Добавьте в `config.yaml`:
```yaml
myfeature:
  enabled: true
  timeout: 30s
```

2. Установите значения по умолчанию в `cmd/config.go`:
```go
viper.SetDefault("myfeature.enabled", false)
viper.SetDefault("myfeature.timeout", "10s")
```

3. Используйте в коде:
```go
if viper.GetBool("myfeature.enabled") {
    timeout := viper.GetDuration("myfeature.timeout")
    // ...
}
```

## Полезные материалы

- [Документация Cobra](https://cobra.dev/)
- [Документация Viper](https://github.com/spf13/viper)
- [Примеры использования Cobra](https://github.com/spf13/cobra/blob/master/user_guide.md)
- [12-Factor App Config](https://12factor.net/config)

## Часто задаваемые вопросы

### Как Viper находит конфигурационный файл?
Viper ищет файл в текущей директории с именем `config` и расширениями `.json`, `.toml`, `.yaml`, `.yml`, `.properties`, `.props`, `.prop`, `.hcl`, `.tfvars`, `.dotenv`, `.env`, `.ini`.

### Как отладить загрузку конфигурации?
Используйте флаг `--verbose` для подробного вывода:
```bash
./app --verbose serve
```

### Можно ли использовать JSON вместо YAML?
Да, просто переименуйте файл в `config.json` и используйте JSON формат.

### Как организовать конфигурацию для разных окружений?
1. Используйте разные файлы конфигурации
2. Используйте переменные окружения
3. Комбинируйте оба подхода

## Лицензия

MIT

## Автор

Создано для демонстрации на вебинаре по Golang.