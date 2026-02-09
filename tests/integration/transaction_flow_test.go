package integration

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/animesh/transaction_routine/internal/api"
	"github.com/animesh/transaction_routine/internal/domain"
	"github.com/animesh/transaction_routine/internal/repository/postgres"
	"github.com/animesh/transaction_routine/internal/service"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TransactionIntegrationTestSuite struct {
	suite.Suite
	pgContainer *tcpostgres.PostgresContainer
	db          *sql.DB
	router      *gin.Engine
}

func (s *TransactionIntegrationTestSuite) SetupSuite() {
	ctx := context.Background()

	// 1. Start Postgres Container
	pgContainer, err := tcpostgres.Run(ctx,
		"postgres:15-alpine",
		tcpostgres.WithDatabase("transaction_db"),
		tcpostgres.WithUsername("postgres"),
		tcpostgres.WithPassword("postgres"),
		tcpostgres.WithInitScripts("../../scripts/init.sql"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	s.NoError(err)
	s.pgContainer = pgContainer

	// 2. Get Connection String
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	s.NoError(err)

	// 3. Connect to DB
	db, err := sql.Open("postgres", connStr)
	s.NoError(err)
	s.db = db

	// 4. Setup App
	accountRepo := postgres.NewAccountRepository(db)
	txRepo := postgres.NewTransactionRepository(db)
	accountService := service.NewAccountService(accountRepo)
	transactionService := service.NewTransactionService(txRepo, accountRepo)
	handler := api.NewHandler(accountService, transactionService)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/accounts", handler.CreateAccount)
	r.GET("/accounts/:accountId", handler.GetAccount)
	r.POST("/transactions", handler.CreateTransaction)
	s.router = r
}

func (s *TransactionIntegrationTestSuite) TearDownSuite() {
	ctx := context.Background()
	if s.db != nil {
		s.db.Close()
	}
	if s.pgContainer != nil {
		s.pgContainer.Terminate(ctx)
	}
}

func (s *TransactionIntegrationTestSuite) SetupTest() {
	// Clean up tables before each test for isolation
	_, err := s.db.Exec("TRUNCATE TABLE transactions, accounts RESTART IDENTITY CASCADE")
	s.NoError(err)
}

func (s *TransactionIntegrationTestSuite) TestFullTransactionFlow() {
	// 1. Create Account
	accountReq := map[string]string{"document_number": "12345678900"}
	body, _ := json.Marshal(accountReq)
	req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusCreated, w.Code)

	var account domain.Account
	json.Unmarshal(w.Body.Bytes(), &account)
	s.Equal("12345678900", account.DocumentNumber)
	s.NotZero(account.ID)

	// 2. Create Transaction (Purchase - should be negative)
	txReq := map[string]interface{}{
		"account_id":        account.ID,
		"operation_type_id": 1, // Normal Purchase
		"amount":            100.50,
	}
	body, _ = json.Marshal(txReq)
	req, _ = http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusCreated, w.Code)

	var tx domain.Transaction
	json.Unmarshal(w.Body.Bytes(), &tx)
	s.Equal(-100.50, tx.Amount)
	s.Equal(account.ID, tx.AccountID)

	// 3. Verify in Database
	var dbAmount float64
	err := s.db.QueryRow("SELECT amount FROM transactions WHERE transaction_id = $1", tx.ID).Scan(&dbAmount)
	s.NoError(err)
	s.Equal(-100.50, dbAmount)
}

func (s *TransactionIntegrationTestSuite) TestAccountNotFound() {
	txReq := map[string]interface{}{
		"account_id":        999, // Non-existent
		"operation_type_id": 1,
		"amount":            100.0,
	}
	body, _ := json.Marshal(txReq)
	req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *TransactionIntegrationTestSuite) TestInvalidPayload() {
	// Missing Amount
	txReq := map[string]interface{}{
		"account_id":        1,
		"operation_type_id": 1,
	}
	body, _ := json.Marshal(txReq)
	req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *TransactionIntegrationTestSuite) TestInvalidOperationType() {
	// Create Account first
	accountReq := map[string]string{"document_number": "99999"}
	body, _ := json.Marshal(accountReq)
	req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusCreated, w.Code)

	var account domain.Account
	json.Unmarshal(w.Body.Bytes(), &account)

	txReq := map[string]interface{}{
		"account_id":        account.ID,
		"operation_type_id": 99, // Invalid
		"amount":            100.0,
	}
	body, _ = json.Marshal(txReq)
	req, _ = http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *TransactionIntegrationTestSuite) TestConcurrency() {
	// Create Account
	accountReq := map[string]string{"document_number": "88888"}
	body, _ := json.Marshal(accountReq)
	req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusCreated, w.Code)

	var account domain.Account
	json.Unmarshal(w.Body.Bytes(), &account)

	// Run 10 concurrent transactions
	concurrency := 10
	done := make(chan bool)

	for i := 0; i < concurrency; i++ {
		go func() {
			txReq := map[string]interface{}{
				"account_id":        account.ID,
				"operation_type_id": 4, // Payment (Credit)
				"amount":            10.0,
			}
			body, _ := json.Marshal(txReq)
			req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
			rec := httptest.NewRecorder()
			s.router.ServeHTTP(rec, req)
			s.Equal(http.StatusCreated, rec.Code)
			done <- true
		}()
	}

	for i := 0; i < concurrency; i++ {
		<-done
	}

	// Verify total transactions
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM transactions WHERE account_id = $1", account.ID).Scan(&count)
	s.NoError(err)
	s.Equal(concurrency, count)
}

func TestTransactionIntegrationSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	suite.Run(t, new(TransactionIntegrationTestSuite))
}
