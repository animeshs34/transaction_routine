package postgres

import (
	"context"
	"database/sql"

	"github.com/animesh/transaction_routine/internal/domain"
	customErrors "github.com/animesh/transaction_routine/pkg/errors"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(ctx context.Context, account *domain.Account) error {
	query := `INSERT INTO accounts (document_number) VALUES ($1) RETURNING account_id`
	err := r.db.QueryRowContext(ctx, query, account.DocumentNumber).Scan(&account.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *AccountRepository) FindByID(ctx context.Context, id int64) (*domain.Account, error) {
	query := `SELECT account_id, document_number FROM accounts WHERE account_id = $1`
	account := &domain.Account{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&account.ID, &account.DocumentNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customErrors.ErrAccountNotFound
		}
		return nil, err
	}
	return account, nil
}
