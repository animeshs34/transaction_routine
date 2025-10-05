package respository

import (
	"github.com/animeshs34/transaction_routine/internal/domain"
	"testing"
	"time"
)

func TestMemoryStore(t *testing.T) {
	r := NewInMemoryStore()

	acc, err := r.CreateAccount("doc1")
	if err != nil {
		t.Fatalf("CreateAccount failed: %v", err)
	}
	if acc.ID == 0 || acc.DocumentNumber != "doc1" {
		t.Errorf("unexpected account: %+v", acc)
	}

	got, err := r.GetAccount(acc.ID)
	if err != nil || got.ID != acc.ID {
		t.Errorf("GetAccount failed: %v", err)
	}

	_, err = r.GetAccount(9999)
	if err == nil {
		t.Errorf("expected error for missing account")
	}

	if !r.HasOperationType(domain.OpCashPurchase) {
		t.Errorf("expected true for OpCashPurchase")
	}
	if r.HasOperationType(999) {
		t.Errorf("expected false for unknown op type")
	}

	tx := domain.Transaction{AccountID: acc.ID, OperationTypeID: domain.OpCashPurchase, Amount: 100}
	txResult, err := r.CreateTransaction(tx)
	if err != nil {
		t.Errorf("CreateTransaction failed: %v", err)
	}
	if txResult.ID == 0 || txResult.AccountID != acc.ID {
		t.Errorf("unexpected transaction: %+v", txResult)
	}

	tx = domain.Transaction{AccountID: 9999, OperationTypeID: domain.OpCashPurchase, Amount: 100}
	_, err = r.CreateTransaction(tx)
	if err == nil {
		t.Errorf("expected error for missing account")
	}

	tx = domain.Transaction{AccountID: acc.ID, OperationTypeID: 999, Amount: 100}
	_, err = r.CreateTransaction(tx)
	if err == nil {
		t.Errorf("expected error for missing operation type")
	}

	tx = domain.Transaction{AccountID: acc.ID, OperationTypeID: domain.OpCashPurchase, Amount: 100}
	tx.EventDate = time.Time{} // zero
	txResult, err = r.CreateTransaction(tx)
	if err != nil {
		t.Errorf("CreateTransaction failed: %v", err)
	}
	if txResult.EventDate.IsZero() {
		t.Errorf("expected EventDate to be set")
	}
}
