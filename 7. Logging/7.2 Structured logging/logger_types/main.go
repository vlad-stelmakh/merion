package main

import (
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	usr := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "str",
		Age:  123,
	}

	logger.Info("Production логгер",
		zap.String("mode", "production"),
		zap.String("format", "JSON"),
		zap.Bool("structured", true),
		zap.Any("user", usr),
	)

	logger.Error("Ошибка в production режиме",
		zap.String("error", "database connection failed"),
		zap.Int("retry_count", 3),
	)
}
