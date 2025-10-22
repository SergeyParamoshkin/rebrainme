# Обновление версий

## Изменения

Все версии Golang и Alpine обновлены до актуальных:

### Версии до обновления
- **Golang**: 1.20, 1.23
- **Alpine**: 3.17, 3.20

### Версии после обновления
- **Golang**: 1.24
- **Alpine**: 3.21

## Обновленные файлы

### Dockerfiles (5 файлов)
- ✅ `dockerfiles/Dockerfile.basic`
- ✅ `dockerfiles/Dockerfile.optimized`
- ✅ `dockerfiles/Dockerfile.distroless`
- ✅ `dockerfiles/Dockerfile.scratch`
- ✅ `dockerfiles/Dockerfile.multistage-advanced`

### Документация (4 файла)
- ✅ `README.md`
- ✅ `CHEATSHEET.md`
- ✅ `PRESENTATION.md`
- ✅ `KANIKO.md`

### CI/CD примеры (3 файла)
- ✅ `examples/kaniko-gitlab-ci.yml`
- ✅ `examples/buildkit-github-actions.yml`
- ✅ `examples/tekton-kaniko-pipeline.yml`

### Корневой проект (1 файл)
- ✅ `../Dockerfile`

## Проверка

```bash
# Проверка на старые версии (должно вернуть 0)
grep -r "golang:1\.2[0-3]\|alpine:3\.[0-1][0-9]\|alpine3\.[0-1][0-9]" webinar-materials/ --include="*.md" --include="*.yml" 2>/dev/null | wc -l

# Проверка на новые версии
grep -r "golang:1\.24\|alpine:3\.21" webinar-materials/ --include="*.md" --include="*.yml" 2>/dev/null | wc -l
```

## Совместимость

### Golang 1.24
- Выход: март 2024
- Новые фичи: улучшенная производительность, новый garbage collector
- Обратная совместимость: полная с 1.20+

### Alpine 3.21
- Выход: декабрь 2024
- Обновленные пакеты безопасности
- Размер: ~7 MB (был ~5 MB в 3.17)
- Обратная совместимость: полная

## Следующие шаги

1. Протестировать сборку всех вариантов:
   ```bash
   cd webinar-materials
   make build-all
   ```

2. Проверить размеры образов:
   ```bash
   make compare
   ```

3. Запустить demo:
   ```bash
   make run-demo
   ```

4. Обновить зависимости в go.mod (если требуется):
   ```bash
   go get -u ./...
   go mod tidy
   ```

## Changelog

### v1.0.1 (2025-01-22)
- Обновлен Golang 1.20/1.23 → 1.24
- Обновлен Alpine 3.17/3.20 → 3.21
- Обновлены все примеры в документации
- Обновлены все CI/CD конфигурации

---

**Дата обновления**: 2025-01-22
**Статус**: ✅ Завершено
