package respository

import (
	"database/sql"
	"fmt"

	"github.com/animeshs34/transaction_routine/internal/domain"
	_ "github.com/lib/pq"
)

type DBConn struct {
	db *sql.DB
}

func NewPostgresConn(host string, port int, user, password, dbname, sslmode string) (*DBConn, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := initSchema(db); err != nil {
		return nil, fmt.Errorf("failed to initialize database schema: %w", err)
	}

	if err := seedOperationTypes(db); err != nil {
		return nil, fmt.Errorf("failed to seed operation types: %w", err)
	}

	return &DBConn{db: db}, nil
}

func (c *DBConn) GetDB() *sql.DB {
	return c.db
}

func (c *DBConn) Close() error {
	return c.db.Close()
}

func initSchema(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS accounts (
			id SERIAL PRIMARY KEY,
			document_number TEXT NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create accounts table: %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS operation_types (
			id INT PRIMARY KEY,
			description TEXT NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create operation_types table: %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS transactions (
			id SERIAL PRIMARY KEY,
			account_id INT NOT NULL REFERENCES accounts(id),
			operation_type_id INT NOT NULL REFERENCES operation_types(id),
			amount DECIMAL(15,2) NOT NULL,
			event_date TIMESTAMP WITH TIME ZONE NOT NULL,
			CONSTRAINT fk_account FOREIGN KEY (account_id) REFERENCES accounts(id),
			CONSTRAINT fk_operation_type FOREIGN KEY (operation_type_id) REFERENCES operation_types(id)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create transactions table: %w", err)
	}

	return nil
}

func seedOperationTypes(db *sql.DB) error {

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM operation_types").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check operation types: %w", err)
	}

	if count > 0 {
		return nil
	}

	operationTypes := []struct {
		ID          int
		Description string
	}{
		{domain.OpCashPurchase, "CASH PURCHASE"},
		{domain.OpInstallmentPurchase, "INSTALLMENT PURCHASE"},
		{domain.OpWithdrawal, "WITHDRAWAL"},
		{domain.OpPayment, "PAYMENT"},
	}

	for _, ot := range operationTypes {
		_, err := db.Exec("INSERT INTO operation_types (id, description) VALUES ($1, $2)", ot.ID, ot.Description)
		if err != nil {
			return fmt.Errorf("failed to insert operation type %d: %w", ot.ID, err)
		}
	}

	return nil
}
