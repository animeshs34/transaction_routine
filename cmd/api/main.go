package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	api "github.com/animeshs34/transaction_routine/internal/api"
	"github.com/animeshs34/transaction_routine/internal/config"
	"github.com/animeshs34/transaction_routine/internal/logger"
	"github.com/animeshs34/transaction_routine/internal/respository"
	"github.com/animeshs34/transaction_routine/internal/service"
	"go.uber.org/zap"
)

func main() {
	configFile := flag.String("config", "", "config file path")
	flag.Parse()
	var cfg *config.Config
	var err error

	if *configFile != "" {
		cfg, err = config.LoadFromFile(*configFile)
		if err != nil {
			fmt.Printf("Failed to load configuration from file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Loaded configuration from file: %s\n", *configFile)
	} else {
		cfg, err = config.Load()
		if err != nil {
			fmt.Printf("Failed to load configuration: %v\n", err)
			os.Exit(1)
		}
	}

	if err := logger.Init(cfg.Logging.Level); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	var repo respository.Respository
	var dbConn *respository.DBConn
	switch cfg.Database.Type {
	case "memory":
		repo = respository.NewInMemoryStore()
	case "postgres":
		var err error
		dbConn, err = respository.NewPostgresConn(
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.DBName,
			cfg.Database.SSLMode,
		)
		if err != nil {
			logger.Fatal("Failed to initialize PostgresStore connection", zap.Error(err))
		}
		repo = respository.NewPostgresStore(dbConn)
	default:
		logger.Fatal("Unsupported database type", zap.String("type", cfg.Database.Type))
	}

	svc := service.New(repo)
	handler := api.New(svc)

	middlewareChainedHandler := api.Chain(
		handler.Router(),
		api.Recoverer(),
		api.LoggingMiddleware,
	)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      middlewareChainedHandler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		logger.Info("HTTP server starting", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Server error", zap.Error(err))
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Graceful shutdown failed", zap.Error(err))
	}

	if dbConn != nil {
		if err := dbConn.Close(); err != nil {
			logger.Error("Failed to close PostgresStore connection", zap.Error(err))
		} else {
			logger.Info("PostgresStore connection closed")
		}
	}

	logger.Info("Server stopped")
}
