package main

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(file),
			zap.InfoLevel,
		),
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.AddSync(os.Stdout),
			zap.DebugLevel,
		),
	)

	logger := zap.New(core)
	defer logger.Sync()

	logger.Debug("Отладочная информация (только в консоли)")
	logger.Info("Информационное сообщение (в файл и консоль)")
	logger.Error("Ошибка (в файл и консоль)")
}
