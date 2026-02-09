package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransaction_CalculateFinalAmount(t *testing.T) {
	tests := []struct {
		name            string
		operationTypeID OperationType
		inputAmount     float64
		expectedAmount  float64
	}{
		{
			name:            "Normal Purchase - Positive Input",
			operationTypeID: NormalPurchase,
			inputAmount:     100.0,
			expectedAmount:  -100.0,
		},
		{
			name:            "Normal Purchase - Negative Input",
			operationTypeID: NormalPurchase,
			inputAmount:     -100.0,
			expectedAmount:  -100.0,
		},
		{
			name:            "Withdrawal - Positive Input",
			operationTypeID: Withdrawal,
			inputAmount:     50.0,
			expectedAmount:  -50.0,
		},
		{
			name:            "Credit Voucher - Positive Input",
			operationTypeID: CreditVoucher,
			inputAmount:     200.0,
			expectedAmount:  200.0,
		},
		{
			name:            "Credit Voucher - Negative Input",
			operationTypeID: CreditVoucher,
			inputAmount:     -200.0,
			expectedAmount:  200.0, // Should probably be forced positive
		},
		{
			name:            "Purchase with Installments - Positive Input",
			operationTypeID: PurchaseWithInstallments,
			inputAmount:     300.0,
			expectedAmount:  -300.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &Transaction{
				OperationTypeID: tt.operationTypeID,
				Amount:          tt.inputAmount,
			}
			tx.CalculateFinalAmount()
			assert.Equal(t, tt.expectedAmount, tx.Amount)
		})
	}
}

func TestNewTransaction_InvalidType(t *testing.T) {
	tx, err := NewTransaction(1, OperationType(99), 100.0)
	assert.Error(t, err)
	assert.Nil(t, tx)
}
