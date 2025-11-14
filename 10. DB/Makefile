# Makefile для управления проектом

.PHONY: help install-goose migrate-up migrate-down migrate-status migrate-version migrate-reset migrate-create run build clean

# Показать справку
help:
	@echo "🔧 Доступные команды:"
	@echo ""
	@echo "  install-goose    - Установить goose в bin/"
	@echo "  migrate-up       - Применить все миграции"
	@echo "  migrate-down     - Откатить последнюю миграцию"
	@echo "  migrate-status   - Показать статус миграций"
	@echo "  migrate-version  - Показать версию БД"
	@echo "  migrate-reset    - Сбросить все миграции"
	@echo "  migrate-create   - Создать новую миграцию (make migrate-create NAME=имя_миграции)"
	@echo "  run             - Запустить приложение"
	@echo "  build           - Собрать приложение"
	@echo "  clean           - Очистить сборочные файлы"
	@echo ""

# Установить goose
install-goose:
	@echo "📦 Установка goose..."
	@mkdir -p bin
	@GOBIN=$(PWD)/bin go install github.com/pressly/goose/v3/cmd/goose@latest
	@echo "✅ Goose установлен в bin/goose"

# Миграции
migrate-up:
	@./scripts/migrate.sh up

migrate-down:
	@./scripts/migrate.sh down

migrate-status:
	@./scripts/migrate.sh status

migrate-version:
	@./scripts/migrate.sh version

migrate-reset:
	@./scripts/migrate.sh reset

migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "❌ Укажите имя миграции: make migrate-create NAME=имя_миграции"; \
		exit 1; \
	fi
	@./scripts/migrate.sh create $(NAME)

# Приложение
run:
	@echo "🚀 Запуск приложения..."
	@go run main.go

build:
	@echo "🔨 Сборка приложения..."
	@go build -o bin/app main.go
	@echo "✅ Приложение собрано в bin/app"

# Очистка
clean:
	@echo "🧹 Очистка..."
	@rm -f bin/app
	@echo "✅ Очистка завершена"

# Установка зависимостей
deps:
	@echo "📦 Установка зависимостей..."
	@go mod download
	@go mod tidy
	@echo "✅ Зависимости установлены"

# Форматирование кода
fmt:
	@echo "🎨 Форматирование кода..."
	@go fmt ./...
	@echo "✅ Код отформатирован"

# Проверка кода
vet:
	@echo "🔍 Проверка кода..."
	@go vet ./...
	@echo "✅ Проверка завершена"

# Тесты
test:
	@echo "🧪 Запуск тестов..."
	@go test ./...
	@echo "✅ Тесты завершены"

