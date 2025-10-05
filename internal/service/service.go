package service

import (
	"errors"
	"math"
	"strings"
	"time"

	"github.com/animeshs34/transaction_routine/internal/domain"
	"github.com/animeshs34/transaction_routine/internal/respository"
)

var (
	ErrInvalidDocument      = errors.New("invalid document_number")
	ErrInvalidOperationType = errors.New("invalid operation_type_id")
	ErrInvalidAmount        = errors.New("amount must be greater than zero")
)

type Repository = respository.Respository

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateAccount(document string) (domain.Account, error) {
	document = strings.TrimSpace(document)
	if document == "" {
		return domain.Account{}, ErrInvalidDocument
	}
	return s.repo.CreateAccount(document)
}

func (s *Service) GetAccount(id int64) (domain.Account, error) {
	return s.repo.GetAccount(id)
}

func (s *Service) CreateTransaction(accountID int64, operationTypeID int, amount float64, eventTime *time.Time) (domain.Transaction, error) {
	if !s.repo.HasOperationType(operationTypeID) {
		return domain.Transaction{}, ErrInvalidOperationType
	}
	a := math.Abs(amount)
	if a <= 0 {
		return domain.Transaction{}, ErrInvalidAmount
	}

	if domain.IsDebitOperation(operationTypeID) {
		a = -a
	} else if domain.IsCreditOperation(operationTypeID) {
		// positive, already correct value is here.
	} else {
		return domain.Transaction{}, ErrInvalidOperationType
	}

	var ts time.Time
	if eventTime != nil && !eventTime.IsZero() {
		ts = eventTime.UTC()
	}

	tx := domain.Transaction{
		AccountID:       accountID,
		OperationTypeID: operationTypeID,
		Amount:          a,
		EventDate:       ts,
	}
	created, err := s.repo.CreateTransaction(tx)
	if err != nil {
		if errors.Is(err, respository.ErrAccountNotFound) {
			return domain.Transaction{}, respository.ErrAccountNotFound
		}
		if errors.Is(err, respository.ErrOperationTypeNotFound) {
			return domain.Transaction{}, ErrInvalidOperationType
		}
		return domain.Transaction{}, err
	}
	return created, nil
}
