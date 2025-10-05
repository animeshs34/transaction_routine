package respository

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/animeshs34/transaction_routine/internal/domain"
	"regexp"
	"testing"
	"time"
)

func TestPostgresStore(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	store := &PostgresStore{db: db}

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO accounts (document_number) VALUES ($1) RETURNING id, document_number")).
		WithArgs("doc1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "document_number"}).AddRow(1, "doc1"))
	acc, err := store.CreateAccount("doc1")
	if err != nil || acc.ID != 1 || acc.DocumentNumber != "doc1" {
		t.Errorf("CreateAccount failed: %v", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO accounts (document_number) VALUES ($1) RETURNING id, document_number")).
		WithArgs("fail").
		WillReturnError(errors.New("fail"))
	_, err = store.CreateAccount("fail")
	if err == nil {
		t.Errorf("expected error for CreateAccount fail")
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, document_number FROM accounts WHERE id = $1")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "document_number"}).AddRow(1, "doc1"))
	acc, err = store.GetAccount(1)
	if err != nil || acc.ID != 1 {
		t.Errorf("GetAccount failed: %v", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, document_number FROM accounts WHERE id = $1")).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)
	_, err = store.GetAccount(999)
	if !errors.Is(err, ErrAccountNotFound) {
		t.Errorf("expected ErrAccountNotFound, got %v", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, document_number FROM accounts WHERE id = $1")).
		WithArgs(2).
		WillReturnError(errors.New("fail"))
	_, err = store.GetAccount(2)
	if err == nil {
		t.Errorf("expected error for GetAccount fail")
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM operation_types WHERE id = $1)")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	if !store.HasOperationType(1) {
		t.Errorf("expected HasOperationType true")
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM operation_types WHERE id = $1)")).
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
	if store.HasOperationType(2) {
		t.Errorf("expected HasOperationType false")
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM operation_types WHERE id = $1)")).
		WithArgs(3).
		WillReturnError(errors.New("fail"))
	if store.HasOperationType(3) {
		t.Errorf("expected HasOperationType false on error")
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM accounts WHERE id = $1)")).WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM operation_types WHERE id = $1)")).WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO transactions (account_id, operation_type_id, amount, event_date)
		VALUES ($1, $2, $3, $4)
		RETURNING id, account_id, operation_type_id, amount, event_date`)).
		WithArgs(1, 1, 100.0, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "operation_type_id", "amount", "event_date"}).AddRow(1, 1, 1, 100.0, time.Now()))

	tx := domain.Transaction{AccountID: 1, OperationTypeID: 1, Amount: 100}

	txResult, err := store.CreateTransaction(tx)
	if err != nil || txResult.ID != 1 {
		t.Errorf("CreateTransaction failed: %v", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM accounts WHERE id = $1)")).WithArgs(999).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
	tx = domain.Transaction{AccountID: 999, OperationTypeID: 1, Amount: 100}
	_, err = store.CreateTransaction(tx)
	if !errors.Is(err, ErrAccountNotFound) {
		t.Errorf("expected ErrAccountNotFound, got %v", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM accounts WHERE id = $1)")).WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM operation_types WHERE id = $1)")).WithArgs(999).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
	tx = domain.Transaction{AccountID: 1, OperationTypeID: 999, Amount: 100}
	_, err = store.CreateTransaction(tx)

	if !errors.Is(err, ErrOperationTypeNotFound) {
		t.Errorf("expected ErrOperationTypeNotFound, got %v", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM accounts WHERE id = $1)")).WithArgs(1).
		WillReturnError(errors.New("fail"))
	tx = domain.Transaction{AccountID: 1, OperationTypeID: 1, Amount: 100}
	_, err = store.CreateTransaction(tx)
	if err == nil {
		t.Errorf("expected error for account check fail")
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM accounts WHERE id = $1)")).WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM operation_types WHERE id = $1)")).WithArgs(1).
		WillReturnError(errors.New("fail"))
	tx = domain.Transaction{AccountID: 1, OperationTypeID: 1, Amount: 100}
	_, err = store.CreateTransaction(tx)
	if err == nil {
		t.Errorf("expected error for operation type check fail")
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM accounts WHERE id = $1)")).WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM operation_types WHERE id = $1)")).WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO transactions (account_id, operation_type_id, amount, event_date)
		VALUES ($1, $2, $3, $4)
		RETURNING id, account_id, operation_type_id, amount, event_date`)).
		WithArgs(1, 1, 100.0, sqlmock.AnyArg()).
		WillReturnError(errors.New("fail"))
	tx = domain.Transaction{AccountID: 1, OperationTypeID: 1, Amount: 100}
	_, err = store.CreateTransaction(tx)
	if err == nil {
		t.Errorf("expected error for insert fail")
	}
}
