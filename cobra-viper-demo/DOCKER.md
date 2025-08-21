# Docker инструкция для демо-приложения

## Быстрый старт

### 1. Запуск базы данных и сервисов

```bash
# Запустить PostgreSQL, Redis и Adminer
docker-compose up -d

# Проверить статус
docker-compose ps

# Посмотреть логи
docker-compose logs -f postgres
```

### 2. Доступ к сервисам

- **PostgreSQL**: `localhost:5432`
  - User: `demo_user`
  - Password: `demo_password`
  - Database: `demo_db`

- **Adminer** (Web UI для БД): http://localhost:8088
  - System: PostgreSQL
  - Server: postgres
  - Username: demo_user
  - Password: demo_password
  - Database: demo_db

- **Redis**: `localhost:6379`

## Использование с приложением

### Запуск приложения локально с Docker БД

```bash
# 1. Запустить инфраструктуру
docker-compose up -d

# 2. Дождаться готовности БД
docker-compose exec postgres pg_isready

# 3. Запустить приложение
go run main.go serve --verbose

# 4. Проверить подключение к БД
go run main.go db status
go run main.go db connect
```

### Запуск приложения в Docker

```bash
# Собрать образ приложения
docker build -t cobra-viper-demo:latest .

# Запустить всё вместе
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up

# Или запустить только приложение
docker run --rm \
  --network cobra-viper-demo_demo-network \
  -e DEMO_DATABASE_HOST=postgres \
  -e DEMO_DATABASE_PORT=5432 \
  -e DEMO_DATABASE_NAME=demo_db \
  -e DEMO_DATABASE_USER=demo_user \
  -e DEMO_DATABASE_PASSWORD=demo_password \
  -p 8080:8080 \
  cobra-viper-demo:latest serve --verbose
```

## Команды для работы с БД

### Через приложение

```bash
# Проверить статус
go run main.go db status

# Показать строку подключения
go run main.go db connect

# Выполнить миграции
go run main.go db migrate

# Создать backup
go run main.go db backup
```

### Через Docker

```bash
# Подключиться к PostgreSQL через psql
docker-compose exec postgres psql -U demo_user -d demo_db

# Выполнить SQL команду
docker-compose exec postgres psql -U demo_user -d demo_db -c "SELECT * FROM users;"

# Создать дамп БД
docker-compose exec postgres pg_dump -U demo_user demo_db > backup.sql

# Восстановить из дампа
docker-compose exec -T postgres psql -U demo_user demo_db < backup.sql

# Подключиться к Redis
docker-compose exec redis redis-cli

# Проверить Redis
docker-compose exec redis redis-cli ping
```

## Мониторинг и отладка

### Логи

```bash
# Все сервисы
docker-compose logs

# Конкретный сервис
docker-compose logs postgres
docker-compose logs -f postgres  # следить за логами

# Последние 100 строк
docker-compose logs --tail=100 postgres
```

### Статистика

```bash
# Использование ресурсов
docker stats

# Информация о контейнерах
docker-compose ps
docker inspect demo-postgres
```

### Подключение к контейнеру

```bash
# Bash в контейнере PostgreSQL
docker-compose exec postgres sh

# Выполнить команду
docker-compose exec postgres ls -la /var/lib/postgresql/data
```

## Production конфигурация

### Использование .env файла

```bash
# Создать .env.prod
cp .env.example .env.prod
# Отредактировать с production значениями

# Запустить с production конфигурацией
docker-compose --env-file .env.prod -f docker-compose.prod.yml up -d
```

### Переменные окружения для production

```bash
# .env.prod
DB_USER=production_user
DB_PASSWORD=strong_password_here
DB_NAME=production_db
DB_PORT=5432
APP_PORT=8080
REDIS_PASSWORD=redis_password_here
LOG_LEVEL=warn
```

## Управление данными

### Volumes

```bash
# Список volumes
docker volume ls

# Информация о volume
docker volume inspect cobra-viper-demo_postgres_data

# Создать backup volume
docker run --rm \
  -v cobra-viper-demo_postgres_data:/data \
  -v $(pwd)/backups:/backup \
  alpine tar czf /backup/postgres_data_$(date +%Y%m%d).tar.gz -C /data .

# Очистить неиспользуемые volumes
docker volume prune
```

### Сброс данных

```bash
# Остановить и удалить контейнеры
docker-compose down

# Удалить контейнеры и volumes (ВНИМАНИЕ: удалит все данные!)
docker-compose down -v

# Пересоздать с нуля
docker-compose down -v
docker-compose up -d
```

## Полезные алиасы

Добавьте в ~/.bashrc или ~/.zshrc:

```bash
# Docker Compose алиасы
alias dc='docker-compose'
alias dcup='docker-compose up -d'
alias dcdown='docker-compose down'
alias dclogs='docker-compose logs -f'
alias dcps='docker-compose ps'

# Быстрый доступ к БД
alias demo-psql='docker-compose exec postgres psql -U demo_user -d demo_db'
alias demo-redis='docker-compose exec redis redis-cli'

# Приложение
alias demo-status='go run main.go db status'
alias demo-serve='go run main.go serve --verbose'
```

## Troubleshooting

### Проблема: Порт уже занят

```bash
# Проверить, что занимает порт
lsof -i :5432
netstat -an | grep 5432

# Изменить порт в docker-compose.yml
ports:
  - "5433:5432"  # использовать 5433 на хосте
```

### Проблема: Контейнер не запускается

```bash
# Посмотреть логи
docker-compose logs postgres

# Проверить конфигурацию
docker-compose config

# Пересоздать контейнер
docker-compose up -d --force-recreate postgres
```

### Проблема: Нет доступа к БД

```bash
# Проверить сеть
docker network ls
docker network inspect cobra-viper-demo_demo-network

# Проверить переменные окружения
docker-compose exec postgres env | grep POSTGRES

# Тест подключения изнутри контейнера
docker-compose exec postgres pg_isready -U demo_user
```

### Проблема: Медленная работа на Mac/Windows

```bash
# Использовать delegated для volumes на Mac
volumes:
  - postgres_data:/var/lib/postgresql/data:delegated

# Или использовать named volumes вместо bind mounts
```

## Интеграция с CI/CD

### GitHub Actions пример

```yaml
name: Test with Docker

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Start services
        run: docker-compose up -d
      
      - name: Wait for PostgreSQL
        run: |
          until docker-compose exec -T postgres pg_isready; do
            sleep 1
          done
      
      - name: Run tests
        run: |
          docker-compose exec -T postgres psql -U demo_user -d demo_db -c "SELECT COUNT(*) FROM users;"
      
      - name: Cleanup
        run: docker-compose down -v
```

## Оптимизация

### Уменьшение размера образа

```dockerfile
# Использовать distroless образ
FROM gcr.io/distroless/static:nonroot
COPY --from=builder /build/app /app
ENTRYPOINT ["/app"]
```

### Кэширование слоёв

```dockerfile
# Сначала копировать go.mod/go.sum
COPY go.mod go.sum ./
RUN go mod download

# Потом исходники (для лучшего кэширования)
COPY . .
```

### Health checks

```yaml
healthcheck:
  test: ["CMD", "wget", "--spider", "http://localhost:8080/health"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 40s
```

## Безопасность

1. **Не коммитить .env файлы с паролями**
2. **Использовать secrets в production**
3. **Запускать от непривилегированного пользователя**
4. **Ограничивать ресурсы контейнеров**
5. **Использовать проверенные базовые образы**
6. **Регулярно обновлять образы**

```bash
# Сканирование уязвимостей
docker scan cobra-viper-demo:latest
```