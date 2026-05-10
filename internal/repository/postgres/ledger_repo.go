package postgres

import (
	"context"
	"encoding/json"

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
	defer dbTx.Rollback(ctx)

	// 1. Insert Idempotency Key (The Gatekeeper)
	_, err = dbTx.Exec(ctx, `
        INSERT INTO idempotency_keys (key, request_path) 
        VALUES ($1, $2)`,
		tx.IdempotencyKey, "/v1/transactions")

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			// Duplicate detected: check if we already have a response saved
			var savedBody []byte
			err := dbTx.QueryRow(ctx, `
				SELECT response_body FROM idempotency_keys WHERE key = $1`,
				tx.IdempotencyKey).Scan(&savedBody)

			if err == nil && savedBody != nil {
				// We already processed this! Return the saved transaction.
				return json.Unmarshal(savedBody, tx)
			}
			return domain.ErrDuplicateRequest
		}
		return err
	}

	// 2. Insert Ledger Entry
	err = dbTx.QueryRow(ctx, `
        INSERT INTO ledger_entries (
            id, account_id, amount, currency, type, funding_source, status, idempotency_key
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING created_at, update_at`,
		tx.ID, tx.AccountID, tx.Amount, tx.Currency, tx.Type, tx.FundingSource, tx.Status, tx.IdempotencyKey,
	).Scan(&tx.CreatedAt, &tx.UpdatedAt)

	if err != nil {
		return err
	}

	// 3. Update Idempotency Key with the response
	responseBody, err := json.Marshal(tx)
	if err != nil {
		return err
	}

	_, err = dbTx.Exec(ctx, `
		UPDATE idempotency_keys 
		SET response_code = 201, response_body = $1 
		WHERE key = $2`,
		responseBody, tx.IdempotencyKey)

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
