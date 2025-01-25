package handler

import (
	"context"
	"fmt"

	"github.com/vctrl/currency-service/currency/internal/dto"
	"github.com/vctrl/currency-service/pkg/currency"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// Должна быть в dto
const (
	defaultBaseCurrency = "RUB"
)

func (s CurrencyServer) GetRate(ctx context.Context, request *currency.RateRequest) (*currency.RateResponse, error) {
	reqDTO := dto.CurrencyRequestDTOFromProtobuf(request, defaultBaseCurrency)

	rates, err := s.service.GetCurrencyRatesInInterval(ctx, reqDTO)
	if err != nil {
		return nil, fmt.Errorf("service.GetCurrencyRatesInInterval: %w", err)
	}

	rateRecords := make([]*currency.RateRecord, len(rates))
	for i, rate := range rates {
		rateRecords[i] = &currency.RateRecord{
			Date: timestamppb.New(rate.Date),
			Rate: rate.Rate,
		}
	}

	return &currency.RateResponse{
		Currency: reqDTO.TargetCurrency,
		Rates:    rateRecords,
	}, nil
}
