package dto

import (
	"time"

	"github.com/vctrl/currency-service/pkg/currency"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type CurrencyRequestDTO struct {
	BaseCurrency   string
	TargetCurrency string
	DateFrom       time.Time
	DateTo         time.Time
}

type CurrencyResponseDTO struct {
	Currency string
	Rates    []RateRecordDTO
}

type RateRecordDTO struct {
	Date time.Time
	Rate float32
}

// Должно быть в handler
func CurrencyRequestDTOFromProtobuf(req *currency.RateRequest, baseCurrency string) *CurrencyRequestDTO {
	return &CurrencyRequestDTO{
		BaseCurrency:   baseCurrency,
		TargetCurrency: req.Currency,
		DateFrom:       req.DateFrom.AsTime(),
		DateTo:         req.DateTo.AsTime(),
	}
}

// Должно быть в handler
func (dto *CurrencyResponseDTO) ToProtobuf() *currency.RateResponse {
	rateRecords := make([]*currency.RateRecord, 0, len(dto.Rates))
	for _, record := range dto.Rates {
		rateRecords = append(
			rateRecords, &currency.RateRecord{
				Date: timestamppb.New(record.Date),
				Rate: record.Rate,
			},
		)
	}

	return &currency.RateResponse{
		Currency: dto.Currency,
		Rates:    rateRecords,
	}
}
