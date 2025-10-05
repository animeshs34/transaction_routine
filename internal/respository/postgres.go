package respository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/animeshs34/transaction_routine/internal/domain"
	"github.com/animeshs34/transaction_routine/internal/logger"
	"go.uber.org/zap"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(conn *DBConn) *PostgresStore {
	return &PostgresStore{db: conn.GetDB()}
}

func (r *PostgresStore) CreateAccount(document string) (domain.Account, error) {
	var acc domain.Account
	err := r.db.QueryRow("INSERT INTO accounts (document_number) VALUES ($1) RETURNING id, document_number", document).Scan(&acc.ID, &acc.DocumentNumber)
	if err != nil {
		return domain.Account{}, fmt.Errorf("failed to create account: %w", err)
	}
	return acc, nil
}

func (r *PostgresStore) GetAccount(id int64) (domain.Account, error) {
	var acc domain.Account
	err := r.db.QueryRow("SELECT id, document_number FROM accounts WHERE id = $1", id).Scan(&acc.ID, &acc.DocumentNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Account{}, ErrAccountNotFound
		}
		return domain.Account{}, fmt.Errorf("failed to get account: %w", err)
	}
	return acc, nil
}

func (r *PostgresStore) HasOperationType(id int) bool {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM operation_types WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		logger.Error("Failed to check operation type", zap.Error(err), zap.Int("operation_type_id", id))
		return false
	}
	return exists
}

func (r *PostgresStore) CreateTransaction(t domain.Transaction) (domain.Transaction, error) {

	var accountExists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM accounts WHERE id = $1)", t.AccountID).Scan(&accountExists)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("failed to check account: %w", err)
	}
	if !accountExists {
		return domain.Transaction{}, ErrAccountNotFound
	}

	var operationTypeExists bool
	err = r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM operation_types WHERE id = $1)", t.OperationTypeID).Scan(&operationTypeExists)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("failed to check operation type: %w", err)
	}
	if !operationTypeExists {
		return domain.Transaction{}, ErrOperationTypeNotFound
	}

	if t.EventDate.IsZero() {
		t.EventDate = time.Now().UTC()
	}

	err = r.db.QueryRow(`
		INSERT INTO transactions (account_id, operation_type_id, amount, event_date)
		VALUES ($1, $2, $3, $4)
		RETURNING id, account_id, operation_type_id, amount, event_date
	`, t.AccountID, t.OperationTypeID, t.Amount, t.EventDate).Scan(
		&t.ID, &t.AccountID, &t.OperationTypeID, &t.Amount, &t.EventDate)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	return t, nil
}
