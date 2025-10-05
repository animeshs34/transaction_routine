package domain

import (
	"testing"
)

func TestIsDebitOperation(t *testing.T) {
	if !IsDebitOperation(OpCashPurchase) {
		t.Errorf("OpCashPurchase should be debit")
	}
	if !IsDebitOperation(OpInstallmentPurchase) {
		t.Errorf("OpInstallmentPurchase should be debit")
	}
	if !IsDebitOperation(OpWithdrawal) {
		t.Errorf("OpWithdrawal should be debit")
	}
	if IsDebitOperation(OpPayment) {
		t.Errorf("OpPayment should not be debit")
	}
	if IsDebitOperation(999) {
		t.Errorf("Unknown op should not be debit")
	}
}

func TestIsCreditOperation(t *testing.T) {
	if !IsCreditOperation(OpPayment) {
		t.Errorf("OpPayment should be credit")
	}
	if IsCreditOperation(OpCashPurchase) {
		t.Errorf("OpCashPurchase should not be credit")
	}
	if IsCreditOperation(OpInstallmentPurchase) {
		t.Errorf("OpInstallmentPurchase should not be credit")
	}
	if IsCreditOperation(OpWithdrawal) {
		t.Errorf("OpWithdrawal should not be credit")
	}
	if IsCreditOperation(999) {
		t.Errorf("Unknown op should not be credit")
	}
}
