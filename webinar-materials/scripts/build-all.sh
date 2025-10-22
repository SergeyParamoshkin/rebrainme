#!/usr/bin/env bash

# ============================================
# Скрипт для сборки всех вариантов Dockerfile
# и сравнения их размеров
# ============================================

set -e

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Переменные
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"
DOCKERFILES_DIR="$SCRIPT_DIR/../dockerfiles"
IMAGE_PREFIX="golang-demo"
VERSION="${VERSION:-1.0.0}"
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Building all Docker images${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "Project root: ${GREEN}$PROJECT_ROOT${NC}"
echo -e "Version: ${GREEN}$VERSION${NC}"
echo -e "Build time: ${GREEN}$BUILD_TIME${NC}"
echo -e "Git commit: ${GREEN}$GIT_COMMIT${NC}"
echo ""

# Функция для сборки образа
build_image() {
    local dockerfile=$1
    local image_name=$2
    local build_args=$3

    echo -e "${YELLOW}Building: ${image_name}${NC}"
    echo -e "Dockerfile: ${dockerfile}"

    if [ -n "$build_args" ]; then
        docker build \
            $build_args \
            -t "${IMAGE_PREFIX}:${image_name}" \
            -f "$dockerfile" \
            "$PROJECT_ROOT" 2>&1 | tail -n 5
    else
        docker build \
            -t "${IMAGE_PREFIX}:${image_name}" \
            -f "$dockerfile" \
            "$PROJECT_ROOT" 2>&1 | tail -n 5
    fi

    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Successfully built: ${image_name}${NC}"
    else
        echo -e "${RED}✗ Failed to build: ${image_name}${NC}"
        return 1
    fi
    echo ""
}

# Массив с информацией о сборках
declare -a BUILDS=(
    "Dockerfile.basic:basic:"
    "Dockerfile.optimized:optimized:--build-arg VERSION=$VERSION --build-arg BUILD_TIME=$BUILD_TIME"
    "Dockerfile.distroless:distroless:"
    "Dockerfile.scratch:scratch:"
    "Dockerfile.multistage-advanced:advanced:--build-arg VERSION=$VERSION --build-arg BUILD_TIME=$BUILD_TIME --build-arg GIT_COMMIT=$GIT_COMMIT --build-arg SKIP_TESTS=true --build-arg SKIP_LINT=true"
)

# Сборка всех образов
for build in "${BUILDS[@]}"; do
    IFS=':' read -r dockerfile tag args <<< "$build"
    build_image "$DOCKERFILES_DIR/$dockerfile" "$tag" "$args"
done

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Build Summary${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Вывод информации о размерах
echo -e "${YELLOW}Image sizes:${NC}"
docker images "${IMAGE_PREFIX}" --format "table {{.Repository}}:{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}" | \
    grep -E "(REPOSITORY|${IMAGE_PREFIX})" | \
    column -t

echo ""
echo -e "${GREEN}All images built successfully!${NC}"
echo ""
echo -e "${YELLOW}Useful commands:${NC}"
echo -e "  View detailed info: ${BLUE}docker images ${IMAGE_PREFIX}${NC}"
echo -e "  Inspect image: ${BLUE}docker inspect ${IMAGE_PREFIX}:basic${NC}"
echo -e "  Run container: ${BLUE}docker run -p 8080:8080 ${IMAGE_PREFIX}:basic${NC}"
echo -e "  Clean up: ${BLUE}docker rmi \$(docker images ${IMAGE_PREFIX} -q)${NC}"
echo ""
