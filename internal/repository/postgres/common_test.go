package postgres

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func setupTestDB(t *testing.T, ctx context.Context) *pgxpool.Pool {
	// 1. Skip the Docker container startup for now
	// 2. Use your local Postgres connection string
	// Replace 'postgres', 'password', and '5432' with your actual local credentials
	connStr := "postgres://postgres:password@localhost:5432/fintech_test?sslmode=disable"

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("Failed to connect to local DB: %v. Make sure 'fintech_test' database is created!", err)
	}

	// 3. Setup Cleanup (only close the pool, no container to terminate)
	t.Cleanup(func() {
		pool.Close()
	})

	// 4. Run Migrations (Keep this part!)
	schemaPath := filepath.Join("..", "..", "..", "migrations", "000001_init_ledger.up.sql")
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("failed to read migration file: %v", err)
	}

	if _, err := pool.Exec(ctx, string(schema)); err != nil {
		t.Fatalf("failed to apply schema: %v", err)
	}

	return pool
}
