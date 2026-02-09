package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/animesh/transaction_routine/internal/api"
	"github.com/animesh/transaction_routine/internal/config"
	"github.com/animesh/transaction_routine/internal/infrastructure/observability"
	"github.com/animesh/transaction_routine/internal/repository/postgres"
	"github.com/animesh/transaction_routine/internal/service"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	// 1. Initialize Configuration
	cfg := config.Load()

	// 2. Initialize Logger
	observability.InitLogger()
	logger := observability.GetLogger()
	defer logger.Sync()

	logger.Info("Starting Transaction Routine API", zap.String("port", cfg.AppPort))

	// 3. Database Connection
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Fatal("could not connect to database", zap.Error(err))
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.Fatal("could not ping database", zap.Error(err))
	}

	// 4. Initialize Dependency Injection
	accountRepo := postgres.NewAccountRepository(db)
	txRepo := postgres.NewTransactionRepository(db)

	accountService := service.NewAccountService(accountRepo)
	transactionService := service.NewTransactionService(txRepo, accountRepo)

	handler := api.NewHandler(accountService, transactionService)

	// 5. Router Setup
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// Structured logging middleware
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Info("Request processed",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
		)
	})

	r.POST("/accounts", handler.CreateAccount)
	r.GET("/accounts/:accountId", handler.GetAccount)
	r.POST("/transactions", handler.CreateTransaction)

	// 6. Server Setup with Graceful Shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen: %s\n", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: ", zap.Error(err))
	}

	logger.Info("Server exiting")
}
