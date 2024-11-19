package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/vctrl/currency-service/currency/internal/pkg/config"
	"github.com/vctrl/currency-service/currency/internal/service"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type CurrencyWorker struct {
	CurrencyService service.CurrencyService
	Cron            *cron.Cron
	Schedule        string
	BaseCurrency    string
	TargetCurrency  string
	logger          *zap.Logger
}

func NewCurrencyWorker(cfg config.WorkerConfig,
	service service.CurrencyService,
	cron *cron.Cron,
	logger *zap.Logger) *CurrencyWorker {
	return &CurrencyWorker{
		CurrencyService: service,
		Cron:            cron,
		Schedule:        cfg.Schedule,
		BaseCurrency:    cfg.CurrencyPair.BaseCurrency,
		TargetCurrency:  cfg.CurrencyPair.TargetCurrency,
		logger:          logger,
	}
}

func (w *CurrencyWorker) StartFetchingCurrencyRates() error {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5) // todo move to config
		defer cancel()

		err := w.CurrencyService.FetchAndSaveCurrencyRates(ctx, w.BaseCurrency)

		if err != nil {
			w.logger.Error(
				"Failed to fetch currency rate immediately on startup",
				zap.Time("timestamp", time.Now()),
				zap.Error(err),
			)
		}
	}()

	_, err := w.Cron.AddFunc(
		w.Schedule, func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5) // todo move to config
			defer cancel()

			err := w.CurrencyService.FetchAndSaveCurrencyRates(ctx, w.BaseCurrency)
			if err != nil {
				w.logger.Error(
					"Failed to fetch currency rate on scheduled run",
					zap.Time("timestamp", time.Now()),
					zap.Error(err),
					zap.String("schedule", w.Schedule),
				)
			}
		},
	)

	if err != nil {
		return fmt.Errorf("Cron.AddFunc: %w", err)
	}

	w.Cron.Start()

	return nil
}

func (w *CurrencyWorker) Stop() {
	w.Cron.Stop()
}
