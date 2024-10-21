package main

import (
	"github.com/vctrl/currency-service/currency/internal/handler"
	"github.com/vctrl/currency-service/currency/internal/pkg/config"
	currencyClient "github.com/vctrl/currency-service/currency/internal/pkg/currency"
	"github.com/vctrl/currency-service/currency/internal/repository"
	"github.com/vctrl/currency-service/currency/internal/service"
	"github.com/vctrl/currency-service/pkg/currency"

	"flag"
	"fmt"
	"log"
	"net"

	"go.uber.org/zap"

	"google.golang.org/grpc"
)

func main() {
	configPath := flag.String("config", "./config", "path to the config file")

	flag.Parse()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("init logger: %v", err)
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	db, _, err := repository.NewDatabaseConnection(cfg.Database)
	//if err := migrations.RunPgMigrations(dsn, cfg.Database.MigrationsPath); err != nil {
	//	log.Fatalf("RunPgMigrations failed: %v", err)
	//}

	repo, err := repository.NewExchangeRateRepository(db)
	if err != nil {
		log.Fatalf("error init exchange rate repository: %v", err)
	}

	client, err := currencyClient.NewCurrencyClient(cfg.API, logger)
	if err != nil {
		log.Fatalf("error creating currency client: %v", err)
	}

	svc := service.NewCurrencyService(repo, client, logger)

	currencyServer := handler.NewCurrencyServer(svc, logger)

	if err := startGRPCServer(cfg, currencyServer); err != nil {
		log.Fatalf("Error starting GRPC server: %s", err)
	}
}

func startGRPCServer(cfg config.AppConfig, srv handler.CurrencyServer) error {
	lis, err := net.Listen("tcp", ":"+cfg.Service.ServerPort)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s := grpc.NewServer()
	currency.RegisterCurrencyServiceServer(s, srv)

	log.Printf("gRPC server is listening on :%s", cfg.Service.ServerPort)
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}
