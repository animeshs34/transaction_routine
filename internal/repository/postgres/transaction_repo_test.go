package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/animesh/transaction_routine/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestTransactionRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewTransactionRepository(db)

	t.Run("Success", func(t *testing.T) {
		tx := &domain.Transaction{
			AccountID:       1,
			OperationTypeID: domain.NormalPurchase,
			Amount:          -100.0,
			EventDate:       time.Now(),
		}

		mock.ExpectQuery(`INSERT INTO transactions`).
			WithArgs(tx.AccountID, tx.OperationTypeID, tx.Amount, sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"transaction_id"}).AddRow(1))

		err := repo.Create(context.Background(), tx)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), tx.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DBError", func(t *testing.T) {
		tx := &domain.Transaction{
			AccountID:       1,
			OperationTypeID: domain.NormalPurchase,
			Amount:          -100.0,
		}

		mock.ExpectQuery(`INSERT INTO transactions`).
			WithArgs(tx.AccountID, tx.OperationTypeID, tx.Amount, sqlmock.AnyArg()).
			WillReturnError(errors.New("db error"))

		err := repo.Create(context.Background(), tx)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
