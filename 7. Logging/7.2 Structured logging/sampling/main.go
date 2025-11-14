package main

import (
	"go.uber.org/zap"
)

func main() {
	config := zap.NewProductionConfig()

	config.Sampling = &zap.SamplingConfig{
		Initial:    100,
		Thereafter: 100,
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	for i := 0; i < 1000; i++ {
		logger.Info("Тестовое сообщение для демонстрации выборки",
			zap.Int("iteration", i),
		)
	}

	logger.Info("Выборка логов помогает уменьшить объем в высоконагруженных системах")
}
