package handler

import (
	"github.com/vctrl/currency-service/currency/internal/service"
	"github.com/vctrl/currency-service/pkg/currency"

	"go.uber.org/zap"
)

type CurrencyServer struct {
	currency.UnimplementedCurrencyServiceServer
	service service.CurrencyService
	logger  *zap.Logger
}

func NewCurrencyServer(svc service.CurrencyService, logger *zap.Logger) CurrencyServer {
	return CurrencyServer{
		service: svc,
		logger:  logger,
	}
}
