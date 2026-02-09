package service

import (
	"context"
	"testing"

	"github.com/animesh/transaction_routine/internal/domain"
	customErrors "github.com/animesh/transaction_routine/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTransactionRepository
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(ctx context.Context, transaction *domain.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func TestTransactionService_CreateTransaction(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockAccRepo := new(MockAccountRepository)
		mockTxRepo := new(MockTransactionRepository)
		service := NewTransactionService(mockTxRepo, mockAccRepo)

		mockAccRepo.On("FindByID", mock.Anything, int64(1)).Return(&domain.Account{ID: 1}, nil)
		mockTxRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Transaction")).Return(nil)

		tx, err := service.CreateTransaction(context.Background(), 1, 1, 100.0)

		assert.NoError(t, err)
		assert.NotNil(t, tx)
		assert.Equal(t, -100.0, tx.Amount) // Check logic
		mockAccRepo.AssertExpectations(t)
		mockTxRepo.AssertExpectations(t)
	})

	t.Run("AccountNotFound", func(t *testing.T) {
		mockAccRepo := new(MockAccountRepository)
		mockTxRepo := new(MockTransactionRepository)
		service := NewTransactionService(mockTxRepo, mockAccRepo)

		mockAccRepo.On("FindByID", mock.Anything, int64(1)).Return(nil, customErrors.ErrAccountNotFound)

		tx, err := service.CreateTransaction(context.Background(), 1, 1, 100.0)

		assert.ErrorIs(t, err, customErrors.ErrAccountNotFound)
		assert.Nil(t, tx)
		mockAccRepo.AssertExpectations(t)
	})

	t.Run("InvalidOperationType", func(t *testing.T) {
		mockAccRepo := new(MockAccountRepository)
		mockTxRepo := new(MockTransactionRepository)
		service := NewTransactionService(mockTxRepo, mockAccRepo)

		mockAccRepo.On("FindByID", mock.Anything, int64(1)).Return(&domain.Account{ID: 1}, nil)

		tx, err := service.CreateTransaction(context.Background(), 1, 99, 100.0)

		assert.ErrorIs(t, err, customErrors.ErrOperationType)
		assert.Nil(t, tx)
		mockAccRepo.AssertExpectations(t)
	})
}
