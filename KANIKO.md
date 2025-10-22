# Kaniko и альтернативы Docker для сборки образов

## Содержание

1. [Введение](#введение)
2. [Kaniko](#kaniko)
3. [BuildKit](#buildkit)
4. [Podman](#podman)
5. [Buildah](#buildah)
6. [img](#img)
7. [Сравнение инструментов](#сравнение-инструментов)
8. [Практические примеры](#практические-примеры)

---

## Введение

### Проблема с Docker daemon

**Традиционный подход**:
```bash
docker build -t myapp:latest .
```

**Проблемы**:
- ❌ Требует Docker daemon (root privileges)
- ❌ Не работает в Kubernetes pods без DinD (Docker-in-Docker)
- ❌ Безопасность: доступ к Docker socket = root доступ к хосту
- ❌ DinD сложен в настройке и небезопасен

### Решение: Daemonless сборка

Инструменты для сборки Docker образов **без** Docker daemon:
- Kaniko (Google)
- BuildKit (Docker Inc)
- Buildah (Red Hat)
- Podman (Red Hat)
- img

---

## Kaniko

### Что такое Kaniko?

**Kaniko** - инструмент от Google для сборки Docker образов из Dockerfile внутри контейнера или Kubernetes cluster без привилегированного доступа.

### Ключевые особенности

✅ **Не требует Docker daemon**
✅ **Работает в Kubernetes/OpenShift**
✅ **Безопасность**: не нужен privileged mode
✅ **Кэширование**: поддержка remote cache (registry)
✅ **Multi-stage builds**: полная поддержка
✅ **100% совместим с Dockerfile**

### Как работает?

1. Парсит Dockerfile
2. Выполняет каждую команду в userspace
3. Создает snapshot файловой системы после каждого шага
4. Push образа в registry

### Установка и использование

#### Docker

```bash
docker run \
  -v $(pwd):/workspace \
  gcr.io/kaniko-project/executor:latest \
  --dockerfile=/workspace/Dockerfile \
  --context=/workspace \
  --destination=myregistry/myapp:latest
```

#### Kubernetes Job

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: kaniko-build
spec:
  template:
    spec:
      containers:
      - name: kaniko
        image: gcr.io/kaniko-project/executor:latest
        args:
        - "--dockerfile=/workspace/Dockerfile"
        - "--context=dir:///workspace"
        - "--destination=myregistry/myapp:1.0.0"
        - "--cache=true"
        - "--cache-repo=myregistry/cache"
        volumeMounts:
        - name: dockerfile
          mountPath: /workspace
        - name: docker-config
          mountPath: /kaniko/.docker/
      restartPolicy: Never
      volumes:
      - name: dockerfile
        configMap:
          name: dockerfile
      - name: docker-config
        secret:
          secretName: regcred
```

#### GitLab CI пример

```yaml
# .gitlab-ci.yml
build:
  stage: build
  image:
    name: gcr.io/kaniko-project/executor:debug
    entrypoint: [""]
  script:
    - mkdir -p /kaniko/.docker
    - echo "{\"auths\":{\"$CI_REGISTRY\":{\"auth\":\"$(echo -n $CI_REGISTRY_USER:$CI_REGISTRY_PASSWORD | base64)\"}}}" > /kaniko/.docker/config.json
    - /kaniko/executor
      --context $CI_PROJECT_DIR
      --dockerfile $CI_PROJECT_DIR/Dockerfile
      --destination $CI_REGISTRY_IMAGE:$CI_COMMIT_TAG
      --cache=true
      --cache-repo=$CI_REGISTRY_IMAGE/cache
```

### Кэширование в Kaniko

```bash
# С кэшированием в registry
docker run \
  -v $(pwd):/workspace \
  gcr.io/kaniko-project/executor:latest \
  --dockerfile=/workspace/Dockerfile \
  --context=/workspace \
  --destination=myregistry/myapp:latest \
  --cache=true \
  --cache-repo=myregistry/cache \
  --cache-ttl=24h
```

### Преимущества

✅ **Безопасность**: не нужен root или Docker socket
✅ **Kubernetes native**: идеально для CI/CD в K8s
✅ **Кэширование**: поддержка remote cache layers
✅ **Reproducible builds**: одинаковые образы везде

### Недостатки

❌ Медленнее чем BuildKit в некоторых случаях
❌ Некоторые edge cases с Dockerfile могут не работать
❌ Требует registry для кэша (нет local cache)

---

## BuildKit

### Что такое BuildKit?

**BuildKit** - следующее поколение Docker builder от Docker Inc. Встроен в Docker 18.09+.

### Ключевые особенности

✅ **Параллельная сборка**: независимые stages собираются одновременно
✅ **Улучшенное кэширование**: inline cache, remote cache
✅ **Новый синтаксис**: BuildKit-specific features
✅ **Rootless mode**: сборка без root
✅ **Экспорт в разные форматы**: OCI, tar, registry
✅ **Secrets management**: безопасная передача секретов

### Использование

#### Включение BuildKit в Docker

```bash
# Временно
DOCKER_BUILDKIT=1 docker build -t myapp .

# Постоянно (в ~/.docker/config.json)
{
  "features": {
    "buildkit": true
  }
}
```

#### BuildKit CLI (buildctl)

```bash
# Standalone BuildKit
buildctl build \
  --frontend dockerfile.v0 \
  --local context=. \
  --local dockerfile=. \
  --output type=image,name=myapp:latest,push=true
```

#### Dockerfile с BuildKit features

```dockerfile
# syntax=docker/dockerfile:1.4

FROM golang:1.24-alpine AS builder

# BuildKit cache mount - кэширование go mod cache
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# BuildKit secret - безопасная передача токенов
RUN --mount=type=secret,id=github_token \
    GITHUB_TOKEN=$(cat /run/secrets/github_token) \
    go get private-repo

# BuildKit SSH mount - для git clone
RUN --mount=type=ssh \
    git clone git@github.com:user/private-repo.git

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o app ./cmd
```

#### GitHub Actions с BuildKit

```yaml
name: BuildKit Build

on: [push]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: .
          push: true
          tags: myregistry/myapp:latest
          cache-from: type=registry,ref=myregistry/myapp:buildcache
          cache-to: type=registry,ref=myregistry/myapp:buildcache,mode=max
          secrets: |
            github_token=${{ secrets.GITHUB_TOKEN }}
```

### BuildKit advanced features

#### 1. Cache mounts

```dockerfile
# Кэширование go modules между сборками
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o app
```

#### 2. Inline cache

```bash
# Сборка с inline cache
docker buildx build \
  --cache-from=type=registry,ref=myapp:latest \
  --cache-to=type=inline \
  --push \
  -t myapp:latest .
```

#### 3. Multi-platform builds

```bash
docker buildx build \
  --platform linux/amd64,linux/arm64,linux/arm/v7 \
  -t myapp:latest \
  --push .
```

### Преимущества

✅ **Скорость**: параллельная сборка
✅ **Продвинутое кэширование**: cache mounts, remote cache
✅ **Современные фичи**: secrets, SSH, multi-platform
✅ **Совместимость**: работает с Docker

### Недостатки

❌ Требует Docker (в обычном режиме)
❌ Rootless mode все еще экспериментальный
❌ Сложнее настройка в CI/CD

---

## Podman

### Что такое Podman?

**Podman** - daemonless альтернатива Docker от Red Hat. "POD Manager".

### Ключевые особенности

✅ **Daemonless**: нет фонового процесса
✅ **Rootless**: работает без root по умолчанию
✅ **Docker-compatible**: `alias docker=podman`
✅ **Systemd integration**: управление как systemd services
✅ **Pods support**: Kubernetes-style pods
✅ **Безопасность**: user namespaces, SELinux

### Использование

```bash
# Установка (RHEL/Fedora)
sudo dnf install podman

# Установка (Ubuntu)
sudo apt install podman

# Использование (как Docker)
podman build -t myapp .
podman run -p 8080:8080 myapp
podman push myapp:latest

# Rootless
podman build --userns=keep-id -t myapp .
```

### Podman в CI/CD

```yaml
# GitLab CI
build:
  image: quay.io/podman/stable
  script:
    - podman login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY
    - podman build -t $CI_REGISTRY_IMAGE:$CI_COMMIT_TAG .
    - podman push $CI_REGISTRY_IMAGE:$CI_COMMIT_TAG
```

### Преимущества

✅ **Безопасность**: rootless по умолчанию
✅ **Простота**: нет daemon
✅ **Совместимость**: Docker CLI compatible
✅ **Pods**: поддержка Kubernetes pods

### Недостатки

❌ Меньше распространен чем Docker
❌ Некоторые фичи Docker могут отсутствовать
❌ Performance в некоторых случаях ниже

---

## Buildah

### Что такое Buildah?

**Buildah** - инструмент от Red Hat для создания OCI-совместимых образов. Фокус на scriptable builds.

### Ключевые особенности

✅ **Daemonless и rootless**
✅ **Scripting**: создание образов через Bash
✅ **Dockerfile support**: может использовать Dockerfile
✅ **Fine-grained control**: полный контроль над слоями
✅ **OCI compliant**: создает OCI образы

### Использование

#### С Dockerfile

```bash
# Установка
sudo dnf install buildah

# Сборка из Dockerfile
buildah bud -t myapp .

# Push
buildah push myapp docker://registry.example.com/myapp:latest
```

#### Scripting (без Dockerfile)

```bash
#!/bin/bash

# Создать новый контейнер
container=$(buildah from golang:1.24-alpine)

# Смонтировать файловую систему
buildah copy $container . /app

# Выполнить команды
buildah run $container go build -o /app/bin ./cmd

# Создать финальный образ
final=$(buildah from alpine:latest)
buildah copy $final --from=$container /app/bin /usr/local/bin/app

# Настроить
buildah config --entrypoint '["/usr/local/bin/app"]' $final
buildah config --author "Me" $final

# Commit
buildah commit $final myapp:latest

# Cleanup
buildah rm $container $final
```

### Buildah в CI/CD

```yaml
# GitHub Actions
name: Buildah Build

on: [push]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install Buildah
        run: |
          sudo apt-get update
          sudo apt-get install -y buildah

      - name: Build
        run: |
          buildah bud -t myapp .

      - name: Push
        run: |
          buildah login -u ${{ secrets.REGISTRY_USER }} -p ${{ secrets.REGISTRY_PASSWORD }} registry.example.com
          buildah push myapp docker://registry.example.com/myapp:latest
```

### Преимущества

✅ **Гибкость**: scriptable builds
✅ **Безопасность**: rootless
✅ **Точный контроль**: управление каждым слоем
✅ **OCI**: полная совместимость с OCI

### Недостатки

❌ Более сложен для простых случаев
❌ Требует больше знаний
❌ Меньше документации

---

## img

### Что такое img?

**img** - standalone daemon-less unprivileged Dockerfile builder от Jess Frazelle.

### Ключевые особенности

✅ **Полностью unprivileged**
✅ **Standalone**: один бинарник
✅ **Dockerfile compatible**
✅ **Быстрый**: использует runc

### Использование

```bash
# Установка
curl -fSL "https://github.com/genuinetools/img/releases/download/v0.5.11/img-linux-amd64" -o /usr/local/bin/img
chmod +x /usr/local/bin/img

# Сборка
img build -t myapp .

# Push
img push myapp:latest
```

### Преимущества

✅ **Простота**: один бинарник
✅ **Безопасность**: unprivileged
✅ **Совместимость**: Dockerfile

### Недостатки

❌ Проект менее активно развивается
❌ Меньше функций чем у конкурентов
❌ Ограниченная поддержка

---

## Сравнение инструментов

### Таблица сравнения

| Критерий | Kaniko | BuildKit | Podman | Buildah | img |
|----------|--------|----------|--------|---------|-----|
| **Daemonless** | ✅ | ❌ (rootless mode) | ✅ | ✅ | ✅ |
| **Rootless** | ✅ | ⚠️ (experimental) | ✅ | ✅ | ✅ |
| **Скорость** | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ |
| **Kubernetes** | ✅✅ | ⚠️ | ✅ | ✅ | ✅ |
| **Кэширование** | Registry | Local/Remote | Local | Local | Local |
| **Dockerfile** | 100% | 100%+ | 100% | 100% | 100% |
| **Multi-platform** | ❌ | ✅ | ✅ | ✅ | ❌ |
| **Сложность** | Средняя | Средняя | Низкая | Высокая | Низкая |
| **Поддержка** | Google | Docker Inc | Red Hat | Red Hat | Community |
| **Use Case** | K8s CI/CD | Docker user | RHEL/Podman | Scripting | Simple |

### Рекомендации по выбору

#### Kaniko - выбирайте если:
- ✅ Сборка в Kubernetes/OpenShift
- ✅ Нужна полная безопасность (rootless)
- ✅ GitLab CI, Tekton, Argo Workflows
- ✅ Remote cache в registry

#### BuildKit - выбирайте если:
- ✅ Используете Docker
- ✅ Нужна максимальная скорость
- ✅ Multi-platform builds
- ✅ Продвинутые фичи (cache mounts, secrets)
- ✅ GitHub Actions, Jenkins

#### Podman - выбирайте если:
- ✅ Red Hat/Fedora/RHEL окружение
- ✅ Нужна полная замена Docker
- ✅ Rootless по умолчанию
- ✅ Systemd integration

#### Buildah - выбирайте если:
- ✅ Нужен scripting образов
- ✅ Точный контроль над слоями
- ✅ Автоматизация без Dockerfile
- ✅ Red Hat экосистема

#### img - выбирайте если:
- ✅ Нужна простота
- ✅ Standalone бинарник
- ✅ Базовые требования

---

## Практические примеры

### Пример 1: Kaniko в GitHub Actions

```yaml
name: Kaniko Build

on:
  push:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    container:
      image: gcr.io/kaniko-project/executor:debug
      options: --user root
    steps:
      - uses: actions/checkout@v4

      - name: Build and push
        env:
          REGISTRY_USER: ${{ secrets.REGISTRY_USER }}
          REGISTRY_PASSWORD: ${{ secrets.REGISTRY_PASSWORD }}
        run: |
          echo "{\"auths\":{\"registry.example.com\":{\"auth\":\"$(echo -n $REGISTRY_USER:$REGISTRY_PASSWORD | base64)\"}}}" > /kaniko/.docker/config.json

          /kaniko/executor \
            --context $GITHUB_WORKSPACE \
            --dockerfile $GITHUB_WORKSPACE/Dockerfile \
            --destination registry.example.com/myapp:$GITHUB_SHA \
            --destination registry.example.com/myapp:latest \
            --cache=true \
            --cache-repo=registry.example.com/myapp/cache
```

### Пример 2: BuildKit с cache mounts

```dockerfile
# syntax=docker/dockerfile:1.4

FROM golang:1.24-alpine AS builder

WORKDIR /build

# Cache go modules
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download

COPY . .

# Cache go build
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o app ./cmd

FROM scratch
COPY --from=builder /build/app /app
ENTRYPOINT ["/app"]
```

```bash
# Сборка с BuildKit
DOCKER_BUILDKIT=1 docker build -t myapp .
```

### Пример 3: Tekton Pipeline с Kaniko

```yaml
apiVersion: tekton.dev/v1beta1
kind: Pipeline
metadata:
  name: build-pipeline
spec:
  params:
    - name: IMAGE
      type: string
      default: registry.example.com/myapp
  workspaces:
    - name: shared-workspace
  tasks:
    - name: fetch-repository
      taskRef:
        name: git-clone
      workspaces:
        - name: output
          workspace: shared-workspace
      params:
        - name: url
          value: https://github.com/user/repo.git

    - name: build-push
      runAfter: [fetch-repository]
      taskRef:
        name: kaniko
      workspaces:
        - name: source
          workspace: shared-workspace
      params:
        - name: IMAGE
          value: $(params.IMAGE):$(tasks.fetch-repository.results.commit)
        - name: EXTRA_ARGS
          value:
            - --cache=true
            - --cache-repo=$(params.IMAGE)/cache
```

### Пример 4: Podman в systemd

```bash
# Создать Containerfile (Dockerfile)
cat > Containerfile <<EOF
FROM golang:1.24-alpine
WORKDIR /app
COPY . .
RUN go build -o myapp
CMD ["./myapp"]
EOF

# Собрать
podman build -t myapp .

# Создать systemd service
podman generate systemd --new --name myapp > ~/.config/systemd/user/myapp.service

# Запустить как service
systemctl --user enable --now myapp.service
```

### Пример 5: Buildah script для multi-stage

```bash
#!/bin/bash
set -e

# Stage 1: Build
builder=$(buildah from golang:1.24-alpine)
buildah copy $builder . /src
buildah config --workingdir /src $builder
buildah run $builder go build -o app ./cmd

# Stage 2: Runtime
runtime=$(buildah from alpine:3.21)
buildah copy $runtime --from=$builder /src/app /usr/local/bin/app
buildah config --entrypoint '["/usr/local/bin/app"]' $runtime
buildah config --port 8080 $runtime

# Commit and push
buildah commit $runtime myapp:latest
buildah push myapp:latest docker://registry.example.com/myapp:latest

# Cleanup
buildah rm $builder $runtime
```

---

## Выводы

### Когда использовать что?

| Сценарий | Инструмент | Причина |
|----------|-----------|---------|
| **Kubernetes CI/CD** | Kaniko | Rootless, K8s native |
| **GitHub Actions** | BuildKit | Скорость, cache, официальная поддержка |
| **GitLab CI в K8s** | Kaniko | Интеграция, безопасность |
| **Jenkins** | BuildKit или Podman | Гибкость |
| **RHEL/Fedora** | Podman/Buildah | Нативная поддержка |
| **Scripted builds** | Buildah | Полный контроль |
| **Простые случаи** | Podman или img | Простота |
| **Multi-platform** | BuildKit | Лучшая поддержка |
| **Maximum security** | Kaniko или Buildah | Rootless, unprivileged |

### Тренды

1. **Daemonless становится стандартом** - уход от Docker daemon
2. **Rootless по умолчанию** - безопасность на первом месте
3. **Kubernetes-native** - инструменты для K8s
4. **BuildKit features** - становятся индустриальным стандартом
5. **OCI compliance** - все движутся к OCI спецификации

---

## Дополнительные ресурсы

### Документация
- [Kaniko](https://github.com/GoogleContainerTools/kaniko)
- [BuildKit](https://github.com/moby/buildkit)
- [Podman](https://podman.io/)
- [Buildah](https://buildah.io/)
- [img](https://github.com/genuinetools/img)

### Сравнительные статьи
- [Kaniko vs BuildKit](https://blog.alexellis.io/building-containers-without-docker/)
- [Podman vs Docker](https://developers.redhat.com/blog/2020/11/19/transitioning-from-docker-to-podman)

---

**Заключение**: Выбор инструмента зависит от вашей инфраструктуры, требований безопасности и use case. Для большинства CI/CD в Kubernetes - Kaniko, для локальной разработки с Docker - BuildKit, для Red Hat экосистемы - Podman/Buildah.
