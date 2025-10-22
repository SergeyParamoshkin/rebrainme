# Быстрый старт

Краткое руководство для начала работы с материалами вебинара.

## Предварительные требования

```bash
# Проверка Docker
docker --version  # Должно быть >= 20.10

# Проверка Docker Compose
docker-compose --version  # Должно быть >= 2.0

# Опционально: полезные инструменты
make --version
```

## Шаг 1: Сборка примеров (5 минут)

### Вариант A: Используя Makefile (рекомендуется)

```bash
cd webinar-materials

# Посмотреть все доступные команды
make help

# Собрать все варианты Dockerfile
make build-all

# Сравнить размеры
make compare
```

### Вариант B: Используя скрипты

```bash
cd webinar-materials

# Собрать все образы
./scripts/build-all.sh

# Сравнить размеры
./scripts/compare-sizes.sh
```

### Вариант C: Вручную

```bash
cd webinar-materials

# Собрать базовый вариант
docker build -t golang-demo:basic -f dockerfiles/Dockerfile.basic ..

# Собрать оптимизированный
docker build -t golang-demo:optimized -f dockerfiles/Dockerfile.optimized ..

# И так далее...
```

## Шаг 2: Посмотреть результаты

```bash
# Список образов
docker images golang-demo

# Ожидаемый вывод:
# golang-demo  basic        ... 25MB
# golang-demo  optimized    ... 20MB
# golang-demo  distroless   ... 15MB
# golang-demo  scratch      ... 8MB
# golang-demo  advanced     ... 20MB
```

## Шаг 3: Запустить demo (опционально)

```bash
# Запустить все варианты через docker-compose
make run-demo

# Или напрямую
docker-compose -f docker-compose.demo.yml up -d

# Проверить работу
curl http://localhost:8081/health  # basic
curl http://localhost:8082/health  # optimized
# и т.д.

# Посмотреть логи
docker-compose -f docker-compose.demo.yml logs -f

# Остановить
make stop-demo
# или
docker-compose -f docker-compose.demo.yml down
```

## Шаг 4: Анализ образов (опционально)

### Использование dive

```bash
# Установка (MacOS)
brew install dive

# Анализ образа
dive golang-demo:basic

# Сравнение с оптимизированным
dive golang-demo:scratch
```

### Использование trivy

```bash
# Установка (MacOS)
brew install trivy

# Сканирование безопасности
trivy image golang-demo:basic

# Только критичные уязвимости
trivy image --severity CRITICAL,HIGH golang-demo:basic
```

## Шаг 5: Применить к своему проекту

### 1. Скопируйте нужный Dockerfile

```bash
# Для большинства случаев рекомендуется optimized
cp webinar-materials/dockerfiles/Dockerfile.optimized ./Dockerfile

# Или для максимальной безопасности
cp webinar-materials/dockerfiles/Dockerfile.distroless ./Dockerfile
```

### 2. Скопируйте .dockerignore

```bash
cp webinar-materials/.dockerignore ./.dockerignore
```

### 3. Адаптируйте под свой проект

Измените в Dockerfile:
- Пути к вашему main файлу (`./cmd` → ваш путь)
- Порты (если нужно)
- Переменные окружения

### 4. Соберите и проверьте

```bash
# Сборка
docker build -t myapp:1.0.0 .

# Проверка размера
docker images myapp

# Запуск
docker run -p 8080:8080 myapp:1.0.0

# Проверка работы
curl http://localhost:8080
```

## Частые проблемы

### Проблема 1: "no such file or directory" при запуске

**Причина**: Бинарник динамически слинкован, но библиотек нет в scratch/distroless

**Решение**: Убедитесь, что используете `CGO_ENABLED=0`:
```dockerfile
RUN CGO_ENABLED=0 go build ...
```

### Проблема 2: "x509: certificate signed by unknown authority"

**Причина**: Отсутствуют CA сертификаты

**Решение**: Добавьте в Dockerfile:
```dockerfile
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
```

### Проблема 3: Медленная сборка

**Причина**: Неправильный порядок COPY команд

**Решение**: Копируйте go.mod/go.sum перед исходным кодом:
```dockerfile
COPY go.mod go.sum ./
RUN go mod download
COPY . .
```

### Проблема 4: Большой размер образа

**Решение**: Проверьте:
1. Используете ли multi-stage build?
2. Добавлен ли .dockerignore?
3. Используете ли флаги `-ldflags="-s -w"`?
4. Выбран ли правильный базовый образ?

## Следующие шаги

1. **Изучите документацию**
   - Откройте [README.md](README.md) для детального руководства
   - Посмотрите [CHEATSHEET.md](CHEATSHEET.md) для быстрой справки

2. **Настройте CI/CD**
   - Используйте примеры из [examples/github-actions.yml](examples/github-actions.yml)
   - Адаптируйте под вашу CI/CD систему

3. **Проведите security audit**
   - Запустите trivy на ваших образах
   - Исправьте найденные проблемы

4. **Оптимизируйте production сборку**
   - Добавьте версионирование
   - Настройте health checks
   - Добавьте метрики и логирование

## Полезные команды

```bash
# Быстрая справка
make help

# Сборка всех вариантов
make build-all

# Сравнение размеров
make compare

# Запуск demo окружения
make run-demo

# Остановка demo
make stop-demo

# Очистка образов
make clean-images

# Полная очистка
make clean-all

# Security scan (требует trivy)
make security-scan

# Статистика образов
make stats
```

## Структура материалов

```
webinar-materials/
├── README.md                    # Подробная документация
├── QUICKSTART.md               # Этот файл
├── CHEATSHEET.md               # Шпаргалка по Docker
├── PRESENTATION.md             # Презентация вебинара
├── Makefile                    # Автоматизация команд
├── .dockerignore               # Пример .dockerignore
├── docker-compose.demo.yml     # Demo окружение
├── dockerfiles/                # Варианты Dockerfile
│   ├── Dockerfile.basic
│   ├── Dockerfile.optimized
│   ├── Dockerfile.distroless
│   ├── Dockerfile.scratch
│   └── Dockerfile.multistage-advanced
├── scripts/                    # Утилиты
│   ├── build-all.sh
│   └── compare-sizes.sh
└── examples/                   # Примеры
    └── github-actions.yml
```

## Получение помощи

Если что-то не работает:

1. Проверьте версии Docker и Docker Compose
2. Посмотрите логи: `docker logs <container-name>`
3. Проверьте [README.md](README.md) для детальной информации
4. Откройте [CHEATSHEET.md](CHEATSHEET.md) для справки по командам

## Обратная связь

Если вы нашли ошибку или у вас есть предложения по улучшению материалов, создайте issue или pull request в репозитории.

---

Удачи с вашим вебинаром! 🚀
