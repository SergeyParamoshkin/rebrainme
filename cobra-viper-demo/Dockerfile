# Multi-stage build для оптимизации размера образа

# Стадия сборки
FROM golang:1.21-alpine AS builder

# Установка необходимых пакетов для сборки
RUN apk add --no-cache git make

# Рабочая директория
WORKDIR /build

# Копирование файлов модулей и загрузка зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=$(git describe --tags --always --dirty)" \
    -o app \
    main.go

# Финальный образ
FROM alpine:3.18

# Установка ca-certificates для HTTPS запросов и timezone data
RUN apk --no-cache add ca-certificates tzdata

# Создание пользователя для запуска приложения
RUN addgroup -g 1000 -S appuser && \
    adduser -u 1000 -S appuser -G appuser

# Рабочая директория
WORKDIR /app

# Копирование собранного приложения из builder
COPY --from=builder /build/app .
COPY --from=builder /build/config.yaml .

# Изменение владельца файлов
RUN chown -R appuser:appuser /app

# Переключение на непривилегированного пользователя
USER appuser

# Порт приложения
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Запуск приложения
ENTRYPOINT ["./app"]
CMD ["serve"]