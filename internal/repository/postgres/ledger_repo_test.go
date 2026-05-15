package postgres

import (
	"context"
	"testing"

	"github.com/ginko97/fintech-playground-v1/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestLedgerRepo_Create_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	// NOTE: In a real project, we use a helper to start the container.
	// For now, assume 'pool' is a connection to a local test DB.
	// We will implement the TestContainer helper in the next step.
	pool := setupTestDB(t, ctx)
	defer pool.Close()

	repo := NewLedgerRepo(pool)

	tx, _ := domain.NewTransaction(uuid.New(), 5000, "SGD", domain.TypeDebit, "CRYPTO", "idem-789")

	// Test the actual DB Insert
	err := repo.Create(ctx, tx)

	assert.NoError(t, err)
	assert.NotZero(t, tx.CreatedAt) // Proves the DB generated the time!
}
