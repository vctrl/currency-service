package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/vctrl/currency-service/currency/internal/pkg/config"
	currencyClient "github.com/vctrl/currency-service/currency/internal/pkg/currency"
	"github.com/vctrl/currency-service/currency/internal/repository"
	"github.com/vctrl/currency-service/currency/internal/service"
	"github.com/vctrl/currency-service/currency/internal/worker"

	"go.uber.org/zap"

	"github.com/robfig/cron/v3"
)

func main() {
	configPath := flag.String("config", "./config", "path to the config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	db, _, err := repository.NewDatabaseConnection(cfg.Database)
	if err != nil {
		log.Fatalf("error init database connection: %v", err)
	}
	//if err := migrations.RunPgMigrations(dsn, cfg.Database.MigrationsPath); err != nil {
	//	log.Fatalf("RunPgMigrations failed: %v", err)
	//}

	repo, err := repository.NewExchangeRateRepository(db)
	if err != nil {
		log.Fatalf("error creating repository: %v", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("init logger: %w", err)
	}

	client, err := currencyClient.NewCurrencyClient(cfg.API, logger)
	if err != nil {
		log.Fatalf("error creating currency client: %v", err)
	}

	svc := service.NewCurrencyService(repo, client, logger)

	c := cron.New()

	currencyWorker := worker.NewCurrencyWorker(cfg.Worker, svc, c, logger)

	if err != nil {
		log.Fatalf("error adding cron job: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := currencyWorker.StartFetchingCurrencyRates(); err != nil {
		log.Fatalf("error start fetching currency rates: %v", err)
	}

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	currencyWorker.Stop()
}
