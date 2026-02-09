package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/animesh/transaction_routine/internal/domain"
	customErrors "github.com/animesh/transaction_routine/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestAccountRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewAccountRepository(db)

	t.Run("Success", func(t *testing.T) {
		account := &domain.Account{DocumentNumber: "12345678900"}
		mock.ExpectQuery("INSERT INTO accounts").
			WithArgs(account.DocumentNumber).
			WillReturnRows(sqlmock.NewRows([]string{"account_id"}).AddRow(1))

		err := repo.Create(context.Background(), account)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), account.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DBError", func(t *testing.T) {
		account := &domain.Account{DocumentNumber: "12345678900"}
		mock.ExpectQuery("INSERT INTO accounts").
			WithArgs(account.DocumentNumber).
			WillReturnError(errors.New("db error"))

		err := repo.Create(context.Background(), account)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAccountRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewAccountRepository(db)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery("SELECT account_id, document_number FROM accounts").
			WithArgs(int64(1)).
			WillReturnRows(sqlmock.NewRows([]string{"account_id", "document_number"}).AddRow(1, "12345678900"))

		acc, err := repo.FindByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.NotNil(t, acc)
		assert.Equal(t, int64(1), acc.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery("SELECT account_id, document_number FROM accounts").
			WithArgs(int64(1)).
			WillReturnRows(sqlmock.NewRows([]string{"account_id", "document_number"}))

		acc, err := repo.FindByID(context.Background(), 1)

		assert.ErrorIs(t, err, customErrors.ErrAccountNotFound)
		assert.Nil(t, acc)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectQuery("SELECT account_id, document_number FROM accounts").
			WithArgs(int64(1)).
			WillReturnError(errors.New("db error"))

		acc, err := repo.FindByID(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, acc)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
