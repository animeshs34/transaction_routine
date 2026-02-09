package postgres

import (
	"context"
	"database/sql"

	"github.com/animesh/transaction_routine/internal/domain"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	query := `INSERT INTO transactions (account_id, operation_type_id, amount, event_date) 
	          VALUES ($1, $2, $3, $4) RETURNING transaction_id`

	err := r.db.QueryRowContext(ctx, query,
		tx.AccountID,
		tx.OperationTypeID,
		tx.Amount,
		tx.EventDate,
	).Scan(&tx.ID)

	if err != nil {
		return err
	}
	return nil
}
