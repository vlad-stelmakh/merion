package main

import (
	"context"
	"fmt"

	log "github.com/vostelmakh/ctxlog"
	"go.uber.org/zap"

	"ctxlog/service"
)

func main() {
	logger := zap.Must(zap.NewProduction())

	ctxLogger := log.New(logger)

	orderNr := `ORD-12345`

	ctx := context.Background()

	ctx = log.WithOrderNr(ctx, orderNr)
	ctx = log.WithMessageID(ctx, "MSG-1289")

	ctxLogger.Info(ctx, "Start handling request")

	orderService := service.NewOrderService(ctxLogger)

	err := orderService.UpdateOrderWithError(ctx, orderNr)
	if err != nil {
		errCtx := log.ErrorCtx(ctx, err)
		f := log.Fields(errCtx, nil)

		for k, v := range f {
			fmt.Println(k, v)
		}

		ctxLogger.Fatal(ctx, "Failed to update order", zap.Error(err))
	}
}
