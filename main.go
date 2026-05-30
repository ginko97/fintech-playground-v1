package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/ginko97/fintech-playground-v1/internal/api/router"
	"github.com/ginko97/fintech-playground-v1/internal/domain"
	"github.com/ginko97/fintech-playground-v1/internal/repository/postgres"
	"github.com/ginko97/fintech-playground-v1/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LocalMockGateway struct{}

func (g *LocalMockGateway) Authorize(ctx context.Context, tx domain.Transaction) (*domain.PaymentResponse, error) {
	return &domain.PaymentResponse{
		ExternalID:  "mock_bank_ref_" + tx.IdempotencyKey, // Using ExternalID from gateway.go
		Status:      domain.StatusSuccess,                 // Using your domain's status enum/type
		RawResponse: `{"status":"approved","ref":"mock"}`,
	}, nil
}

func (g *LocalMockGateway) Inquiry(ctx context.Context, tx domain.Transaction) (*domain.PaymentResponse, error) {
	return &domain.PaymentResponse{
		ExternalID:  "mock_bank_ref_" + tx.IdempotencyKey,
		Status:      domain.StatusSuccess,
		RawResponse: `{"status":"approved","ref":"mock"}`,
	}, nil
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo, // Capture INFO, WARN, and ERROR logs
	}))

	// Set it globally so any file can access it using slog.Info() or slog.Error()
	slog.SetDefault(logger)

	slog.Info("Starting fintech-playground-v1 API server", "port", "8080")

	ctx := context.Background()

	// Database Connection Pool
	connStr := "postgres://postgres:password@localhost:5432/fintech_test?sslmode=disable"
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	// Layer Initialization
	repo := postgres.NewLedgerRepo(pool)
	mockGateway := &LocalMockGateway{}

	// Declare uc as the correct concrete pointer type
	var uc *usecase.TransactionUsecase = usecase.NewTransactionUsecase(repo, mockGateway)

	r := router.NewRouter(uc)

	log.Println("Server running on port :8080...")
	if err := r.Run(":8080"); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server crash error: %v", err)
	}
}
