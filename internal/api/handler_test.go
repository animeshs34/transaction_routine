package api

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/animesh/transaction_routine/internal/domain"
	customErrors "github.com/animesh/transaction_routine/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks
type MockAccountService struct {
	mock.Mock
}

func (m *MockAccountService) CreateAccount(ctx context.Context, documentNumber string) (*domain.Account, error) {
	args := m.Called(ctx, documentNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Account), args.Error(1)
}

func (m *MockAccountService) GetAccount(ctx context.Context, id int64) (*domain.Account, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Account), args.Error(1)
}

type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) CreateTransaction(ctx context.Context, accountID int64, operationTypeID int, amount float64) (*domain.Transaction, error) {
	args := m.Called(ctx, accountID, operationTypeID, amount)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Transaction), args.Error(1)
}

func TestHandler_CreateAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockAccUC := new(MockAccountService)
		mockTxUC := new(MockTransactionService)
		handler := NewHandler(mockAccUC, mockTxUC)

		r := gin.Default()
		r.POST("/accounts", handler.CreateAccount)

		expectedAcc := &domain.Account{ID: 1, DocumentNumber: "123"}
		mockAccUC.On("CreateAccount", mock.Anything, "123").Return(expectedAcc, nil)

		body := []byte(`{"document_number": "123"}`)
		req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockAccUC.AssertExpectations(t)
	})

	t.Run("BadRequest", func(t *testing.T) {
		mockAccUC := new(MockAccountService)
		mockTxUC := new(MockTransactionService)
		handler := NewHandler(mockAccUC, mockTxUC)

		r := gin.Default()
		r.POST("/accounts", handler.CreateAccount)

		body := []byte(`{"invalid": "json"}`) // Missing required field
		req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("InternalError", func(t *testing.T) {
		mockAccUC := new(MockAccountService)
		mockTxUC := new(MockTransactionService)
		handler := NewHandler(mockAccUC, mockTxUC)

		r := gin.Default()
		r.POST("/accounts", handler.CreateAccount)

		mockAccUC.On("CreateAccount", mock.Anything, "123").Return(nil, errors.New("db error"))

		body := []byte(`{"document_number": "123"}`)
		req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockAccUC.AssertExpectations(t)
	})
}

func TestHandler_GetAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockAccUC := new(MockAccountService)
		mockTxUC := new(MockTransactionService)
		handler := NewHandler(mockAccUC, mockTxUC)

		r := gin.Default()
		r.GET("/accounts/:accountId", handler.GetAccount)

		expectedAcc := &domain.Account{ID: 1, DocumentNumber: "123"}
		mockAccUC.On("GetAccount", mock.Anything, int64(1)).Return(expectedAcc, nil)

		req, _ := http.NewRequest("GET", "/accounts/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockAccUC.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockAccUC := new(MockAccountService)
		mockTxUC := new(MockTransactionService)
		handler := NewHandler(mockAccUC, mockTxUC)

		r := gin.Default()
		r.GET("/accounts/:accountId", handler.GetAccount)

		mockAccUC.On("GetAccount", mock.Anything, int64(1)).Return(nil, customErrors.ErrAccountNotFound)

		req, _ := http.NewRequest("GET", "/accounts/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockAccUC.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockAccUC := new(MockAccountService)
		mockTxUC := new(MockTransactionService)
		handler := NewHandler(mockAccUC, mockTxUC)

		r := gin.Default()
		r.GET("/accounts/:accountId", handler.GetAccount)

		req, _ := http.NewRequest("GET", "/accounts/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestHandler_CreateTransaction(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockAccUC := new(MockAccountService)
		mockTxUC := new(MockTransactionService)
		handler := NewHandler(mockAccUC, mockTxUC)

		r := gin.Default()
		r.POST("/transactions", handler.CreateTransaction)

		expectedTx := &domain.Transaction{ID: 1, Amount: -100.0}
		mockTxUC.On("CreateTransaction", mock.Anything, int64(1), 1, 100.0).Return(expectedTx, nil)

		body := []byte(`{"account_id": 1, "operation_type_id": 1, "amount": 100.0}`)
		req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockTxUC.AssertExpectations(t)
	})

	t.Run("BadRequest", func(t *testing.T) {
		mockAccUC := new(MockAccountService)
		mockTxUC := new(MockTransactionService)
		handler := NewHandler(mockAccUC, mockTxUC)

		r := gin.Default()
		r.POST("/transactions", handler.CreateTransaction)

		mockTxUC.On("CreateTransaction", mock.Anything, int64(1), 99, 100.0).Return(nil, customErrors.ErrOperationType)

		body := []byte(`{"account_id": 1, "operation_type_id": 99, "amount": 100.0}`)
		req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockTxUC.AssertExpectations(t)
	})
}
