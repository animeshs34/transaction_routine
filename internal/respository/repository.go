package respository

import (
	"errors"

	"github.com/animeshs34/transaction_routine/internal/domain"
)

var (
	ErrAccountNotFound       = errors.New("account not found")
	ErrOperationTypeNotFound = errors.New("operation type not found")
)

type Respository interface {
	CreateAccount(document string) (domain.Account, error)
	GetAccount(id int64) (domain.Account, error)
	HasOperationType(id int) bool
	CreateTransaction(t domain.Transaction) (domain.Transaction, error)
}
