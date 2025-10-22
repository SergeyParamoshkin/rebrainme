# Docker Best Practices для Golang

## Презентация вебинара

---

## Слайд 1: Введение

### Тема вебинара
**Сборка Docker образов для Golang приложений: Best Practices**

### О чем поговорим
- Multi-stage builds
- Оптимизация размера образов
- Безопасность
- Кэширование и ускорение сборки
- Практические примеры

---

## Слайд 2: Проблема

### Типичный Dockerfile (без оптимизации)

```dockerfile
FROM golang:1.24
WORKDIR /app
COPY . .
RUN go build -o app ./cmd
CMD ["./app"]
```

### Результат
- **Размер образа**: ~800 MB
- **Время сборки**: медленно при каждом изменении
- **Безопасность**: низкая (много ненужных пакетов)
- **Production-ready**: ❌

---

## Слайд 3: Multi-Stage Build - Решение

### Идея
Разделить сборку и runtime на разные этапы

```dockerfile
# Этап 1: Сборка
FROM golang:1.24-alpine3.21 AS builder
WORKDIR /build
COPY . .
RUN go build -o app ./cmd

# Этап 2: Runtime
FROM alpine:3.21
COPY --from=builder /build/app /app
CMD ["/app"]
```

### Результат
- **Размер**: ~800 MB → ~25 MB (97% уменьшение!)
- В финальном образе нет Go compiler
- Только необходимое для запуска

---

## Слайд 4: Оптимизация кэширования

### Проблема
```dockerfile
COPY . .              # ← Копируем ВСЁ
RUN go mod download   # ← Кэш инвалидируется при любом изменении
```

### Решение
```dockerfile
# Сначала зависимости
COPY go.mod go.sum ./
RUN go mod download    # ← Кэш сохраняется, если зависимости не менялись

# Потом код
COPY . .
RUN go build ...
```

### Результат
- **Ускорение пересборки** в 10-100 раз
- Зависимости скачиваются только при их изменении

---

## Слайд 5: Флаги компиляции Go

### Базовая сборка
```bash
go build -o app ./cmd
```
**Размер**: ~10 MB

### Оптимизированная сборка
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -ldflags="-s -w" -trimpath -o app ./cmd
```
**Размер**: ~6 MB (40% меньше!)

### Что делают флаги
- `CGO_ENABLED=0` - статический бинарник
- `-ldflags="-s -w"` - удаление debug info
- `-trimpath` - удаление file paths

---

## Слайд 6: Выбор базового образа

| Образ | Размер | Особенности | Использование |
|-------|--------|-------------|---------------|
| **golang:1.24** | 800 MB | Полный набор инструментов | Development |
| **alpine** | 25 MB | Минимальный Linux + shell | Production (standard) |
| **distroless** | 15 MB | Нет shell, package manager | Production (secure) |
| **scratch** | 8 MB | Только бинарник | Production (minimal) |

### Рекомендация
- Development: Alpine
- Production (стандарт): Alpine или Distroless
- Production (максимум безопасности): Distroless или Scratch

---

## Слайд 7: Scratch образ

### Что такое scratch?
- Абсолютно пустой базовый образ
- Размер: 0 байт
- Только ваш бинарник

### Dockerfile
```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build \
    -ldflags="-s -w -extldflags '-static'" \
    -o app ./cmd

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/app /app
ENTRYPOINT ["/app"]
```

### Важно
- ⚠️ Нужен **полностью статический** бинарник
- ⚠️ Нельзя использовать CGO
- ⚠️ Нет shell для отладки

---

## Слайд 8: Distroless

### Google Distroless
- Минимальные образы от Google
- Содержат только runtime зависимости
- Нет package manager, shell

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o app ./cmd

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /build/app /app
ENTRYPOINT ["/app"]
```

### Преимущества
- ✅ Меньше уязвимостей
- ✅ Меньше поверхность атаки
- ✅ Хороший баланс размер/функциональность

---

## Слайд 9: Безопасность

### 1. Непривилегированный пользователь

❌ **Плохо**:
```dockerfile
FROM alpine
COPY app /app
CMD ["/app"]  # Запуск от root!
```

✅ **Хорошо**:
```dockerfile
FROM alpine
RUN adduser -D -u 1000 appuser
USER appuser
COPY --chown=appuser:appuser app /app
CMD ["/app"]
```

### 2. Не включайте секреты
```dockerfile
# ❌ НИКОГДА
COPY .env /app/.env
COPY secrets.json /app/

# ✅ Используйте ENV или Docker secrets
```

### 3. .dockerignore
```
.git
*.md
.env
test/
```

---

## Слайд 10: .dockerignore

### Зачем нужен?
1. **Ускорение сборки** - меньше контекст
2. **Безопасность** - не попадут секреты
3. **Размер** - меньше лишних файлов

### Пример
```
# Git
.git
.gitignore

# IDE
.vscode
.idea

# Тесты
*_test.go
testdata/

# Секреты
.env
*.key
*.pem

# Docker
Dockerfile*
docker-compose*.yml
```

---

## Слайд 11: Сравнение размеров

### Наши результаты

| Подход | Размер | Уменьшение |
|--------|--------|------------|
| Без multi-stage | 800 MB | baseline |
| Alpine + multi-stage | 25 MB | **-97%** |
| Optimized | 18 MB | **-98%** |
| Distroless | 15 MB | **-98.5%** |
| Scratch | 8 MB | **-99%** |

### Экономия
- **Меньше трафика** при pull
- **Быстрее deployment**
- **Меньше затраты** на storage
- **Быстрее старт** контейнеров

---

## Слайд 12: Health Checks

### Зачем?
- Kubernetes / Docker Swarm знают, когда приложение готово
- Автоматический restart при проблемах
- Правильный load balancing

### Dockerfile
```dockerfile
HEALTHCHECK --interval=30s --timeout=3s \
    --start-period=5s --retries=3 \
    CMD wget --spider http://localhost:8080/health || exit 1
```

### В приложении
```go
http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
})
```

---

## Слайд 13: Метаданные и версионирование

### Labels
```dockerfile
LABEL maintainer="team@example.com" \
      version="1.0.0" \
      description="My awesome app"
```

### Build-time версионирование
```dockerfile
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

RUN go build \
    -ldflags="-X main.Version=${VERSION} \
              -X main.BuildTime=${BUILD_TIME} \
              -X main.GitCommit=${GIT_COMMIT}" \
    -o app ./cmd
```

### Использование
```bash
docker build \
  --build-arg VERSION=1.0.0 \
  --build-arg BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
  --build-arg GIT_COMMIT=$(git rev-parse --short HEAD) \
  -t myapp:1.0.0 .
```

---

## Слайд 14: Инструменты

### dive - анализ слоев
```bash
brew install dive
dive myapp:latest
```
- Показывает каждый слой
- Помогает найти "тяжелые" места
- Efficiency score

### trivy - сканер безопасности
```bash
brew install trivy
trivy image myapp:latest
```
- Поиск уязвимостей
- CVE database
- Рекомендации по исправлению

### hadolint - линтер Dockerfile
```bash
brew install hadolint
hadolint Dockerfile
```
- Best practices проверка
- Советы по оптимизации

---

## Слайд 15: CI/CD интеграция

### GitHub Actions пример
```yaml
name: Docker Build

on:
  push:
    tags: ['v*']

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: .
          push: true
          tags: myapp:${{ github.ref_name }}
          cache-from: type=gha
          cache-to: type=gha,mode=max
```

### BuildKit кэш
- Кэш между сборками
- Ускорение CI в 10+ раз

---

## Слайд 16: Практический пример

### Было (800 MB, медленно, небезопасно)
```dockerfile
FROM golang:1.24
WORKDIR /app
COPY . .
RUN go build -o app ./cmd
CMD ["./app"]
```

### Стало (8 MB, быстро, безопасно)
```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /build
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o app ./cmd

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/app /app
USER 65534:65534
ENTRYPOINT ["/app"]
```

---

## Слайд 17: Checklist для Production

✅ Multi-stage build
✅ Минимальный базовый образ (alpine/distroless/scratch)
✅ Оптимизация флагов компиляции (`-ldflags="-s -w"`)
✅ .dockerignore файл
✅ Непривилегированный пользователь
✅ Нет секретов в образе
✅ HEALTHCHECK
✅ Метаданные (LABEL)
✅ Версионирование
✅ Security scan (trivy)
✅ Тестирование образа

---

## Слайд 18: Частые ошибки

### ❌ 1. Копировать все сразу
```dockerfile
COPY . .
RUN go mod download
```

### ❌ 2. Забыть .dockerignore
→ Медленная сборка, большой контекст

### ❌ 3. Запуск от root
→ Security проблемы

### ❌ 4. Использовать :latest в production
→ Непредсказуемое поведение

### ❌ 5. Включать секреты в образ
→ Утечка данных

---

## Слайд 19: Когда использовать какой образ?

### Development
```dockerfile
FROM golang:1.24-alpine
# Нужны debug утилиты, shell
```

### Staging
```dockerfile
FROM alpine:3.21
# Баланс между debugging и production
```

### Production (стандарт)
```dockerfile
FROM gcr.io/distroless/static:nonroot
# Безопасность + разумный размер
```

### Production (максимальная оптимизация)
```dockerfile
FROM scratch
# Только для статических бинарников
# Минимальный размер и максимальная безопасность
```

---

## Слайд 20: Демо

### Что покажем
1. Сборка разных вариантов Dockerfile
2. Сравнение размеров с помощью dive
3. Security scan с trivy
4. Запуск через docker-compose

### Команды
```bash
# Сборка всех вариантов
make build-all

# Сравнение
make compare

# Анализ
dive golang-demo:basic

# Security scan
trivy image golang-demo:basic

# Запуск demo
make run-demo
```

---

## Слайд 21: Резюме

### Ключевые выводы

1. **Multi-stage builds** - обязательно для production
2. **Выбор базового образа** - alpine/distroless/scratch
3. **Оптимизация сборки** - правильный порядок COPY
4. **Безопасность** - непривилегированный пользователь, сканирование
5. **Инструменты** - dive, trivy, hadolint
6. **CI/CD** - автоматизация и кэширование

### Результаты
- **99% уменьшение** размера (800MB → 8MB)
- **10-100x ускорение** пересборки
- **Повышение безопасности**
- **Production-ready** образы

---

## Слайд 22: Дополнительные ресурсы

### Документация
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [Go Official Images](https://hub.docker.com/_/golang)
- [Distroless Images](https://github.com/GoogleContainerTools/distroless)

### Инструменты
- [dive](https://github.com/wagoodman/dive) - анализ слоев
- [trivy](https://github.com/aquasecurity/trivy) - security scanning
- [hadolint](https://github.com/hadolint/hadolint) - Dockerfile linter

### Этот репозиторий
- Все примеры Dockerfile
- Скрипты для сборки и сравнения
- Подробная документация
- docker-compose для демо

---

## Слайд 23: Q&A

### Вопросы?

📧 Контакты для связи
📦 GitHub репозиторий с материалами
💬 Обсуждение в комментариях

### Спасибо за внимание!

---

## Слайд 24: Практические задания

### Для самостоятельной работы

1. **Задание 1**: Возьмите ваш проект и оптимизируйте Dockerfile
   - Цель: уменьшить размер минимум на 80%

2. **Задание 2**: Настройте CI/CD
   - GitHub Actions или GitLab CI
   - Автоматическая сборка и push в registry

3. **Задание 3**: Security audit
   - Просканируйте образы с trivy
   - Исправьте найденные проблемы

4. **Задание 4**: Мониторинг
   - Добавьте health checks
   - Настройте метрики

---

**Конец презентации**

Материалы доступны в репозитории: `webinar-materials/`
