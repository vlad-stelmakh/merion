package service

import (
	"context"
	"errors"

	log "github.com/vostelmakh/ctxlog"
)

type OrderService struct {
	logger log.Logger
}

func NewOrderService(logger log.Logger) *OrderService {
	return &OrderService{
		logger: logger,
	}
}

func (s *OrderService) UpdateOrder(ctx context.Context, orderNr string) error {
	s.logger.Info(ctx, "Start updating order")

	/**
	 */

	s.logger.Info(ctx, "Order updated")

	s.logger.Info(ctx, "Finish updating order")

	return nil
}

func (s *OrderService) UpdateOrderWithError(ctx context.Context, orderNr string) error {
	s.logger.Info(ctx, "Start updating order")

	return log.WrapErrorCtx(ctx, errors.New("cannot update order"))
}
