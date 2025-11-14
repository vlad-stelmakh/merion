# Database Examples

Примеры работы с базой данных PostgreSQL на Go с использованием лучших практик.

## Архитектура проекта

Проект организован по следующей структуре:

```
.
├── bin/                    # Бинарные файлы
│   └── goose              # Goose для миграций
├── internal/              # Внутренние пакеты
│   ├── config/           # Конфигурация
│   ├── models/           # Модели данных
│   └── repository/       # Репозитории
│       └── postgres/     # PostgreSQL реализация
├── migrations/           # SQL миграции
├── scripts/              # Скрипты
│   └── migrate.sh        # Скрипт для миграций
├── main.go              # Основное приложение
├── Makefile             # Команды для разработки
├── go.mod               # Go модули
└── README.md           # Документация
```

## Особенности

- **Паттерн Repository**: Абстракция для работы с данными
- **Goose миграции**: Управление схемой базы данных
- **Конфигурация через ENV**: Гибкая настройка подключения
- **Контекст**: Поддержка отмены операций
- **Транзакции**: Безопасные операции с данными
- **Пул соединений**: Оптимизированная работа с БД

## Установка и запуск

### 1. Установка зависимостей

```bash
# Установка Go зависимостей
make deps

# Установка goose (если не установлен)
make install-goose
```

### 2. Настройка базы данных

Создайте файл `.env` в корне проекта:

```env
# Настройки базы данных
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=testdb
DB_SSLMODE=disable

# Настройки пула соединений
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=25
DB_CONN_MAX_LIFETIME=5m
```

### 3. Запуск миграций

```bash
# Применить все миграции
make migrate-up

# Посмотреть статус миграций
make migrate-status

# Откатить последнюю миграцию
make migrate-down

# Сбросить все миграции
make migrate-reset

# Создать новую миграцию
make migrate-create NAME=add_user_avatar
```

### 4. Запуск приложения

```bash
# Запуск через make
make run

# Или напрямую
go run main.go

# Сборка приложения
make build
```

## Использование

### Работа с пользователями

```go
// Создание пользователя
userReq := &models.CreateUserRequest{
    Name:     "Иван Иванов",
    Email:    "ivan@example.com",
    Age:      sql.NullInt64{Int64: 30, Valid: true},
    Bio:      sql.NullString{String: "Разработчик", Valid: true},
    IsActive: true,
}

user, err := repo.User.Create(ctx, userReq)

// Получение пользователя
user, err := repo.User.GetByID(ctx, 1)

// Обновление пользователя
updateReq := &models.UpdateUserRequest{
    Name: &newName,
    Age:  &newAge,
}
updatedUser, err := repo.User.Update(ctx, 1, updateReq)

// Удаление пользователя
err := repo.User.Delete(ctx, 1)
```


## Миграции

Миграции находятся в папке `migrations/` и используют формат goose:

```sql
-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE users;
```

## Лучшие практики

1. **Использование контекста**: Все операции с БД принимают `context.Context`
2. **Обработка NULL значений**: Использование `sql.NullString`, `sql.NullInt64`
3. **Транзакции**: Критические операции выполняются в транзакциях
4. **Пул соединений**: Настроен для оптимальной производительности
5. **Структурированные ошибки**: Понятные сообщения об ошибках
6. **Разделение ответственности**: Четкое разделение на слои

## Команды для разработки

### Makefile команды

```bash
# Показать все доступные команды
make help

# Установка и настройка
make deps              # Установить зависимости
make install-goose     # Установить goose

# Миграции
make migrate-up        # Применить миграции
make migrate-down      # Откатить последнюю миграцию
make migrate-status    # Статус миграций
make migrate-create NAME=имя_миграции  # Создать новую миграцию

# Разработка
make run              # Запустить приложение
make build            # Собрать приложение
make clean            # Очистить сборочные файлы

# Качество кода
make fmt              # Форматировать код
make vet              # Проверить код
make test             # Запустить тесты
```

### Прямые команды

```bash
# Форматирование кода
go fmt ./...

# Проверка кода
go vet ./...

# Запуск тестов
go test ./...

# Работа с миграциями через скрипт
./scripts/migrate.sh help
```

## Структура базы данных

### Таблица users
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    age INTEGER,
    bio TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```


## Пример вывода

```
✅ Успешно подключились к PostgreSQL!

🔄 Демонстрация операций с пользователями:

📝 Создание пользователей...
✅ Создан пользователь с ID: 1
✅ Создан пользователь с ID: 2

📖 Чтение пользователя по ID...
✅ Найден пользователь: Алексей Иванов (alexey@example.com)
   Возраст: 28
   Биография: Разработчик Go

✏️  Обновление пользователя...
✅ Пользователь обновлен: Алексей Иванович Иванов, возраст: 29

🗑️  Удаление пользователя...
✅ Пользователь удален

📊 Всего пользователей в системе: 1

📊 Статистика пула соединений:
   Открытых соединений: 1
   Используемых соединений: 0
   Простаивающих соединений: 1
```