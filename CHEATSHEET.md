# Docker для Golang - Шпаргалка

## Основные команды Docker

### Сборка образов

```bash
# Базовая сборка
docker build -t myapp:latest .

# Сборка с конкретным Dockerfile
docker build -t myapp:latest -f Dockerfile.optimized .

# Сборка с build arguments
docker build --build-arg VERSION=1.0.0 -t myapp:1.0.0 .

# Сборка без кэша
docker build --no-cache -t myapp:latest .

# Multi-platform сборка
docker buildx build --platform linux/amd64,linux/arm64 -t myapp:latest .
```

### Запуск контейнеров

```bash
# Простой запуск
docker run myapp:latest

# С пробросом портов
docker run -p 8080:8080 myapp:latest

# В фоновом режиме
docker run -d --name myapp -p 8080:8080 myapp:latest

# С переменными окружения
docker run -e DB_HOST=localhost -e DB_PORT=5432 myapp:latest

# С монтированием volumes
docker run -v $(pwd)/data:/app/data myapp:latest

# Интерактивный режим
docker run -it myapp:latest /bin/sh
```

### Управление образами

```bash
# Список образов
docker images

# Удалить образ
docker rmi myapp:latest

# Удалить все неиспользуемые образы
docker image prune -a

# Информация об образе
docker inspect myapp:latest

# История слоев
docker history myapp:latest

# Размер образа
docker images myapp --format "{{.Repository}}:{{.Tag}} - {{.Size}}"

# Сохранить образ в файл
docker save myapp:latest -o myapp.tar

# Загрузить образ из файла
docker load -i myapp.tar
```

### Работа с контейнерами

```bash
# Список запущенных контейнеров
docker ps

# Список всех контейнеров
docker ps -a

# Остановить контейнер
docker stop myapp

# Запустить остановленный контейнер
docker start myapp

# Перезапустить контейнер
docker restart myapp

# Удалить контейнер
docker rm myapp

# Удалить все остановленные контейнеры
docker container prune

# Логи контейнера
docker logs myapp
docker logs -f myapp  # follow mode

# Выполнить команду в контейнере
docker exec -it myapp /bin/sh

# Скопировать файл из контейнера
docker cp myapp:/app/config.json ./config.json

# Скопировать файл в контейнер
docker cp ./config.json myapp:/app/config.json

# Статистика контейнера
docker stats myapp
```

### Docker Compose

```bash
# Запустить все сервисы
docker-compose up

# Запустить в фоновом режиме
docker-compose up -d

# Остановить все сервисы
docker-compose down

# Остановить и удалить volumes
docker-compose down -v

# Пересобрать образы
docker-compose build

# Пересобрать без кэша
docker-compose build --no-cache

# Логи всех сервисов
docker-compose logs -f

# Логи конкретного сервиса
docker-compose logs -f app

# Масштабирование сервиса
docker-compose up -d --scale app=3

# Список сервисов
docker-compose ps

# Выполнить команду в сервисе
docker-compose exec app /bin/sh
```

---

## Оптимальные флаги компиляции Go

```bash
# Минимальный размер
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -ldflags="-s -w" -trimpath -o app

# С версионированием
CGO_ENABLED=0 go build \
  -ldflags="-s -w -X main.Version=1.0.0 -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  -trimpath \
  -o app

# Для максимальной оптимизации
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build \
  -a \
  -installsuffix cgo \
  -ldflags="-s -w -extldflags '-static'" \
  -trimpath \
  -tags netgo \
  -o app
```

### Объяснение флагов

- `CGO_ENABLED=0` - отключает CGO, создает статический бинарник
- `GOOS=linux` - целевая ОС
- `GOARCH=amd64` - целевая архитектура
- `-ldflags="-s -w"` - удаляет debug info и symbol table
- `-trimpath` - удаляет file paths из бинарника
- `-a` - пересобрать все пакеты
- `-installsuffix cgo` - суффикс для разделения output
- `-extldflags '-static'` - статическая линковка
- `-tags netgo` - использовать pure Go network resolver

---

## Шаблоны Dockerfile

### Базовый Multi-Stage

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /build
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o app ./cmd

FROM alpine:3.21
RUN apk --no-cache add ca-certificates
COPY --from=builder /build/app /app
USER nobody
CMD ["/app"]
```

### Distroless

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /build
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o app ./cmd

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /build/app /app
ENTRYPOINT ["/app"]
```

### Scratch

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /build
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build \
    -ldflags="-s -w -extldflags '-static'" \
    -tags netgo \
    -o app ./cmd

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/app /app
USER 65534:65534
ENTRYPOINT ["/app"]
```

---

## .dockerignore

```
# Git
.git
.gitignore

# IDE
.vscode
.idea
*.swp

# Документация
*.md
docs/

# Тесты
*_test.go
testdata/

# Артефакты
*.exe
bin/
dist/

# Docker
Dockerfile*
docker-compose*.yml
.dockerignore

# Env
.env
*.env

# Vendor
vendor/
```

---

## Полезные инструменты

### dive - анализ слоев

```bash
# Установка
brew install dive

# Использование
dive myapp:latest

# Клавиши в dive:
# Tab - переключение между панелями
# Ctrl+U/Ctrl+D - прокрутка вверх/вниз
# Space - collapse/expand директорию
# Ctrl+C - выход
```

### trivy - сканирование безопасности

```bash
# Установка
brew install trivy

# Сканирование образа
trivy image myapp:latest

# Только критичные уязвимости
trivy image --severity CRITICAL,HIGH myapp:latest

# Сохранить результат в JSON
trivy image -f json -o results.json myapp:latest

# Сканирование Dockerfile
trivy config Dockerfile
```

### hadolint - линтинг Dockerfile

```bash
# Установка
brew install hadolint

# Проверка Dockerfile
hadolint Dockerfile

# С конкретным форматом вывода
hadolint --format json Dockerfile

# Игнорировать правила
hadolint --ignore DL3006 --ignore DL3018 Dockerfile
```

---

## Best Practices Checklist

### Оптимизация размера

- [ ] Используйте multi-stage build
- [ ] Выберите минимальный базовый образ (alpine/distroless/scratch)
- [ ] Компилируйте с флагами `-ldflags="-s -w"`
- [ ] Используйте .dockerignore
- [ ] Объединяйте RUN команды
- [ ] Удаляйте ненужные файлы в том же слое
- [ ] Рассмотрите UPX сжатие

### Кэширование

- [ ] Копируйте go.mod/go.sum перед исходным кодом
- [ ] Размещайте часто меняющиеся слои в конце
- [ ] Используйте BuildKit для параллельных сборок
- [ ] Настройте cache mount для go mod download

### Безопасность

- [ ] Не запускайте от root (USER nobody/65534)
- [ ] Не включайте секреты в образ
- [ ] Используйте конкретные версии базовых образов
- [ ] Сканируйте на уязвимости (trivy)
- [ ] Подписывайте образы
- [ ] Минимизируйте поверхность атаки

### Production Ready

- [ ] Добавьте HEALTHCHECK
- [ ] Добавьте метаданные (LABEL)
- [ ] Версионируйте образы (не только :latest)
- [ ] Документируйте EXPOSE порты
- [ ] Используйте ENTRYPOINT + CMD
- [ ] Тестируйте образы перед deploy

---

## Переменные окружения для Go приложений

```dockerfile
# В Dockerfile
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Runtime переменные
ENV TZ=UTC \
    APP_ENV=production \
    LOG_LEVEL=info
```

```bash
# При запуске
docker run \
  -e DB_HOST=postgres \
  -e DB_PORT=5432 \
  -e DB_USER=app \
  -e DB_PASSWORD=secret \
  -e GOMAXPROCS=4 \
  myapp:latest
```

---

## Debugging

### Вход в контейнер

```bash
# Alpine-based
docker run -it myapp:alpine /bin/sh

# Distroless :debug
docker run -it myapp:distroless-debug /busybox/sh

# Для running контейнера
docker exec -it myapp /bin/sh
```

### Просмотр логов

```bash
# Все логи
docker logs myapp

# Последние N строк
docker logs --tail 100 myapp

# Follow mode
docker logs -f myapp

# С timestamp
docker logs -t myapp
```

### Инспекция

```bash
# Полная информация
docker inspect myapp

# Конкретное поле (с jq)
docker inspect myapp | jq '.[0].Config.Env'

# IP адрес
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' myapp

# Используемые порты
docker inspect -f '{{.NetworkSettings.Ports}}' myapp
```

---

## Registry Operations

### Docker Hub

```bash
# Login
docker login

# Tag
docker tag myapp:latest username/myapp:1.0.0

# Push
docker push username/myapp:1.0.0

# Pull
docker pull username/myapp:1.0.0
```

### Приватный Registry

```bash
# Login
docker login registry.example.com

# Tag
docker tag myapp:latest registry.example.com/myapp:1.0.0

# Push
docker push registry.example.com/myapp:1.0.0
```

### GitHub Container Registry

```bash
# Login
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin

# Tag
docker tag myapp:latest ghcr.io/username/myapp:1.0.0

# Push
docker push ghcr.io/username/myapp:1.0.0
```

---

## Quick Reference

### Размеры базовых образов

| Образ | Размер | Shell | Package Manager |
|-------|--------|-------|-----------------|
| `scratch` | 0 MB | ❌ | ❌ |
| `busybox` | 1-2 MB | ✅ | ❌ |
| `alpine` | 5-7 MB | ✅ | apk |
| `distroless/static` | 2-3 MB | ❌ | ❌ |
| `distroless/base` | 20 MB | ❌ | ❌ |
| `debian:slim` | 80 MB | ✅ | apt |
| `ubuntu` | 100 MB | ✅ | apt |

### Типичные размеры Go приложений

| Подход | Размер образа |
|--------|---------------|
| golang:1.24 (без multi-stage) | ~800 MB |
| golang:1.24-alpine (без multi-stage) | ~350 MB |
| alpine + multi-stage | ~20 MB |
| distroless + multi-stage | ~15 MB |
| scratch + multi-stage | ~5-10 MB |

---

## Troubleshooting

### Проблема: Образ слишком большой

**Решение**:
1. Используйте multi-stage build
2. Проверьте с помощью `dive`
3. Добавьте .dockerignore
4. Используйте Alpine/Distroless/Scratch

### Проблема: Медленная сборка

**Решение**:
1. Оптимизируйте порядок COPY команд
2. Используйте BuildKit
3. Включите experimental features
4. Используйте cache mount

### Проблема: "standard_init_linux.go: exec user process caused: no such file or directory"

**Причина**: Динамическая линковка, отсутствие библиотек

**Решение**:
```bash
CGO_ENABLED=0 go build ...
```

### Проблема: "x509: certificate signed by unknown authority"

**Решение**:
```dockerfile
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
```

---

**Совет**: Добавьте эту шпаргалку в закладки для быстрого доступа!
