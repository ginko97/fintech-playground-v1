package postgres

import (
	"context"

	"github.com/ginko97/fintech-playground-v1/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LedgerRepo struct {
	pool *pgxpool.Pool
}

func NewLedgerRepo(pool *pgxpool.Pool) *LedgerRepo {
	return &LedgerRepo{pool: pool}
}

func (r *LedgerRepo) Create(ctx context.Context, tx *domain.Transaction) error {
	dbTx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer dbTx.Rollback(ctx) // Safe: rollback does nothing if committed

	// 2. Insert Idempotency Key (The Gatekeeper)
	// If this fails with 23505, someone is retrying a request currently in progress.
	_, err = dbTx.Exec(ctx, `
        INSERT INTO idempotency_keys (key, request_path) 
        VALUES ($1, $2)`,
		tx.IdempotencyKey, "/v1/transactions")

	if err != nil {
		// Handle PostgreSQL Unique Violation
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return domain.ErrDuplicateRequest
		}
		return err
	}

	// 3. Insert Ledger Entry
	// Notice we use 'now()' for timestamps if not provided, but here we scan them back
	err = dbTx.QueryRow(ctx, `
        INSERT INTO ledger_entries (
            id, account_id, amount, currency, type, funding_source, status, idempotency_key
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING created_at, updated_at`,
		tx.ID, tx.AccountID, tx.Amount, tx.Currency, tx.Type, tx.FundingSource, tx.Status, tx.IdempotencyKey,
	).Scan(&tx.CreatedAt, &tx.UpdatedAt)

	if err != nil {
		return err
	}

	// 4. Commit
	return dbTx.Commit(ctx)
}

func (r *LedgerRepo) GetByInternalID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	return nil, nil // TODO: Implement
}

func (r *LedgerRepo) GetByIdempotencyKey(ctx context.Context, key string) (*domain.Transaction, error) {
	return nil, nil // TODO: Implement
}

func (r *LedgerRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.TransactionStatus) error {
	return nil // TODO: Implement
}
