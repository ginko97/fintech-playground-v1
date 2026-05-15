package domain

import (
	"context"

	"github.com/google/uuid"
)

type TransactionRepository interface {
	// Create handles the idempotency logic and the ledger entry in one TX
	Create(ctx context.Context, tx *Transaction) error
	GetByInternalID(ctx context.Context, id uuid.UUID) (*Transaction, error)
	GetByIdempotencyKey(ctx context.Context, key string) (*Transaction, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status TransactionStatus) error
}
