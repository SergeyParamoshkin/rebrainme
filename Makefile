# ============================================
# Makefile для демонстрации Docker best practices
# ============================================

.PHONY: help build-all compare clean run-demo test-all

# Переменные
IMAGE_PREFIX := golang-demo
VERSION := 1.0.0
DOCKERFILES_DIR := dockerfiles
SCRIPTS_DIR := scripts

# Цвета для вывода
CYAN := \033[0;36m
GREEN := \033[0;32m
YELLOW := \033[1;33m
NC := \033[0m

# По умолчанию показываем help
.DEFAULT_GOAL := help

help: ## Показать список доступных команд
	@echo "$(CYAN)=====================================$(NC)"
	@echo "$(CYAN)  Docker Best Practices - Commands  $(NC)"
	@echo "$(CYAN)=====================================$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""

build-all: ## Собрать все варианты Dockerfile
	@echo "$(YELLOW)Building all Docker images...$(NC)"
	@chmod +x $(SCRIPTS_DIR)/build-all.sh
	@$(SCRIPTS_DIR)/build-all.sh

build-basic: ## Собрать базовый вариант
	@echo "$(YELLOW)Building basic image...$(NC)"
	@docker build -t $(IMAGE_PREFIX):basic -f $(DOCKERFILES_DIR)/Dockerfile.basic ..

build-optimized: ## Собрать оптимизированный вариант
	@echo "$(YELLOW)Building optimized image...$(NC)"
	@docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_TIME=$$(date -u +%Y-%m-%dT%H:%M:%SZ) \
		-t $(IMAGE_PREFIX):optimized \
		-f $(DOCKERFILES_DIR)/Dockerfile.optimized ..

build-distroless: ## Собрать distroless вариант
	@echo "$(YELLOW)Building distroless image...$(NC)"
	@docker build -t $(IMAGE_PREFIX):distroless -f $(DOCKERFILES_DIR)/Dockerfile.distroless ..

build-scratch: ## Собрать scratch вариант
	@echo "$(YELLOW)Building scratch image...$(NC)"
	@docker build -t $(IMAGE_PREFIX):scratch -f $(DOCKERFILES_DIR)/Dockerfile.scratch ..

build-advanced: ## Собрать продвинутый multi-stage вариант
	@echo "$(YELLOW)Building advanced multi-stage image...$(NC)"
	@docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_TIME=$$(date -u +%Y-%m-%dT%H:%M:%SZ) \
		--build-arg GIT_COMMIT=$$(git rev-parse --short HEAD) \
		--build-arg SKIP_TESTS=true \
		--build-arg SKIP_LINT=true \
		-t $(IMAGE_PREFIX):advanced \
		-f $(DOCKERFILES_DIR)/Dockerfile.multistage-advanced ..

compare: ## Сравнить размеры образов
	@chmod +x $(SCRIPTS_DIR)/compare-sizes.sh
	@$(SCRIPTS_DIR)/compare-sizes.sh $(IMAGE_PREFIX)

list: ## Показать список собранных образов
	@echo "$(YELLOW)Docker images:$(NC)"
	@docker images $(IMAGE_PREFIX) --format "table {{.Repository}}:{{.Tag}}\t{{.Size}}\t{{.CreatedSince}}"

inspect-basic: ## Детальная информация о basic образе
	@docker inspect $(IMAGE_PREFIX):basic | jq '.[0] | {Size: .Size, Architecture: .Architecture, Os: .Os, Created: .Created}'

inspect-optimized: ## Детальная информация о optimized образе
	@docker inspect $(IMAGE_PREFIX):optimized | jq '.[0] | {Size: .Size, Architecture: .Architecture, Os: .Os, Created: .Created}'

run-basic: ## Запустить basic контейнер
	@echo "$(YELLOW)Running basic container on port 8081...$(NC)"
	@docker run --rm -p 8081:8080 --name demo-basic $(IMAGE_PREFIX):basic

run-optimized: ## Запустить optimized контейнер
	@echo "$(YELLOW)Running optimized container on port 8082...$(NC)"
	@docker run --rm -p 8082:8080 --name demo-optimized $(IMAGE_PREFIX):optimized

run-demo: ## Запустить все контейнеры через docker-compose
	@echo "$(YELLOW)Starting demo environment...$(NC)"
	@docker-compose -f docker-compose.demo.yml up -d
	@echo "$(GREEN)Demo environment started!$(NC)"
	@echo ""
	@echo "Services available at:"
	@echo "  • app-basic:      http://localhost:8081"
	@echo "  • app-optimized:  http://localhost:8082"
	@echo "  • app-distroless: http://localhost:8083"
	@echo "  • app-scratch:    http://localhost:8084"
	@echo "  • app-advanced:   http://localhost:8085"
	@echo "  • Adminer:        http://localhost:8080"
	@echo "  • Portainer:      http://localhost:9000"

stop-demo: ## Остановить demo environment
	@echo "$(YELLOW)Stopping demo environment...$(NC)"
	@docker-compose -f docker-compose.demo.yml down

logs-demo: ## Показать логи demo environment
	@docker-compose -f docker-compose.demo.yml logs -f

clean-images: ## Удалить все собранные образы
	@echo "$(YELLOW)Removing all $(IMAGE_PREFIX) images...$(NC)"
	@docker images $(IMAGE_PREFIX) -q | xargs -r docker rmi -f
	@echo "$(GREEN)Images removed!$(NC)"

clean-all: clean-images ## Полная очистка (образы + volumes)
	@echo "$(YELLOW)Removing all containers and volumes...$(NC)"
	@docker-compose -f docker-compose.demo.yml down -v
	@echo "$(GREEN)Cleanup complete!$(NC)"

security-scan: ## Запустить security scan с trivy (требует установки trivy)
	@echo "$(YELLOW)Scanning images for vulnerabilities...$(NC)"
	@command -v trivy >/dev/null 2>&1 || { echo "$(RED)trivy not installed. Install from https://github.com/aquasecurity/trivy$(NC)"; exit 1; }
	@for tag in basic optimized distroless scratch advanced; do \
		echo "$(CYAN)Scanning $(IMAGE_PREFIX):$$tag...$(NC)"; \
		trivy image $(IMAGE_PREFIX):$$tag; \
	done

dive-basic: ## Анализ слоев basic образа (требует dive)
	@command -v dive >/dev/null 2>&1 || { echo "$(RED)dive not installed. Install from https://github.com/wagoodman/dive$(NC)"; exit 1; }
	@dive $(IMAGE_PREFIX):basic

dive-optimized: ## Анализ слоев optimized образа
	@command -v dive >/dev/null 2>&1 || { echo "$(RED)dive not installed. Install from https://github.com/wagoodman/dive$(NC)"; exit 1; }
	@dive $(IMAGE_PREFIX):optimized

test-all: build-all ## Собрать и протестировать все образы
	@echo "$(YELLOW)Testing all images...$(NC)"
	@for tag in basic optimized distroless scratch advanced; do \
		echo "$(CYAN)Testing $(IMAGE_PREFIX):$$tag...$(NC)"; \
		docker run --rm $(IMAGE_PREFIX):$$tag --version || echo "No version command available"; \
	done
	@echo "$(GREEN)All tests completed!$(NC)"

stats: ## Показать статистику образов
	@echo "$(YELLOW)Docker images statistics:$(NC)"
	@docker images $(IMAGE_PREFIX) --format "{{.Repository}}:{{.Tag}}" | while read img; do \
		echo "$(CYAN)$$img:$(NC)"; \
		docker inspect $$img --format='  Size: {{.Size}} bytes ({{div .Size 1048576}} MB)'; \
		docker inspect $$img --format='  Layers: {{len .RootFS.Layers}}'; \
		docker inspect $$img --format='  Created: {{.Created}}'; \
		echo ""; \
	done

# Установка необходимых инструментов
install-tools: ## Установить полезные инструменты для работы с Docker
	@echo "$(YELLOW)Installing Docker tools...$(NC)"
	@echo "Please install manually:"
	@echo "  • dive: https://github.com/wagoodman/dive"
	@echo "  • trivy: https://github.com/aquasecurity/trivy"
	@echo "  • hadolint: https://github.com/hadolint/hadolint"
	@echo ""
	@echo "MacOS:"
	@echo "  brew install dive trivy hadolint"
	@echo ""
	@echo "Linux:"
	@echo "  See respective project documentation"
