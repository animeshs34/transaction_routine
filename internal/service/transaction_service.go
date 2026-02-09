package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/animesh/transaction_routine/internal/domain"
	customErrors "github.com/animesh/transaction_routine/pkg/errors"
)

type transactionService struct {
	txRepo      TransactionRepository
	accountRepo AccountRepository
}

func NewTransactionService(txRepo TransactionRepository, accountRepo AccountRepository) TransactionService {
	return &transactionService{
		txRepo:      txRepo,
		accountRepo: accountRepo,
	}
}

func (s *transactionService) CreateTransaction(ctx context.Context, accountID int64, operationTypeID int, amount float64) (*domain.Transaction, error) {
	// 1. Verify if account exists
	_, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		if errors.Is(err, customErrors.ErrAccountNotFound) {
			return nil, fmt.Errorf("transaction service: %w", customErrors.ErrAccountNotFound)
		}
		return nil, fmt.Errorf("transaction service: failed to verify account: %w", err)
	}

	// 2. Create domain object
	tx, err := domain.NewTransaction(accountID, domain.OperationType(operationTypeID), amount)
	if err != nil {
		return nil, fmt.Errorf("transaction service: %w", customErrors.ErrOperationType)
	}

	// 3. Persist
	if err := s.txRepo.Create(ctx, tx); err != nil {
		return nil, fmt.Errorf("transaction service: failed to create transaction: %w", err)
	}

	return tx, nil
}
