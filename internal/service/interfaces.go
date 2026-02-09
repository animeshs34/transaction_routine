package service

import (
	"context"

	"github.com/animesh/transaction_routine/internal/domain"
)

// AccountRepository defines the persistence operations for accounts.
type AccountRepository interface {
	Create(ctx context.Context, account *domain.Account) error
	FindByID(ctx context.Context, id int64) (*domain.Account, error)
}

// TransactionRepository defines the persistence operations for transactions.
type TransactionRepository interface {
	Create(ctx context.Context, transaction *domain.Transaction) error
}

// AccountService defines the business logic for managing accounts.
type AccountService interface {
	CreateAccount(ctx context.Context, documentNumber string) (*domain.Account, error)
	GetAccount(ctx context.Context, id int64) (*domain.Account, error)
}

// TransactionService defines the business logic for processing transactions.
type TransactionService interface {
	CreateTransaction(ctx context.Context, accountID int64, operationTypeID int, amount float64) (*domain.Transaction, error)
}
