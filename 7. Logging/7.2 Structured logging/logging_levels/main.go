package main

import (
	"go.uber.org/zap"
)

func main() {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Debug("Это Debug сообщение не будет показано")
	logger.Info("Info: Это сообщение будет показано")
	logger.Warn("Warn: И это предупреждение тоже")
	logger.Error("Error: И эта ошибка")
}
