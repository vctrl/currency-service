package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/vctrl/currency-service/currency/internal/dto"
	"github.com/vctrl/currency-service/currency/internal/pkg/currency"
	"github.com/vctrl/currency-service/currency/internal/repository"

	"go.uber.org/zap"
)

type CurrencyService struct {
	currencyRepo repository.ExchangeRateRepository
	client       currency.CurrencyClient
	logger       *zap.Logger
}

func NewCurrencyService(repo repository.ExchangeRateRepository,
	client currency.CurrencyClient,
	logger *zap.Logger) CurrencyService {
	return CurrencyService{
		currencyRepo: repo,
		client:       client,
		logger:       logger,
	}
}

func (s *CurrencyService) GetCurrencyRatesInInterval(ctx context.Context, reqDTO *dto.CurrencyRequestDTO) ([]repository.CurrencyRate, error) {
	reqDTO.TargetCurrency = strings.ToLower(reqDTO.TargetCurrency)
	rates, err := s.currencyRepo.FindInInterval(ctx, reqDTO)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch currency rates in interval: %w", err)
	}

	return rates, nil
}

func (s *CurrencyService) FetchAndSaveCurrencyRates(ctx context.Context, baseCurrency string) error {
	rates, err := s.client.FetchCurrentRates(ctx, baseCurrency)
	if err != nil {
		return fmt.Errorf("client.FetchCurrentRates: %s", err)
	}

	date, err := time.Parse("2006-01-02", rates.Date)
	if err != nil {
		return fmt.Errorf("Failed to parse currency date: %v ", err)
	}

	if err := s.currencyRepo.Save(ctx, date, baseCurrency, rates.Rub); err != nil { // todo want to pass struct
		return fmt.Errorf("Failed to save currency rates: %v ", err)
	}

	s.logger.Info("Currency rates fetched and saved", zap.Any("rates", rates))
	return nil
}
