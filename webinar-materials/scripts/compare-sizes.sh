#!/usr/bin/env bash

# ============================================
# Скрипт для детального сравнения размеров образов
# Показывает размер каждого слоя и общую статистику
# ============================================

set -e

# Цвета
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

IMAGE_PREFIX="${1:-golang-demo}"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Docker Images Size Comparison${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Проверка наличия образов
if ! docker images "${IMAGE_PREFIX}" --format "{{.Repository}}" | grep -q "${IMAGE_PREFIX}"; then
    echo -e "${RED}No images found with prefix: ${IMAGE_PREFIX}${NC}"
    echo -e "${YELLOW}Run ./build-all.sh first${NC}"
    exit 1
fi

# Функция для получения размера в байтах
get_size_bytes() {
    local image=$1
    docker inspect "$image" --format='{{.Size}}' 2>/dev/null || echo "0"
}

# Функция для форматирования размера
format_size() {
    local bytes=$1
    if [ "$bytes" -lt 1024 ]; then
        echo "${bytes}B"
    elif [ "$bytes" -lt 1048576 ]; then
        echo "$(($bytes / 1024))KB"
    else
        echo "$(($bytes / 1048576))MB"
    fi
}

# Получение списка образов
echo -e "${YELLOW}Summary of all images:${NC}"
echo ""
docker images "${IMAGE_PREFIX}" --format "table {{.Repository}}:{{.Tag}}\t{{.Size}}\t{{.CreatedSince}}" | \
    head -20

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Detailed Size Analysis${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Массив для хранения размеров
declare -A sizes
declare -A layers_count

# Анализ каждого образа
for tag in basic optimized distroless scratch advanced; do
    image="${IMAGE_PREFIX}:${tag}"

    if docker images "$image" --format "{{.Repository}}" | grep -q "${IMAGE_PREFIX}"; then
        echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo -e "${YELLOW}Image: ${image}${NC}"
        echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

        # Размер образа
        size_bytes=$(get_size_bytes "$image")
        sizes[$tag]=$size_bytes
        size_human=$(format_size "$size_bytes")

        echo -e "Total size: ${GREEN}${size_human}${NC} ($(numfmt --to=iec-i --suffix=B "$size_bytes" 2>/dev/null || echo "$size_bytes bytes"))"

        # Количество слоев
        layers=$(docker history "$image" --format "{{.Size}}" --no-trunc | grep -v "0B" | wc -l | tr -d ' ')
        layers_count[$tag]=$layers
        echo -e "Layers count: ${GREEN}${layers}${NC}"

        # История слоев
        echo ""
        echo -e "${YELLOW}Layer breakdown:${NC}"
        docker history "$image" --format "table {{.Size}}\t{{.CreatedBy}}" --no-trunc | \
            head -20 | \
            awk '{if(NR==1) print $0; else if($1 != "0B") print $0}'

        echo ""
    fi
done

# Сравнительная таблица
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Size Comparison Table${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

printf "%-15s | %-12s | %-8s | %-15s\n" "Image Tag" "Size" "Layers" "Reduction"
echo "----------------+-------------+----------+----------------"

# Базовый размер для сравнения (basic)
base_size=${sizes[basic]:-0}

for tag in basic optimized distroless scratch advanced; do
    if [ -n "${sizes[$tag]}" ]; then
        size=${sizes[$tag]}
        layers=${layers_count[$tag]}
        size_human=$(format_size "$size")

        if [ "$size" -eq "$base_size" ]; then
            reduction="baseline"
        else
            percent=$((100 - (size * 100 / base_size)))
            if [ "$percent" -gt 0 ]; then
                reduction="-${percent}%"
            else
                reduction="+${percent#-}%"
            fi
        fi

        printf "%-15s | %12s | %8s | %15s\n" "$tag" "$size_human" "$layers" "$reduction"
    fi
done

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Recommendations${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Находим самый маленький образ
smallest_tag=""
smallest_size=999999999999

for tag in "${!sizes[@]}"; do
    if [ "${sizes[$tag]}" -lt "$smallest_size" ]; then
        smallest_size=${sizes[$tag]}
        smallest_tag=$tag
    fi
done

echo -e "${GREEN}Smallest image: ${smallest_tag} ($(format_size "$smallest_size"))${NC}"
echo ""
echo -e "${YELLOW}Guidelines:${NC}"
echo -e "  • ${CYAN}basic${NC} - Good for development, easy debugging"
echo -e "  • ${CYAN}optimized${NC} - Balanced approach with good caching"
echo -e "  • ${CYAN}distroless${NC} - Enhanced security, smaller attack surface"
echo -e "  • ${CYAN}scratch${NC} - Minimal size, maximum security (if no CGO)"
echo -e "  • ${CYAN}advanced${NC} - Production-ready with tests and versioning"
echo ""
