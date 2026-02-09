package domain

import (
	"fmt"
	"time"
)

type OperationType int

const (
	NormalPurchase           OperationType = 1
	PurchaseWithInstallments OperationType = 2
	Withdrawal               OperationType = 3
	CreditVoucher            OperationType = 4
)

func (ot OperationType) IsValid() bool {
	switch ot {
	case NormalPurchase, PurchaseWithInstallments, Withdrawal, CreditVoucher:
		return true
	}
	return false
}

func (ot OperationType) String() string {
	switch ot {
	case NormalPurchase:
		return "Normal Purchase"
	case PurchaseWithInstallments:
		return "Purchase with Installments"
	case Withdrawal:
		return "Withdrawal"
	case CreditVoucher:
		return "Credit Voucher"
	default:
		return "Unknown"
	}
}

// Transaction represents a financial record associated with an account.
type Transaction struct {
	ID              int64         `json:"transaction_id"`
	AccountID       int64         `json:"account_id"`
	OperationTypeID OperationType `json:"operation_type_id"`
	Amount          float64       `json:"amount"`
	EventDate       time.Time     `json:"event_date"`
}

// NewTransaction creates and normalizes a new transaction record.
func NewTransaction(accountID int64, operationTypeID OperationType, amount float64) (*Transaction, error) {
	if !operationTypeID.IsValid() {
		return nil, fmt.Errorf("invalid operation type: %d", operationTypeID)
	}

	tx := &Transaction{
		AccountID:       accountID,
		OperationTypeID: operationTypeID,
		Amount:          amount,
		EventDate:       time.Now(),
	}

	tx.CalculateFinalAmount()
	return tx, nil
}

// CalculateFinalAmount applies the business rule:
// NormalPurchase, PurchaseWithInstallments, and Withdrawal MUST be negative.
// CreditVoucher MUST be positive.
func (t *Transaction) CalculateFinalAmount() {
	switch t.OperationTypeID {
	case NormalPurchase, PurchaseWithInstallments, Withdrawal:
		if t.Amount > 0 {
			t.Amount = -t.Amount
		}
	case CreditVoucher:
		if t.Amount < 0 {
			t.Amount = -t.Amount
		}
	}
}
