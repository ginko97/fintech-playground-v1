package usecase

import (
	"context"
	"testing"

	"github.com/ginko97/fintech-playground-v1/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// --- Manual Mocks ---

type mockRepo struct {
	onCreate func(tx *domain.Transaction) error
}

func (m *mockRepo) Create(ctx context.Context, tx *domain.Transaction) error {
	return m.onCreate(tx)
}

// Implement other methods with empty bodies to satisfy the interface
func (m *mockRepo) GetByInternalID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	return nil, nil
}
func (m *mockRepo) GetByIdempotencyKey(ctx context.Context, key string) (*domain.Transaction, error) {
	return nil, nil
}
func (m *mockRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.TransactionStatus) error {
	return nil
}

type mockGateway struct {
	onAuthorize func(tx domain.Transaction) (*domain.PaymentResponse, error)
}

func (m *mockGateway) Authorize(ctx context.Context, tx domain.Transaction) (*domain.PaymentResponse, error) {
	return m.onAuthorize(tx)
}
func (m *mockGateway) Inquiry(ctx context.Context, tx domain.Transaction) (*domain.PaymentResponse, error) {
	return nil, nil
}

// --- The Test ---

func TestProcessPayment_Success(t *testing.T) {
	// 1. Setup
	repo := &mockRepo{
		onCreate: func(tx *domain.Transaction) error {
			return nil // Simulate DB success
		},
	}

	gw := &mockGateway{
		onAuthorize: func(tx domain.Transaction) (*domain.PaymentResponse, error) {
			return &domain.PaymentResponse{
				Status: domain.StatusSuccess,
			}, nil // Simulate Bank success
		},
	}

	uc := NewTransactionUsecase(repo, gw)

	// 2. Execution
	tx, _ := domain.NewTransaction(uuid.New(), 1000, "SGD", domain.TypeDebit, "QR", "idem-123")
	result, err := uc.ProcessPayment(context.Background(), tx)

	// 3. Assertions
	assert.NoError(t, err)
	assert.Equal(t, domain.StatusSuccess, result.Status)
}
