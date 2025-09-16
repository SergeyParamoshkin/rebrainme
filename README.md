# Rate Limiting и Retry паттерны в Go - Материалы для вебинара

## 📚 О чём этот проект

Этот репозиторий содержит демонстрационные материалы для вебинара по реализации паттернов устойчивости (resilience patterns) в Go. Рассматриваются ключевые механизмы для построения надёжных микросервисов: rate limiting, retry с экспоненциальным backoff, circuit breaker, retry budget и правильная работа с context.

## 🎯 Основные темы вебинара

### 1. Rate Limiting
**Зачем нужен:** Защита сервиса от перегрузки, распределение ресурсов между клиентами.

#### Token Bucket Algorithm
- **Принцип работы:** Корзина с токенами, которые пополняются с постоянной скоростью
- **Преимущества:** Позволяет обрабатывать всплески трафика
- **Применение:** API rate limiting, защита от DDoS

```go
// Пример использования
limiter := ratelimit.NewTokenBucket(10, 100) // 10 токенов/сек, ёмкость 100
if limiter.Allow() {
    // Обработка запроса
}
```

#### Sliding Window Log
- **Принцип работы:** Хранение времени каждого запроса в окне
- **Преимущества:** Точный подсчёт запросов
- **Применение:** Точный rate limiting для критичных API

### 2. Retry Strategies

#### Exponential Backoff с Jitter
- **Зачем:** Избежание "thundering herd" проблемы
- **Формула:** `delay = min(initialDelay * (multiplier ^ attempt), maxDelay) + random_jitter`

```go
config := &retry.Config{
    MaxRetries:   5,
    InitialDelay: 100 * time.Millisecond,
    MaxDelay:     10 * time.Second,
    Multiplier:   2.0,
    Jitter:       true,
}
```

### 3. Context в Go

#### Управление таймаутами
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Используем контекст для контроля времени выполнения
result, err := retrier.Do(ctx, someOperation)
```

#### Передача метаданных
```go
ctx = context.WithValue(ctx, "requestID", "abc123")
```

### 4. Retry Budget

**Концепция:** Ограничение количества повторных попыток для защиты downstream сервисов.

```go
budget := retry.NewBudget(
    0.5,              // 50% успешных запросов минимум
    10*time.Second,   // Окно наблюдения
    10,               // Минимум запросов для статистики
)
```

### 5. Circuit Breaker Pattern

**Состояния:**
- **Closed:** Нормальная работа
- **Open:** Блокировка запросов после серии ошибок
- **Half-Open:** Тестирование восстановления сервиса

```go
cb := NewCircuitBreaker(5, 10*time.Second) // 5 ошибок, 10 сек таймаут
result, err := cb.Execute(func() (interface{}, error) {
    return callExternalService()
})
```

## 🚀 Быстрый старт

### Установка и запуск

```bash
# Клонирование репозитория
git clone https://github.com/SergeyParamoshkin/rebrainme.git
cd rebrainme

# Установка зависимостей
go mod download

# Запуск сервера
go run cmd/main.go
```

Сервер запустится на порту 8080.

### Запуск тестов

```bash
# Все тесты
go test ./...

# С покрытием
go test -cover ./...

# Конкретный пакет с подробным выводом
go test -v ./internal/retry/

# Бенчмарки
go test -bench=. ./internal/ratelimit/
```

## 📡 API Endpoints

### Базовые эндпоинты

#### 1. Simple Handler
Простой эндпоинт для тестирования rate limiting.

```bash
curl http://localhost:8080/api/simple
```

**Ответ:**
```json
{
  "message": "Simple response",
  "timestamp": 1699123456,
  "path": "/api/simple"
}
```

#### 2. Flaky Handler
Имитация нестабильного сервиса (50% failure rate).

```bash
curl http://localhost:8080/api/flaky
```

**Возможные ответы:**
- Success (200): `{"message": "Success after potential failure", "lucky": true}`
- Failure (503): `Random failure`

### Демонстрация паттернов

#### 3. Retry Demo
Демонстрация retry механизма с экспоненциальным backoff.

```bash
time curl http://localhost:8080/api/retry-demo
```

**Ответ после успешных повторов:**
```json
{
  "message": "Success after retries",
  "attempts": 3
}
```

#### 4. Context Demo
Демонстрация управления таймаутами через context.

```bash
# С коротким таймаутом клиента
curl -m 1 http://localhost:8080/api/context-demo

# Обычный запрос
curl http://localhost:8080/api/context-demo
```

#### 5. Circuit Breaker Demo
Демонстрация паттерна Circuit Breaker.

```bash
# Генерируем ошибки для открытия circuit breaker
for i in {1..10}; do
  curl http://localhost:8080/api/circuit-breaker
  sleep 0.1
done
```

#### 6. Retry Budget Demo
Демонстрация работы retry budget.

```bash
curl http://localhost:8080/api/budget-demo
```

**Ответ:**
```json
{
  "success": true,
  "budget_stats": {
    "success_rate": 0.6,
    "total_requests": 15
  }
}
```

### Мониторинг

#### Health Check
```bash
curl http://localhost:8080/health
```

#### Prometheus Metrics
```bash
curl http://localhost:8080/metrics | grep -E "(rate_limited|retry_attempts|circuit_breaker)"
```

## 🧪 Примеры тестирования

### Нагрузочное тестирование Rate Limiting

#### Token Bucket Test
```bash
# Отправляем 20 запросов параллельно (лимит 10/сек)
for i in {1..20}; do
  curl -s http://localhost:8080/api/simple &
done
wait

# Ожидаемый результат: ~10 успешных, ~10 с ошибкой 429
```

#### Sliding Window Test
```bash
# Отправляем 110 запросов последовательно (лимит 100/мин)
for i in {1..110}; do
  echo -n "$i: "
  curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8080/api/simple
done
```

### Демонстрация Retry с мониторингом

```bash
# В одном терминале - мониторинг метрик
watch -n 1 'curl -s http://localhost:8080/metrics | grep retry_attempts'

# В другом - генерация запросов
while true; do
  curl -s http://localhost:8080/api/retry-demo
  sleep 2
done
```

### Circuit Breaker в действии

```bash
# Скрипт для демонстрации перехода состояний
for i in {1..30}; do
  echo "Request $i:"
  curl -s http://localhost:8080/api/circuit-breaker | jq .
  echo "---"
  sleep 0.5
done
```

## 📊 Метрики Prometheus

### Основные метрики

| Метрика | Описание | Тип |
|---------|----------|-----|
| `http_requests_total` | Общее количество HTTP запросов | Counter |
| `http_request_duration_seconds` | Длительность обработки запросов | Histogram |
| `rate_limited_requests_total` | Количество заблокированных запросов | Counter |
| `retry_attempts_total` | Количество повторных попыток | Counter |
| `retry_budget_exhausted_total` | Исчерпание retry budget | Counter |
| `circuit_breaker_state` | Состояние circuit breaker (0=closed, 1=open, 2=half-open) | Gauge |
| `token_bucket_available` | Доступные токены в bucket | Gauge |
| `sliding_window_current_requests` | Текущее количество запросов в окне | Gauge |

### Примеры PromQL запросов

```promql
# Rate of requests per second
rate(http_requests_total[1m])

# Процент успешных запросов
sum(rate(http_requests_total{status="200"}[5m])) / sum(rate(http_requests_total[5m]))

# Среднее время ответа
rate(http_request_duration_seconds_sum[5m]) / rate(http_request_duration_seconds_count[5m])

# Процент rate limited запросов
sum(rate(rate_limited_requests_total[5m])) / sum(rate(http_requests_total[5m]))
```

## 🔧 Конфигурация

### Rate Limiter
```yaml
Token Bucket:
  rate: 10 req/sec         # Скорость пополнения токенов
  capacity: 100 tokens      # Максимальная ёмкость корзины

Sliding Window:
  window: 1 minute          # Размер окна
  limit: 100 requests       # Максимум запросов в окне
```

### Retry Configuration
```yaml
Retry:
  max_retries: 5            # Максимальное количество попыток
  initial_delay: 100ms      # Начальная задержка
  max_delay: 10s            # Максимальная задержка
  multiplier: 2.0           # Множитель для экспоненциального роста
  jitter: true              # Добавление случайного разброса
```

### Circuit Breaker
```yaml
Circuit Breaker:
  max_failures: 5           # Порог ошибок для открытия
  reset_timeout: 10s        # Время до перехода в half-open
  half_open_requests: 1     # Количество тестовых запросов
```

### Retry Budget
```yaml
Retry Budget:
  success_threshold: 0.5    # Минимальный процент успешных запросов
  window_size: 10s          # Окно для подсчёта статистики
  min_requests: 10          # Минимум запросов для активации
```

## 📖 Дополнительные материалы

### Статьи и документация
- [Google SRE Book - Handling Overload](https://sre.google/sre-book/handling-overload/)
- [AWS Architecture Blog - Exponential Backoff and Jitter](https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/)
- [Martin Fowler - Circuit Breaker](https://martinfowler.com/bliki/CircuitBreaker.html)

### Полезные библиотеки
- [golang.org/x/time/rate](https://pkg.go.dev/golang.org/x/time/rate) - Официальный rate limiter
- [github.com/sony/gobreaker](https://github.com/sony/gobreaker) - Circuit breaker
- [github.com/cenkalti/backoff](https://github.com/cenkalti/backoff) - Backoff алгоритмы

### Инструменты для тестирования
- [Apache Bench (ab)](https://httpd.apache.org/docs/2.4/programs/ab.html)
- [wrk](https://github.com/wg/wrk) - Modern HTTP benchmarking tool
- [hey](https://github.com/rakyll/hey) - HTTP load generator
- [Grafana k6](https://k6.io/) - Load testing tool

## 🤝 Вопросы для обсуждения на вебинаре

1. **Как выбрать правильные параметры для rate limiting?**
   - Анализ бизнес-требований
   - Мониторинг текущей нагрузки
   - A/B тестирование

2. **Когда использовать retry, а когда circuit breaker?**
   - Retry для временных сбоев
   - Circuit breaker для защиты от каскадных отказов

3. **Как правильно настроить retry budget?**
   - Баланс между доступностью и нагрузкой
   - Учёт SLA downstream сервисов

4. **Распределённый rate limiting**
   - Использование Redis для синхронизации
   - Eventual consistency vs Strong consistency

## 📖 Глоссарий терминов

Подробный глоссарий всех терминов и концепций доступен в файле [GLOSSARY.md](./GLOSSARY.md)

### Ключевые термины:
- **Rate Limiting** - ограничение количества запросов
- **Token Bucket** - алгоритм с корзиной токенов
- **Exponential Backoff** - экспоненциальное увеличение задержки
- **Jitter** - случайный разброс для предотвращения синхронизации
- **Thundering Herd** - проблема одновременного доступа множества клиентов
- **Circuit Breaker** - защита от каскадных сбоев
- **Retry Budget** - ограничение повторов на основе успешности
- **429 Too Many Requests** - HTTP код для rate limiting
- **503 Service Unavailable** - HTTP код временной недоступности

## 👨‍💻 Автор

Sergey Paramoshkin

---