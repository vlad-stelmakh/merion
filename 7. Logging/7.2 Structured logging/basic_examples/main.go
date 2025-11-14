package main

import (
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	logger.Info("Простое информационное сообщение")
	logger.Warn("Предупреждение о потенциальной проблеме")
	logger.Error("Произошла ошибка")

	logger.Info("Сообщение с полями",
		zap.String("user_id", "12356"),
		zap.Int("версия", 1),
	)

	devLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	devLogger.Info("Тот же логгер для разработки", zap.String("user_id", "12356"))
	devLogger.Debug("Отладочное сообщение видно только в dev режиме")
}
