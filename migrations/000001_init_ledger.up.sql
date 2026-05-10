-- 1. Idempotency Table: The Gatekeeper
CREATE TABLE idempotency_keys (
    key VARCHAR(255) PRIMARY KEY,
    request_path TEXT NOT NULL,
    response_code INTEGER,
    response_body JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Ledger Table: The Source of Truth
-- Note: We store 'amount' as BIGINT (cents) to avoid floating point errors.
CREATE TABLE ledger_entries (
    id UUID PRIMARY KEY,
    account_id UUID NOT NULL,
    amount BIGINT NOT NULL, 
    currency VARCHAR(3) NOT NULL,
    type VARCHAR(10) CHECK (type IN ('debit', 'credit')),
    funding_source VARCHAR(50), -- e.g., 'QR_PAYNOW', 'CREDIT_CARD', 'CRYPTO'
    metadata JSONB,
    status VARCHAR(20) CHECK (status IN ('pending', 'success', 'failed', 'uncertain', 'reversed'))
    idempotency_key VARCHAR(255) UNIQUE,
    settled_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create an index for faster account balance lookups later
CREATE INDEX idx_ledger_account_id ON ledger_entries(account_id);