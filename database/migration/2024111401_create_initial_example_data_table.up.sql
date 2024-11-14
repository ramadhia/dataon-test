CREATE TABLE users (
    id VARCHAR(50) PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    address TEXT,
    pin varchar(50) NOT NULL,
    created_date TIMESTAMPTZ NOT NULL,
    updated_date TIMESTAMPTZ
);

CREATE TABLE access_tokens (
    token_id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    access_token TEXT NOT NULL,
    created_date TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE transactions (
    id VARCHAR(50) PRIMARY KEY,
    status VARCHAR(20) NOT NULL,
    user_id VARCHAR(50) NOT NULL,
    transaction_method VARCHAR(50) NOT NULL,
    transaction_type VARCHAR(10) NOT NULL,
    amount NUMERIC(15, 2) NOT NULL,
    remarks TEXT,
    balance_before NUMERIC(15, 2) NOT NULL,
    balance_after NUMERIC(15, 2) NOT NULL,
    created_date TIMESTAMPTZ NOT NULL,
    updated_date TIMESTAMPTZ
);

CREATE TABLE balances (
    user_id VARCHAR(50) NOT NULL,
    amount NUMERIC(15, 2) NOT NULL,
    last_updated TIMESTAMPTZ NULL
);

-- ------------------
-- INITIALIZATION --
-- ------------------

-- Seeder untuk tabel transactions
-- INSERT INTO transactions (id, status, user_id, transaction_method, transaction_type, amount, remarks, balance_before, balance_after, created_date)
-- VALUES
--     ('a7d39cf6-44b6-41fc-b3e9-7b16df5321c5', 'SUCCESS', 'bc1c823e-b0fb-4b20-88c0-dff25e283252', 'BANK_TRANSFER', 'DEBIT', 30000.00, 'Hadiah Ultah', 400000.00, 370000.00, '2021-04-01 22:23:20'),
--     ('13bcb11c-111e-4a65-9afd-90a86a01cd21', 'SUCCESS', 'bc1c823e-b0fb-4b20-88c0-dff25e283252', 'ONLINE_PAYMENT', 'DEBIT', 10000.00, 'Pulsa Telkomsel 100k', 500000.00, 400000.00, '2021-04-01 22:22:00'),
--     ('201ddde1-f797-484b-b1a0-07d1190e790a', 'SUCCESS', 'bc1c823e-b0fb-4b20-88c0-dff25e283252', 'TOP_UP', 'CREDIT', 500000.00, '', 0.00, 500000.00, '2021-04-01 22:21:21'),
--     ('12f34634-00c2-45c9-b3fa-627f1b8634c6', 'FAILED', 'bc1c823e-b0fb-4b20-88c0-dff25e283252', 'BANK_TRANSFER', 'DEBIT', 15000.00, 'Pembayaran Tagihan', 200000.00, 185000.00, '2021-04-02 10:12:00'),
--     ('bb6a88a1-c5b0-49b3-84b5-54ed4dbd9b84', 'PENDING', 'bc1c823e-b0fb-4b20-88c0-dff25e283252', 'PAYPAL', 'CREDIT', 120000.00, 'Refund', 300000.00, 420000.00, '2021-04-03 09:45:00');
