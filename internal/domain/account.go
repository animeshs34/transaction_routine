package domain

type Account struct {
	ID             int64  `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}

func NewAccount(id int64, documentNumber string) *Account {
	return &Account{
		ID:             id,
		DocumentNumber: documentNumber,
	}
}
