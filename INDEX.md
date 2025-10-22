# Материалы вебинара: Docker Best Practices для Golang

## Навигация по материалам

### 📚 Документация

| Файл | Описание | Время чтения |
|------|----------|--------------|
| **[QUICKSTART.md](QUICKSTART.md)** | Быстрый старт (начните отсюда!) | 5 мин |
| **[README.md](README.md)** | Полная документация со всеми практиками | 30 мин |
| **[CHEATSHEET.md](CHEATSHEET.md)** | Шпаргалка по командам Docker | 10 мин |
| **[PRESENTATION.md](PRESENTATION.md)** | Презентация вебинара (слайды) | 20 мин |

### 🐳 Dockerfile примеры

| Файл | Описание | Размер | Сложность |
|------|----------|--------|-----------|
| **[Dockerfile.basic](dockerfiles/Dockerfile.basic)** | Базовый multi-stage с Alpine | ~25 MB | ⭐ Начальный |
| **[Dockerfile.optimized](dockerfiles/Dockerfile.optimized)** | Оптимизированный с кэшированием | ~20 MB | ⭐⭐ Средний |
| **[Dockerfile.distroless](dockerfiles/Dockerfile.distroless)** | Google Distroless для безопасности | ~15 MB | ⭐⭐⭐ Продвинутый |
| **[Dockerfile.scratch](dockerfiles/Dockerfile.scratch)** | Минимальный на базе scratch | ~8 MB | ⭐⭐⭐ Продвинутый |
| **[Dockerfile.multistage-advanced](dockerfiles/Dockerfile.multistage-advanced)** | С тестами и линтингом | ~20 MB | ⭐⭐⭐⭐ Эксперт |

### 🛠 Инструменты и автоматизация

| Файл | Описание | Использование |
|------|----------|---------------|
| **[Makefile](Makefile)** | Автоматизация всех команд | `make help` |
| **[scripts/build-all.sh](scripts/build-all.sh)** | Сборка всех вариантов | `./scripts/build-all.sh` |
| **[scripts/compare-sizes.sh](scripts/compare-sizes.sh)** | Сравнение размеров образов | `./scripts/compare-sizes.sh` |
| **[.dockerignore](.dockerignore)** | Пример .dockerignore файла | Скопируйте в проект |

### 🚀 Демо и примеры

| Файл | Описание | Команда запуска |
|------|----------|-----------------|
| **[docker-compose.demo.yml](docker-compose.demo.yml)** | Demo окружение со всеми вариантами | `make run-demo` |
| **[examples/github-actions.yml](examples/github-actions.yml)** | CI/CD пример для GitHub Actions | Адаптируйте под проект |

---

## Рекомендуемый порядок изучения

### Новичок в Docker

1. ⏱ **5 мин**: Прочитайте [QUICKSTART.md](QUICKSTART.md)
2. ⏱ **10 мин**: Соберите примеры с помощью `make build-all`
3. ⏱ **5 мин**: Посмотрите сравнение размеров `make compare`
4. ⏱ **30 мин**: Изучите [README.md](README.md) - секция "Лучшие практики"
5. ⏱ **15 мин**: Попробуйте адаптировать [Dockerfile.basic](dockerfiles/Dockerfile.basic) для своего проекта

**Итого**: ~1 час базового обучения

### Опытный разработчик

1. ⏱ **5 мин**: Быстрый просмотр [CHEATSHEET.md](CHEATSHEET.md)
2. ⏱ **10 мин**: Сравнение всех вариантов Dockerfile
3. ⏱ **20 мин**: Изучите [Dockerfile.distroless](dockerfiles/Dockerfile.distroless) и [Dockerfile.scratch](dockerfiles/Dockerfile.scratch)
4. ⏱ **15 мин**: Настройте CI/CD по примеру [examples/github-actions.yml](examples/github-actions.yml)
5. ⏱ **20 мин**: Проведите security audit с trivy

**Итого**: ~1 час продвинутого обучения

### Для вебинара / презентации

1. ⏱ **20 мин**: Изучите [PRESENTATION.md](PRESENTATION.md)
2. ⏱ **10 мин**: Подготовьте demo с `make run-demo`
3. ⏱ **30 мин**: Практикуйте live coding сборки образов
4. ⏱ **10 мин**: Подготовьте анализ с dive и trivy

**Итого**: ~1 час подготовки к вебинару

---

## Быстрые команды

### Сборка

```bash
cd webinar-materials

# Все варианты
make build-all

# Конкретный вариант
make build-basic
make build-optimized
make build-distroless
make build-scratch
make build-advanced
```

### Анализ

```bash
# Сравнение размеров
make compare

# Статистика
make stats

# Список образов
make list

# Детальный анализ (требует dive)
make dive-basic

# Security scan (требует trivy)
make security-scan
```

### Запуск

```bash
# Demo окружение
make run-demo

# Конкретный вариант
make run-basic
make run-optimized

# Просмотр логов
make logs-demo
```

### Очистка

```bash
# Удалить образы
make clean-images

# Полная очистка
make clean-all
```

---

## Ключевые концепции

### 1. Multi-Stage Build
**Файлы**: Все Dockerfile
**Результат**: Уменьшение размера с 800MB до 8-25MB

### 2. Кэширование слоев
**Файл**: [Dockerfile.optimized](dockerfiles/Dockerfile.optimized)
**Результат**: Ускорение сборки в 10-100 раз

### 3. Безопасность
**Файлы**: [Dockerfile.distroless](dockerfiles/Dockerfile.distroless), [Dockerfile.scratch](dockerfiles/Dockerfile.scratch)
**Результат**: Минимальная поверхность атаки

### 4. Оптимизация размера
**Все файлы** демонстрируют разные техники:
- Флаги компиляции `-ldflags="-s -w"`
- Выбор базового образа
- .dockerignore
- UPX сжатие (в optimized)

### 5. Production-ready
**Файл**: [Dockerfile.multistage-advanced](dockerfiles/Dockerfile.multistage-advanced)
**Включает**: тесты, линтинг, версионирование

---

## Сравнительная таблица

| Критерий | Basic | Optimized | Distroless | Scratch | Advanced |
|----------|-------|-----------|------------|---------|----------|
| **Размер** | 25 MB | 20 MB | 15 MB | 8 MB | 20 MB |
| **Безопасность** | Средняя | Средняя | Высокая | Максимальная | Высокая |
| **Отладка** | Легко | Легко | Сложно | Очень сложно | Легко |
| **Скорость сборки** | Быстро | Средне | Быстро | Быстро | Медленно |
| **Production** | ✅ | ✅✅ | ✅✅ | ✅✅ | ✅✅✅ |
| **Сложность** | ⭐ | ⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ |

**Рекомендации**:
- 🎓 **Обучение/Development**: Basic
- 🚀 **Production (стандарт)**: Optimized или Distroless
- 🔒 **Production (максимальная безопасность)**: Distroless или Scratch
- 🏭 **Enterprise/CI/CD**: Advanced

---

## Полезные ссылки

### Документация
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [Go Official Images](https://hub.docker.com/_/golang)
- [Distroless Images](https://github.com/GoogleContainerTools/distroless)
- [Multi-Stage Builds](https://docs.docker.com/build/building/multi-stage/)

### Инструменты
- [dive](https://github.com/wagoodman/dive) - Анализ слоев Docker образов
- [trivy](https://github.com/aquasecurity/trivy) - Security сканер
- [hadolint](https://github.com/hadolint/hadolint) - Dockerfile линтер

### Обучающие материалы
- [Docker для Go разработчиков](https://blog.golang.org/docker)
- [Dockerfile Security Best Practices](https://docs.docker.com/develop/security-best-practices/)

---

## Чек-лист для production

Используйте этот чек-лист при подготовке Dockerfile для production:

### Размер и производительность
- [ ] Multi-stage build используется
- [ ] Базовый образ минимален (alpine/distroless/scratch)
- [ ] Флаги компиляции оптимизированы (`-ldflags="-s -w"`)
- [ ] .dockerignore создан и настроен
- [ ] Порядок COPY команд оптимизирован для кэширования
- [ ] go.mod/go.sum копируются отдельно

### Безопасность
- [ ] Приложение запускается от непривилегированного пользователя
- [ ] Нет секретов в образе
- [ ] Использованы конкретные версии базовых образов (не :latest)
- [ ] Проведено сканирование на уязвимости (trivy)
- [ ] Минимизирована поверхность атаки

### Production-ready
- [ ] HEALTHCHECK добавлен (если применимо)
- [ ] LABEL метаданные добавлены
- [ ] Версионирование настроено
- [ ] EXPOSE порты документированы
- [ ] ENTRYPOINT/CMD правильно настроены
- [ ] Тесты проходят в образе
- [ ] Логирование настроено (stdout/stderr)

---

## FAQ

**Q: С чего начать?**
A: Откройте [QUICKSTART.md](QUICKSTART.md) и следуйте инструкциям.

**Q: Какой Dockerfile использовать для production?**
A: Для большинства случаев - [Dockerfile.optimized](dockerfiles/Dockerfile.optimized) или [Dockerfile.distroless](dockerfiles/Dockerfile.distroless).

**Q: Как уменьшить размер образа?**
A: См. секцию "Лучшие практики" в [README.md](README.md).

**Q: Как отлаживать distroless/scratch контейнеры?**
A: Создайте отдельный debug Dockerfile на базе Alpine или используйте distroless:debug тег.

**Q: Нужно ли использовать все эти практики сразу?**
A: Нет, начните с basic и постепенно добавляйте оптимизации.

---

## Обратная связь

Нашли ошибку или есть предложения?
- Создайте issue в репозитории
- Или отправьте pull request с улучшениями

---

## Лицензия

Материалы вебинара доступны для свободного использования в образовательных целях.

---

**Версия**: 1.0.0
**Дата**: 2025
**Автор**: Материалы для вебинара по Docker Best Practices для Golang

---

## Следующие шаги

После изучения материалов:

1. ✅ Примените практики к вашему проекту
2. ✅ Настройте CI/CD pipeline
3. ✅ Проведите security audit
4. ✅ Поделитесь знаниями с командой
5. ✅ Внесите вклад в развитие материалов

**Удачи! 🚀**

---

## Дополнительные материалы

### Kaniko и альтернативы Docker

**[KANIKO.md](KANIKO.md)** - Подробное руководство по инструментам для сборки образов без Docker daemon:
- Kaniko (Google) - для Kubernetes CI/CD
- BuildKit - следующее поколение Docker builder
- Podman - daemonless альтернатива Docker
- Buildah - scripting образов
- Сравнение и рекомендации

### CI/CD примеры

| Файл | Описание |
|------|----------|
| **[kaniko-gitlab-ci.yml](examples/kaniko-gitlab-ci.yml)** | GitLab CI/CD с Kaniko |
| **[buildkit-github-actions.yml](examples/buildkit-github-actions.yml)** | GitHub Actions с BuildKit |
| **[tekton-kaniko-pipeline.yml](examples/tekton-kaniko-pipeline.yml)** | Tekton Pipeline для Kubernetes |

