package domain

import "time"

type Account struct {
	ID             int64  `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}

type Transaction struct {
	ID              int64     `json:"transaction_id"`
	AccountID       int64     `json:"account_id"`
	OperationTypeID int       `json:"operation_type_id"`
	Amount          float64   `json:"amount"`
	EventDate       time.Time `json:"event_date"`
}

type OperationType struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

const (
	OpCashPurchase        = 1
	OpInstallmentPurchase = 2
	OpWithdrawal          = 3
	OpPayment             = 4
)

func IsDebitOperation(opID int) bool {
	return opID == OpCashPurchase || opID == OpInstallmentPurchase || opID == OpWithdrawal
}

func IsCreditOperation(opID int) bool {
	return opID == OpPayment
}
