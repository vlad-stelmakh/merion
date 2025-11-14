package main

import (
	"go.uber.org/zap"
)

func main() {
	globalExample()
}

func globalExample() {
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))

	zap.L().Info("Использование глобального логгера")
	zap.S().Infow("Sugared logger для удобства",
		"ключ1", "значение1",
		"ключ2", 42,
	)
}
