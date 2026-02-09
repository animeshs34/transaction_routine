CREATE TABLE IF NOT EXISTS accounts (
    account_id SERIAL PRIMARY KEY,
    document_number VARCHAR(20) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS operations_types (
    operation_type_id SERIAL PRIMARY KEY,
    description VARCHAR(50) NOT NULL
);

INSERT INTO operations_types (operation_type_id, description) VALUES
(1, 'Normal Purchase'),
(2, 'Purchase with Installments'),
(3, 'Withdrawal'),
(4, 'Credit Voucher')
ON CONFLICT (operation_type_id) DO NOTHING;

CREATE TABLE IF NOT EXISTS transactions (
    transaction_id SERIAL PRIMARY KEY,
    account_id INT NOT NULL,
    operation_type_id INT NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    event_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_account FOREIGN KEY(account_id) REFERENCES accounts(account_id),
    CONSTRAINT fk_operation_type FOREIGN KEY(operation_type_id) REFERENCES operations_types(operation_type_id)
);
