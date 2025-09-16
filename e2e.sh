#!/bin/bash

# Цвета для красивого вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Функция для красивого заголовка
print_header() {
    echo
    echo -e "${BOLD}${CYAN}════════════════════════════════════════════════════════════════${NC}"
    echo -e "${BOLD}${CYAN}  $1${NC}"
    echo -e "${BOLD}${CYAN}════════════════════════════════════════════════════════════════${NC}"
    echo
}

# Функция для подзаголовка
print_subheader() {
    echo -e "${BOLD}${YELLOW}▶ $1${NC}"
    echo
}

# Главное меню
show_menu() {
    print_header "🚀 Rate Limiting & Retry Demo"
    echo -e "${BOLD}Выберите демонстрацию:${NC}"
    echo
    echo -e "  ${GREEN}1)${NC} 🚦 Rate Limiting - Sequential (демонстрация 429 ошибок)"
    echo -e "  ${GREEN}2)${NC} 🔥 Rate Limiting - Parallel Burst (параллельные запросы)"
    echo -e "  ${GREEN}3)${NC} ⏱  Sliding Window (временное окно)"
    echo -e "  ${GREEN}4)${NC} 💥 Flaky Service (случайные 503 ошибки)"
    echo -e "  ${GREEN}5)${NC} 🔄 Circuit Breaker (защита от каскадных сбоев)"
    echo -e "  ${GREEN}6)${NC} 🧪 Run All Tests (запустить все тесты)"
    echo -e "  ${GREEN}7)${NC} 📊 Load Test with Monitoring (нагрузочный тест)"
    echo -e "  ${GREEN}8)${NC} 🎯 Custom Test (настраиваемый тест)"
    echo -e "  ${GREEN}0)${NC} ❌ Exit"
    echo
    echo -n "Ваш выбор: "
}

# Проверка, запущен ли сервер
check_server() {
    if ! curl -s http://localhost:8080/health > /dev/null 2>&1; then
        echo -e "${YELLOW}⚠️  Сервер не запущен. Запускаю...${NC}"
        go run cmd/main.go &
        SERVER_PID=$!
        sleep 2
        if curl -s http://localhost:8080/health > /dev/null 2>&1; then
            echo -e "${GREEN}✓ Сервер запущен (PID: $SERVER_PID)${NC}"
        else
            echo -e "${RED}✗ Не удалось запустить сервер${NC}"
            exit 1
        fi
    else
        echo -e "${GREEN}✓ Сервер уже запущен${NC}"
    fi
}

# Rate limiting sequential test
test_rate_limiting_sequential() {
    print_subheader "Отправляю 15 последовательных запросов (лимит: 10/сек)"

    for i in {1..100}; do
        response=$(curl -s -w "\n%{http_code}" http://localhost:8080/api/simple 2>/dev/null)
        http_code=$(echo "$response" | tail -n1)

        if [ "$http_code" = "200" ]; then
            echo -e "  ${GREEN}✓${NC} Request #$i: ${GREEN}200 OK${NC}"
        elif [ "$http_code" = "429" ]; then
            echo -e "  ${YELLOW}⚠${NC} Request #$i: ${YELLOW}429 Too Many Requests${NC}"
        else
            echo -e "  ${RED}✗${NC} Request #$i: ${RED}$http_code${NC}"
        fi

        sleep 0.05
    done
}

# Rate limiting parallel test
test_rate_limiting_parallel() {
    print_subheader "Запускаю 20 параллельных запросов"

    echo -e "${CYAN}Запуск...${NC}"

    # Запускаем запросы в фоне
    for i in {1..100}; do
        {
            response=$(curl -s -w "\n%{http_code}" http://localhost:8080/api/simple 2>/dev/null)
            http_code=$(echo "$response" | tail -n1)

            if [ "$http_code" = "200" ]; then
                echo -e "  ${GREEN}✓${NC} Request #$i: ${GREEN}200 OK${NC}"
            elif [ "$http_code" = "429" ]; then
                echo -e "  ${YELLOW}⚠${NC} Request #$i: ${YELLOW}429 Rate Limited${NC}"
            else
                echo -e "  ${RED}✗${NC} Request #$i: ${RED}$http_code${NC}"
            fi
        } &
    done

    # Ждём завершения всех запросов
    wait

    echo -e "\n${GREEN}Все запросы завершены${NC}"
}

# Sliding window test
test_sliding_window() {
    print_subheader "Тест sliding window (лимит: 100 запросов/минута)"

    echo "Отправляю 110 запросов..."
    success=0
    limited=0

    for i in {1..110}; do
        response=$(curl -s -w "\n%{http_code}" http://localhost:8080/api/simple 2>/dev/null)
        http_code=$(echo "$response" | tail -n1)

        if [ "$http_code" = "200" ]; then
            ((success++))
            echo -ne "\r  Progress: $i/110 | ${GREEN}✓ Success: $success${NC} | ${YELLOW}⚠ Limited: $limited${NC}"
        elif [ "$http_code" = "429" ]; then
            ((limited++))
            echo -ne "\r  Progress: $i/110 | ${GREEN}✓ Success: $success${NC} | ${YELLOW}⚠ Limited: $limited${NC}"
        fi

        sleep 0.01
    done

    echo -e "\n\n${BOLD}Результаты:${NC}"
    echo -e "  ${GREEN}Успешных запросов: $success${NC}"
    echo -e "  ${YELLOW}Заблокировано: $limited${NC}"
}

# Flaky service test
test_flaky_service() {
    print_subheader "Тестирую нестабильный сервис (50% failure rate)"

    success=0
    failures=0

    for i in {1..20}; do
        response=$(curl -s -w "\n%{http_code}" http://localhost:8080/api/flaky 2>/dev/null)
        http_code=$(echo "$response" | tail -n1)

        if [ "$http_code" = "200" ]; then
            ((success++))
            echo -e "  ${GREEN}✓${NC} Request #$i: ${GREEN}Success!${NC}"
        else
            ((failures++))
            echo -e "  ${RED}✗${NC} Request #$i: ${RED}503 Service Unavailable${NC}"
        fi

        sleep 0.1
    done

    success_rate=$((success * 100 / 20))
    echo -e "\n${BOLD}Статистика:${NC}"
    echo -e "  Success rate: ${BOLD}$success_rate%${NC} (ожидаемо: ~50%)"
}

# Circuit breaker test
test_circuit_breaker() {
    print_subheader "Демонстрация Circuit Breaker"

    echo -e "${CYAN}Circuit breaker открывается после 5 ошибок${NC}\n"

    for i in {1..15}; do
        response=$(curl -s http://localhost:8080/api/circuit-breaker)

        # Попытка распарсить JSON
        if echo "$response" | grep -q "state"; then
            state=$(echo "$response" | grep -oP '"state"\s*:\s*"\K[^"]+' || echo "unknown")

            if echo "$response" | grep -q "error"; then
                echo -e "  ${RED}✗${NC} Request #$i: Circuit $state - Backend error"
            elif [ "$state" = "open" ]; then
                echo -e "  ${YELLOW}⚡${NC} Request #$i: Circuit ${YELLOW}OPEN${NC} - Fast fail!"
            else
                echo -e "  ${GREEN}✓${NC} Request #$i: Circuit ${GREEN}CLOSED${NC} - Success"
            fi
        else
            echo -e "  ${PURPLE}?${NC} Request #$i: Unknown response"
        fi

        sleep 0.5
    done
}

# Load test с мониторингом
load_test_with_monitoring() {
    print_subheader "Нагрузочный тест с мониторингом"

    echo -e "${CYAN}Открываю метрики в новом окне...${NC}"

    # Запускаем мониторинг метрик в новом терминале (если возможно)
    if command -v gnome-terminal &> /dev/null; then
        gnome-terminal -- bash -c "watch -n 1 'curl -s http://localhost:8080/metrics | grep -E \"(rate_limited|http_requests_total|retry)\"'" &
    elif command -v xterm &> /dev/null; then
        xterm -e "watch -n 1 'curl -s http://localhost:8080/metrics | grep -E \"(rate_limited|http_requests_total|retry)\"'" &
    else
        echo -e "${YELLOW}Метрики доступны по адресу: http://localhost:8080/metrics${NC}"
    fi

    echo -e "\n${BOLD}Запускаю нагрузку...${NC}\n"

    # Используем wrk если доступен
    if command -v wrk &> /dev/null; then
        echo "Использую wrk для нагрузочного тестирования..."
        echo -e "${CYAN}Параметры: 12 потоков, 400 соединений, 30 секунд${NC}\n"
        wrk -t6 -c100 -d30s --latency http://localhost:8080/api/simple
    elif command -v ab &> /dev/null; then
        echo "wrk не найден, использую Apache Bench..."
        ab -n 1000 -c 50 -q http://localhost:8080/api/simple
    else
        echo "wrk и Apache Bench не найдены, использую curl..."
        for i in {1..100}; do
            for j in {1..10}; do
                curl -s http://localhost:8080/api/simple > /dev/null &
            done
            wait
            echo -ne "\rПрогресс: $i/100 батчей"
        done
        echo
    fi

    echo -e "\n${GREEN}Нагрузочный тест завершён${NC}"
    echo -e "${CYAN}Проверьте метрики: http://localhost:8080/metrics${NC}"
}

# Custom test
custom_test() {
    print_subheader "Настраиваемый тест"

    echo -n "Введите endpoint (например, /api/simple): "
    read endpoint

    echo -n "Количество запросов: "
    read count

    echo -n "Задержка между запросами (мс): "
    read delay

    echo -n "Параллельные запросы? (y/n): "
    read parallel

    echo

    if [ "$parallel" = "y" ]; then
        echo "Запускаю $count параллельных запросов к $endpoint..."
        for i in $(seq 1 $count); do
            {
                response=$(curl -s -w "\n%{http_code}" http://localhost:8080$endpoint 2>/dev/null)
                http_code=$(echo "$response" | tail -n1)
                echo "Request #$i: $http_code"
            } &
        done
        wait
    else
        echo "Запускаю $count последовательных запросов к $endpoint..."
        for i in $(seq 1 $count); do
            response=$(curl -s -w "\n%{http_code}" http://localhost:8080$endpoint 2>/dev/null)
            http_code=$(echo "$response" | tail -n1)
            echo "Request #$i: $http_code"
            sleep "0.$(printf '%03d' $delay)"
        done
    fi
}

# Запуск всех тестов
run_all_tests() {
    print_header "🧪 Запуск всех тестов"

    go test -v ./integration_test.go -run "Demo"
}

# Главный цикл
main() {
    clear
    print_header "Welcome to Rate Limiting & Retry Demo! 🚀"

    # Проверяем сервер
    check_server
    echo

    while true; do
        show_menu
        read choice

        case $choice in
            1)
                print_header "🚦 Rate Limiting - Sequential Test"
                test_rate_limiting_sequential
                ;;
            2)
                print_header "🔥 Rate Limiting - Parallel Burst"
                test_rate_limiting_parallel
                ;;
            3)
                print_header "⏱ Sliding Window Test"
                test_sliding_window
                ;;
            4)
                print_header "💥 Flaky Service Test"
                test_flaky_service
                ;;
            5)
                print_header "🔄 Circuit Breaker Test"
                test_circuit_breaker
                ;;
            6)
                run_all_tests
                ;;
            7)
                print_header "📊 Load Test with Monitoring"
                load_test_with_monitoring
                ;;
            8)
                print_header "🎯 Custom Test"
                custom_test
                ;;
            0)
                echo -e "\n${GREEN}Спасибо за использование демо! До свидания! 👋${NC}\n"
                if [ ! -z "$SERVER_PID" ]; then
                    echo -e "${YELLOW}Останавливаю сервер...${NC}"
                    kill $SERVER_PID 2>/dev/null
                fi
                exit 0
                ;;
            *)
                echo -e "${RED}Неверный выбор. Попробуйте снова.${NC}"
                ;;
        esac

        echo
        echo -e "${CYAN}Нажмите Enter для продолжения...${NC}"
        read
    done
}

# Запуск
main