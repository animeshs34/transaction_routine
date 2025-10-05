package respository

import (
	"sync"
	"time"

	"github.com/animeshs34/transaction_routine/internal/domain"
)

type InMemoryStore struct {
	mu sync.RWMutex

	accounts       map[int64]*domain.Account
	transactions   map[int64]*domain.Transaction
	operationTypes map[int]domain.OperationType

	nextAccountID     int64
	nextTransactionID int64
}

func NewInMemoryStore() *InMemoryStore {
	r := &InMemoryStore{
		accounts:          make(map[int64]*domain.Account),
		transactions:      make(map[int64]*domain.Transaction),
		operationTypes:    make(map[int]domain.OperationType),
		nextAccountID:     1,
		nextTransactionID: 1,
	}

	r.operationTypes[domain.OpCashPurchase] = domain.OperationType{ID: domain.OpCashPurchase, Description: "CASH PURCHASE"}
	r.operationTypes[domain.OpInstallmentPurchase] = domain.OperationType{ID: domain.OpInstallmentPurchase, Description: "INSTALLMENT PURCHASE"}
	r.operationTypes[domain.OpWithdrawal] = domain.OperationType{ID: domain.OpWithdrawal, Description: "WITHDRAWAL"}
	r.operationTypes[domain.OpPayment] = domain.OperationType{ID: domain.OpPayment, Description: "PAYMENT"}

	return r
}

func (r *InMemoryStore) CreateAccount(document string) (domain.Account, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	acc := domain.Account{
		ID:             r.nextAccountID,
		DocumentNumber: document,
	}
	r.accounts[acc.ID] = &acc
	r.nextAccountID++
	return acc, nil
}

func (r *InMemoryStore) GetAccount(id int64) (domain.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	a, ok := r.accounts[id]
	if !ok {
		return domain.Account{}, ErrAccountNotFound
	}
	return *a, nil
}

func (r *InMemoryStore) HasOperationType(id int) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.operationTypes[id]
	return ok
}

func (r *InMemoryStore) CreateTransaction(t domain.Transaction) (domain.Transaction, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.accounts[t.AccountID]; !ok {
		return domain.Transaction{}, ErrAccountNotFound
	}
	if _, ok := r.operationTypes[t.OperationTypeID]; !ok {
		return domain.Transaction{}, ErrOperationTypeNotFound
	}

	t.ID = r.nextTransactionID
	if t.EventDate.IsZero() {
		t.EventDate = time.Now().UTC()
	}

	r.transactions[t.ID] = &t
	r.nextTransactionID++
	return t, nil
}
