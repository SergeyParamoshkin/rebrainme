# Docker Best Practices для Golang приложений

## Содержание

1. [Введение](#введение)
2. [Лучшие практики](#лучшие-практики)
3. [Варианты Dockerfile](#варианты-dockerfile)
4. [Сравнение подходов](#сравнение-подходов)
5. [Практические примеры](#практические-примеры)
6. [Инструменты и утилиты](#инструменты-и-утилиты)
7. [Дополнительные материалы](#дополнительные-материалы)

---

## Введение

Этот репозиторий содержит материалы для вебинара по сборке Docker образов для Golang приложений с применением лучших практик.

### Цели вебинара

- Изучить различные подходы к созданию Docker образов для Go
- Оптимизировать размер финальных образов
- Повысить безопасность приложений в контейнерах
- Ускорить процесс сборки за счет правильного использования кэширования
- Научиться использовать современные инструменты для работы с Docker

### Предварительные требования

- Docker 20.10+
- Docker Compose 2.0+
- Базовые знания Docker и Golang
- (Опционально) make, jq, trivy, dive

---

## Лучшие практики

### 1. Multi-Stage Build

**Проблема**: Большой размер финального образа из-за включения build-инструментов.

**Решение**: Использование multi-stage build для разделения этапов сборки и runtime.

```dockerfile
# Этап сборки
FROM golang:1.24-alpine AS builder
WORKDIR /build
COPY . .
RUN go build -o app ./cmd

# Финальный этап
FROM alpine:3.21
COPY --from=builder /build/app /app
CMD ["/app"]
```

**Преимущества**:
- ✅ Уменьшение размера образа в 10-50 раз
- ✅ Отсутствие компилятора и исходного кода в production образе
- ✅ Улучшенная безопасность

### 2. Правильное кэширование слоев

**Принцип**: Docker кэширует слои, если они не изменились.

**Плохой подход**:
```dockerfile
COPY . .
RUN go mod download
RUN go build -o app ./cmd
```

**Хороший подход**:
```dockerfile
# Сначала копируем только файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Потом копируем исходный код
COPY . .
RUN go build -o app ./cmd
```

**Результат**: Изменения в коде не инвалидируют кэш зависимостей.

### 3. Оптимизация размера бинарника

**Флаги компиляции**:

```dockerfile
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags="-s -w" \
    -trimpath \
    -o app ./cmd
```

**Что делают флаги**:
- `CGO_ENABLED=0` - отключает CGO, создает полностью статический бинарник
- `-ldflags="-s -w"` - удаляет отладочную информацию и таблицу символов
- `-trimpath` - удаляет пути к файлам из бинарника
- Результат: уменьшение размера на 30-50%

### 4. Выбор базового образа

| Образ | Размер | Безопасность | Использование |
|-------|--------|-------------|---------------|
| `alpine` | ~5 MB | Средняя | Development, легкие приложения |
| `distroless` | ~2 MB | Высокая | Production, без отладки |
| `scratch` | 0 MB | Максимальная | Production, статические бинарники |
| `debian/ubuntu` | ~100 MB | Низкая | Требуются системные утилиты |

### 5. Безопасность

#### 5.1 Непривилегированный пользователь

**Плохо**:
```dockerfile
FROM alpine
COPY app /app
CMD ["/app"]  # Запуск от root
```

**Хорошо**:
```dockerfile
FROM alpine
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser
USER appuser
COPY --chown=appuser:appuser app /app
CMD ["/app"]
```

#### 5.2 Не включайте секреты в образ

```dockerfile
# ❌ ПЛОХО
COPY .env /app/.env

# ✅ ХОРОШО - используйте ENV переменные или secrets
```

#### 5.3 Используйте .dockerignore

```
.git
*.md
test/
.env
```

### 6. Health Checks

```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1
```

### 7. Метаданные и версионирование

```dockerfile
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

LABEL maintainer="your@email.com" \
      version="${VERSION}" \
      description="My application"

RUN go build \
    -ldflags="-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}" \
    -o app ./cmd
```

### 8. Минимизация количества слоев

**Плохо**:
```dockerfile
RUN apk add git
RUN apk add ca-certificates
RUN apk add tzdata
```

**Хорошо**:
```dockerfile
RUN apk add --no-cache git ca-certificates tzdata && \
    rm -rf /var/cache/apk/*
```

---

## Варианты Dockerfile

### 1. Basic (`Dockerfile.basic`)

**Описание**: Базовый multi-stage build с Alpine.

**Размер**: ~20-25 MB

**Использование**: Development, быстрый старт

**Особенности**:
- Multi-stage build
- Alpine Linux в качестве базы
- Непривилегированный пользователь
- Health check

### 2. Optimized (`Dockerfile.optimized`)

**Описание**: Оптимизированная версия с максимальным кэшированием.

**Размер**: ~20 MB

**Использование**: Production, частые пересборки

**Особенности**:
- Продвинутое кэширование слоев
- UPX сжатие бинарника
- Метаданные и labels
- ENTRYPOINT + CMD паттерн

### 3. Distroless (`Dockerfile.distroless`)

**Описание**: Использование Google Distroless образа.

**Размер**: ~15 MB

**Использование**: Production, высокие требования к безопасности

**Особенности**:
- Нет package manager
- Нет shell
- Минимальная поверхность атаки
- Сложная отладка

**Ограничения**:
- Невозможно использовать health checks (нет curl/wget)
- Нельзя зайти в контейнер через exec
- Требуется статическая компиляция

### 4. Scratch (`Dockerfile.scratch`)

**Описание**: Минимальный образ на базе scratch.

**Размер**: ~5-10 MB (только бинарник)

**Использование**: Production, максимальная оптимизация

**Особенности**:
- Абсолютный минимум
- Только бинарник + CA certificates
- Максимальная безопасность

**Требования**:
- Полностью статическая компиляция
- CGO_ENABLED=0
- Нельзя использовать системные библиотеки

### 5. Advanced Multi-Stage (`Dockerfile.multistage-advanced`)

**Описание**: Продвинутый вариант с тестами и линтингом.

**Размер**: ~20 MB

**Использование**: CI/CD pipeline, production

**Особенности**:
- Встроенные тесты
- Линтинг кода
- Версионирование
- Build arguments
- Полная автоматизация

---

## Сравнение подходов

### Размеры образов (примерные)

| Вариант | Размер | Слоев | Уменьшение |
|---------|--------|-------|------------|
| Без multi-stage (golang:1.24) | ~800 MB | 15 | baseline |
| Basic (alpine) | ~25 MB | 8 | -97% |
| Optimized (alpine + upx) | ~18 MB | 7 | -98% |
| Distroless | ~15 MB | 5 | -98.5% |
| Scratch | ~8 MB | 3 | -99% |

### Матрица выбора

| Критерий | Basic | Optimized | Distroless | Scratch | Advanced |
|----------|-------|-----------|------------|---------|----------|
| Размер | 🟡 | 🟢 | 🟢 | 🟢🟢 | 🟡 |
| Безопасность | 🟡 | 🟡 | 🟢 | 🟢🟢 | 🟢 |
| Отладка | 🟢 | 🟢 | 🔴 | 🔴🔴 | 🟢 |
| Скорость сборки | 🟢 | 🟡 | 🟢 | 🟢 | 🔴 |
| Production-ready | 🟡 | 🟢 | 🟢 | 🟢 | 🟢🟢 |

🟢🟢 - Отлично | 🟢 - Хорошо | 🟡 - Средне | 🔴 - Плохо | 🔴🔴 - Очень плохо

---

## Практические примеры

### Быстрый старт

```bash
# 1. Собрать базовый вариант
cd webinar-materials
make build-basic

# 2. Запустить контейнер
make run-basic

# 3. Проверить работу
curl http://localhost:8081/health
```

### Сборка всех вариантов

```bash
# Используя Makefile
make build-all

# Или напрямую скрипт
./scripts/build-all.sh
```

### Сравнение размеров

```bash
# Детальное сравнение
make compare

# Или
./scripts/compare-sizes.sh
```

### Запуск demo окружения

```bash
# Запустить все варианты через docker-compose
make run-demo

# Доступные сервисы:
# - http://localhost:8081 - basic
# - http://localhost:8082 - optimized
# - http://localhost:8083 - distroless
# - http://localhost:8084 - scratch
# - http://localhost:8085 - advanced
# - http://localhost:8080 - Adminer
# - http://localhost:9000 - Portainer
```

### Анализ безопасности

```bash
# Установить trivy
brew install trivy  # MacOS
# или см. https://github.com/aquasecurity/trivy

# Запустить сканирование
make security-scan
```

### Анализ слоев

```bash
# Установить dive
brew install dive  # MacOS

# Анализ образа
make dive-basic
# или
dive golang-demo:basic
```

---

## Инструменты и утилиты

### Обязательные

1. **Docker Desktop** / **Docker Engine**
   - Версия: 20.10+
   - https://www.docker.com/

2. **Docker Compose**
   - Версия: 2.0+
   - Обычно входит в Docker Desktop

### Рекомендуемые

1. **dive** - анализ слоев Docker образов
   ```bash
   brew install dive
   # или
   docker pull wagoodman/dive
   ```

2. **trivy** - сканирование уязвимостей
   ```bash
   brew install trivy
   ```

3. **hadolint** - линтер для Dockerfile
   ```bash
   brew install hadolint
   ```

4. **jq** - обработка JSON
   ```bash
   brew install jq
   ```

### Полезные

1. **ctop** - мониторинг контейнеров
   ```bash
   brew install ctop
   ```

2. **lazydocker** - TUI для Docker
   ```bash
   brew install lazydocker
   ```

---

## Команды Makefile

```bash
make help              # Список всех команд
make build-all         # Собрать все варианты
make compare           # Сравнить размеры
make run-demo          # Запустить demo окружение
make clean-all         # Очистить все
make security-scan     # Сканирование безопасности
make stats             # Статистика образов
```

---

## Дополнительные материалы

### Официальная документация

- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [Go Official Images](https://hub.docker.com/_/golang)
- [Distroless Images](https://github.com/GoogleContainerTools/distroless)

### Статьи и руководства

- [Building Minimal Docker Images for Go Applications](https://blog.golang.org/docker)
- [Dockerfile Security Best Practices](https://docs.docker.com/develop/security-best-practices/)
- [Multi-Stage Builds](https://docs.docker.com/build/building/multi-stage/)

### Инструменты

- [Dive - анализ образов](https://github.com/wagoodman/dive)
- [Trivy - сканер безопасности](https://github.com/aquasecurity/trivy)
- [Hadolint - Dockerfile linter](https://github.com/hadolint/hadolint)

---

## Структура проекта

```
webinar-materials/
├── README.md                    # Эта документация
├── Makefile                     # Команды для работы с проектом
├── .dockerignore                # Исключения для Docker build
├── docker-compose.demo.yml      # Demo окружение
├── dockerfiles/                 # Варианты Dockerfile
│   ├── Dockerfile.basic
│   ├── Dockerfile.optimized
│   ├── Dockerfile.distroless
│   ├── Dockerfile.scratch
│   └── Dockerfile.multistage-advanced
└── scripts/                     # Утилиты
    ├── build-all.sh
    └── compare-sizes.sh
```

---

## Практические задания

### Задание 1: Сравнение размеров

1. Соберите все варианты Dockerfile
2. Сравните их размеры
3. Изучите слои с помощью `dive`
4. Запишите результаты

### Задание 2: Оптимизация существующего Dockerfile

1. Возьмите существующий Dockerfile из вашего проекта
2. Примените техники из вебинара
3. Сравните размеры до и после
4. Проверьте функциональность

### Задание 3: Безопасность

1. Просканируйте образы с помощью trivy
2. Найдите уязвимости
3. Примените исправления
4. Повторите сканирование

### Задание 4: CI/CD интеграция

1. Создайте GitHub Actions workflow
2. Добавьте автоматическую сборку
3. Добавьте тесты и линтинг
4. Настройте push в registry

---

## Контрольный список (Checklist)

При создании production Dockerfile проверьте:

- [ ] Используется multi-stage build
- [ ] Зависимости кэшируются отдельным слоем
- [ ] Бинарник собирается с флагами оптимизации
- [ ] Используется минимальный базовый образ
- [ ] Приложение запускается от непривилегированного пользователя
- [ ] Присутствует .dockerignore
- [ ] Нет секретов в образе
- [ ] Добавлен HEALTHCHECK (если применимо)
- [ ] Добавлены метаданные (LABEL)
- [ ] Проведено сканирование безопасности
- [ ] Протестирована работа приложения
- [ ] Документированы особенности запуска

---

## FAQ

### Q: Какой вариант Dockerfile выбрать для production?

**A**: Зависит от требований:
- Нужна отладка → `Optimized` или `Advanced`
- Максимальная безопасность → `Distroless` или `Scratch`
- Баланс размер/удобство → `Optimized`

### Q: Когда использовать scratch образ?

**A**: Когда:
- Приложение не использует CGO
- Не нужны системные утилиты
- Критичен размер образа
- Максимальные требования к безопасности

### Q: Как отлаживать distroless/scratch контейнеры?

**A**:
1. Используйте `:debug` тег для distroless (содержит busybox)
2. Создайте отдельный debug Dockerfile на базе Alpine
3. Используйте kubectl debug (в Kubernetes)
4. Логируйте в stdout/stderr

### Q: Стоит ли использовать UPX сжатие?

**A**: Зависит:
- ✅ Плюсы: меньше размер, быстрее pull
- ❌ Минусы: медленнее старт, сложнее отладка
- Рекомендация: тестируйте на вашем приложении

---

## Заключение

Правильная сборка Docker образов для Golang приложений позволяет:

1. **Уменьшить размер** образов с ~800MB до ~8MB (99% reduction)
2. **Повысить безопасность** за счет минимизации поверхности атаки
3. **Ускорить deployment** благодаря меньшему размеру
4. **Снизить costs** на хранение и трафик
5. **Улучшить производительность** CI/CD pipeline

Применяйте эти практики в своих проектах и делитесь опытом с командой!

---

**Автор**: Материалы для вебинара по Docker Best Practices
**Дата**: 2025
**Версия**: 1.0.0
