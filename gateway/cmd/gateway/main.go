package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/vctrl/currency-service/gateway/internal/config"
	"github.com/vctrl/currency-service/gateway/internal/handler"
	"github.com/vctrl/currency-service/gateway/internal/middleware"
	"github.com/vctrl/currency-service/gateway/internal/pkg/auth"
	"github.com/vctrl/currency-service/gateway/internal/repository"
	"github.com/vctrl/currency-service/gateway/internal/service"
	"github.com/vctrl/currency-service/pkg/grpc_client"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	configPath := flag.String("config", "./config", "path to the config file")

	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	fmt.Println(cfg)

	router := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("init logger: %w", err)
	}

	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))

	// auth middleware
	authClient, err := auth.NewAuthClient(cfg.Auth)
	if err != nil {
		log.Fatalf("auth.NewAuthClient: %s", err)
	}

	if resp, err := authClient.Ping(); err != nil {
		log.Fatalf("authClient.Ping: %s", err)
	} else {
		if resp != "pong" {
			log.Fatalf("auth client answered with invalid response: %s", resp)
		}
	}

	authMiddleware := middleware.NewAuthorization(authClient, shouldSkipAuthMiddleware, logger)
	router.Use(authMiddleware.Authorize())

	currencyClient, conn, err := grpc_client.NewCurrencyServiceClient(cfg.GRPC.CurrencyServiceURL)
	if err != nil {
		log.Fatalf("can't create currency service client: %v", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			// todo
		}
	}()

	userRepo := repository.NewUserRepository()
	authService := service.NewAuthService(authClient, userRepo)
	currencyService := service.NewCurrencyService(currencyClient)

	srv := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: router,
	}

	server := handler.NewServer(authService, currencyService, router, logger)
	server.RegisterRoutes()

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("srv.ListenAndServe: %s\n", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

func shouldSkipAuthMiddleware(c *gin.Context) bool {
	if strings.HasSuffix(c.Request.URL.Path, "/login") ||
		strings.HasSuffix(c.Request.URL.Path, "/register") {
		return true
	}

	return false
}
