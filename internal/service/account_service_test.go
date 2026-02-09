package service

import (
	"context"
	"testing"

	"github.com/animesh/transaction_routine/internal/domain"
	customErrors "github.com/animesh/transaction_routine/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAccountRepository
type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) Create(ctx context.Context, account *domain.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockAccountRepository) FindByID(ctx context.Context, id int64) (*domain.Account, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Account), args.Error(1)
}

func TestAccountService_CreateAccount(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockAccountRepository)
		serviceInstance := NewAccountService(mockRepo)

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Account")).Return(nil)

		acc, err := serviceInstance.CreateAccount(context.Background(), "12345678900")

		assert.NoError(t, err)
		assert.NotNil(t, acc)
		assert.Equal(t, "12345678900", acc.DocumentNumber)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := new(MockAccountRepository)
		serviceInstance := NewAccountService(mockRepo)

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Account")).Return(customErrors.ErrInternal)

		acc, err := serviceInstance.CreateAccount(context.Background(), "12345678900")

		assert.Error(t, err)
		assert.Nil(t, acc)
		mockRepo.AssertExpectations(t)
	})
}

func TestAccountService_GetAccount(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockAccountRepository)
		serviceInstance := NewAccountService(mockRepo)

		expectedAcc := &domain.Account{ID: 1, DocumentNumber: "12345678900"}
		mockRepo.On("FindByID", mock.Anything, int64(1)).Return(expectedAcc, nil)

		acc, err := serviceInstance.GetAccount(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedAcc, acc)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo := new(MockAccountRepository)
		serviceInstance := NewAccountService(mockRepo)

		mockRepo.On("FindByID", mock.Anything, int64(1)).Return(nil, customErrors.ErrAccountNotFound)

		acc, err := serviceInstance.GetAccount(context.Background(), 1)

		assert.ErrorIs(t, err, customErrors.ErrAccountNotFound)
		assert.Nil(t, acc)
		mockRepo.AssertExpectations(t)
	})
}
