package service

import (
	"testing"
	"time"

	"github.com/animeshs34/transaction_routine/internal/domain"
	"github.com/animeshs34/transaction_routine/internal/respository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) CreateAccount(document string) (domain.Account, error) {
	args := m.Called(document)
	return args.Get(0).(domain.Account), args.Error(1)
}
func (m *mockRepo) GetAccount(id int64) (domain.Account, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Account), args.Error(1)
}
func (m *mockRepo) HasOperationType(id int) bool {
	args := m.Called(id)
	return args.Bool(0)
}
func (m *mockRepo) CreateTransaction(tx domain.Transaction) (domain.Transaction, error) {
	args := m.Called(tx)
	return args.Get(0).(domain.Transaction), args.Error(1)
}

func TestCreateAccount_Valid(t *testing.T) {
	repo := new(mockRepo)
	svc := New(repo)
	doc := "12345"
	acc := domain.Account{ID: 1, DocumentNumber: doc}
	repo.On("CreateAccount", doc).Return(acc, nil)
	result, err := svc.CreateAccount(doc)
	assert.NoError(t, err)
	assert.Equal(t, acc, result)
}

func TestCreateAccount_Invalid(t *testing.T) {
	svc := New(new(mockRepo))
	_, err := svc.CreateAccount("")
	assert.ErrorIs(t, err, ErrInvalidDocument)
}

func TestCreateAccount_Whitespace(t *testing.T) {
	svc := New(new(mockRepo))
	_, err := svc.CreateAccount("   ")
	assert.ErrorIs(t, err, ErrInvalidDocument)
}

func TestGetAccount(t *testing.T) {
	repo := new(mockRepo)
	svc := New(repo)
	acc := domain.Account{ID: 1, DocumentNumber: "doc"}
	repo.On("GetAccount", int64(1)).Return(acc, nil)
	result, err := svc.GetAccount(1)
	assert.NoError(t, err)
	assert.Equal(t, acc, result)
}

func TestCreateTransaction_InvalidOperationType(t *testing.T) {
	repo := new(mockRepo)
	svc := New(repo)
	repo.On("HasOperationType", 99).Return(false)
	_, err := svc.CreateTransaction(1, 99, 100, nil)
	assert.ErrorIs(t, err, ErrInvalidOperationType)
}

func TestCreateTransaction_InvalidAmount(t *testing.T) {
	repo := new(mockRepo)
	svc := New(repo)
	repo.On("HasOperationType", 1).Return(true)
	_, err := svc.CreateTransaction(1, 1, 0, nil)
	assert.ErrorIs(t, err, ErrInvalidAmount)
}

func TestCreateTransaction_Success(t *testing.T) {
	repo := new(mockRepo)
	svc := New(repo)
	repo.On("HasOperationType", 1).Return(true)
	timeNow := time.Now()
	tx := domain.Transaction{AccountID: 1, OperationTypeID: 1, Amount: -100, EventDate: timeNow.UTC()}
	repo.On("CreateTransaction", mock.AnythingOfType("domain.Transaction")).Return(tx, nil)
	result, err := svc.CreateTransaction(1, 1, 100, &timeNow)
	assert.NoError(t, err)
	assert.Equal(t, tx.AccountID, result.AccountID)
	assert.Equal(t, tx.Amount, result.Amount)
}

func TestCreateTransaction_UnknownError(t *testing.T) {
	repo := new(mockRepo)
	svc := New(repo)
	repo.On("HasOperationType", 1).Return(true)
	timeNow := time.Now()
	repo.On("CreateTransaction", mock.AnythingOfType("domain.Transaction")).Return(domain.Transaction{}, assert.AnError)
	result, err := svc.CreateTransaction(1, 1, 100, &timeNow)
	assert.ErrorIs(t, err, assert.AnError)
	assert.Equal(t, domain.Transaction{}, result)
}

func TestCreateTransaction_AccountNotFound(t *testing.T) {
	repo := new(mockRepo)
	svc := New(repo)
	repo.On("HasOperationType", 1).Return(true)
	timeNow := time.Now()
	repo.On("CreateTransaction", mock.AnythingOfType("domain.Transaction")).Return(domain.Transaction{}, respository.ErrAccountNotFound)
	result, err := svc.CreateTransaction(1, 1, 100, &timeNow)
	assert.ErrorIs(t, err, respository.ErrAccountNotFound)
	assert.Equal(t, domain.Transaction{}, result)
}

func TestCreateTransaction_OperationTypeNotFound(t *testing.T) {
	repo := new(mockRepo)
	svc := New(repo)
	repo.On("HasOperationType", 1).Return(true)
	timeNow := time.Now()
	repo.On("CreateTransaction", mock.AnythingOfType("domain.Transaction")).Return(domain.Transaction{}, respository.ErrOperationTypeNotFound)
	result, err := svc.CreateTransaction(1, 1, 100, &timeNow)
	assert.ErrorIs(t, err, ErrInvalidOperationType)
	assert.Equal(t, domain.Transaction{}, result)
}

func TestCreateTransaction_NeitherDebitNorCredit(t *testing.T) {
	repo := new(mockRepo)
	svc := New(repo)
	repo.On("HasOperationType", 5).Return(true)
	// Patch domain.IsDebitOperation and IsCreditOperation to false
	result, err := svc.CreateTransaction(1, 5, 100, nil)
	assert.ErrorIs(t, err, ErrInvalidOperationType)
	assert.Equal(t, domain.Transaction{}, result)
}

func TestCreateTransaction_CreditOperation(t *testing.T) {
	repo := new(mockRepo)
	svc := New(repo)
	repo.On("HasOperationType", domain.OpPayment).Return(true)
	timeNow := time.Now()
	tx := domain.Transaction{AccountID: 1, OperationTypeID: domain.OpPayment, Amount: 100, EventDate: timeNow.UTC()}
	repo.On("CreateTransaction", mock.AnythingOfType("domain.Transaction")).Return(tx, nil)
	result, err := svc.CreateTransaction(1, domain.OpPayment, 100, &timeNow)
	assert.NoError(t, err)
	assert.Equal(t, tx.AccountID, result.AccountID)
	assert.Equal(t, tx.Amount, result.Amount)
}
