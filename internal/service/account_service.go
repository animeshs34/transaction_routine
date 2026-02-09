package service

import (
	"context"
	"fmt"

	"github.com/animesh/transaction_routine/internal/domain"
)

type accountService struct {
	repo AccountRepository
}

func NewAccountService(repo AccountRepository) AccountService {
	return &accountService{repo: repo}
}

func (s *accountService) CreateAccount(ctx context.Context, documentNumber string) (*domain.Account, error) {
	account := &domain.Account{
		DocumentNumber: documentNumber,
	}

	if err := s.repo.Create(ctx, account); err != nil {
		return nil, fmt.Errorf("account service: failed to create account: %w", err)
	}

	return account, nil
}

func (s *accountService) GetAccount(ctx context.Context, id int64) (*domain.Account, error) {
	acc, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("account service: failed to find account: %w", err)
	}
	return acc, nil
}
