package main

import (
	"errors"
	"time"

	"go.uber.org/zap"
)

type User struct {
	ID    int    `json:"id"`
	name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	demonstrateBasicFields(logger)
	demonstrateComplexFields(logger)
	demonstrateErrorLogging(logger)
}

func demonstrateBasicFields(logger *zap.Logger) {
	logger.Info("Базовые типы полей",
		zap.String("name", "Иван"),
		zap.Int("age", 30),
		zap.Bool("active", true),
		zap.Float64("balance", 1234.56),
	)

	logger.Info("Временные поля",
		zap.Time("created_at", time.Now()),
		zap.Duration("response_time", 250*time.Millisecond),
	)
}

func demonstrateComplexFields(logger *zap.Logger) {
	user := User{
		ID:    123,
		name:  "Алексей",
		Email: "alexey@example.com",
	}

	logger.Info("Пользователь зарегистрирован",
		zap.Any("user", user),
		zap.Strings("roles", []string{"user", "moderator"}),
		zap.Ints("permissions", []int{1, 2, 5}),
	)
}

func demonstrateErrorLogging(logger *zap.Logger) {
	err := errors.New("подключение к базе данных потеряно")

	logger.Error("Критическая ошибка",
		zap.Error(err),
		zap.String("database", "users_db"),
		zap.Int("retry_count", 3),
		zap.Bool("will_retry", true),
	)
}
