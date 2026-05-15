package gateway

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/ginko97/fintech-playground-v1/internal/domain"
)

type MockGateway struct{}

func NewMockGateway() *MockGateway {
	return &MockGateway{}
}

func (g *MockGateway) Authorize(ctx context.Context, tx domain.Transaction) (*domain.PaymentResponse, error) {
	// Simulate "Real World" Latency (0.5s to 4s)
	delay := time.Duration(500+rand.Intn(3500)) * time.Millisecond

	select {
	case <-time.After(delay):
		// 10% chance of a random provider error
		if rand.Float32() < 0.1 {
			return nil, fmt.Errorf("provider_internal_error: switcher unavailable")
		}

		return &domain.PaymentResponse{
			ExternalID:  "MOCK_BANK_" + tx.ID.String()[:8],
			Status:      domain.StatusSuccess,
			RawResponse: `{"status": "APPROVED", "code": "00"}`,
		}, nil

	case <-ctx.Done():
		// The Context timed out! This leads to the UNCERTAIN state.
		return nil, ctx.Err()
	}
}

func (g *MockGateway) Inquiry(ctx context.Context, tx domain.Transaction) (*domain.PaymentResponse, error) {
	// Inquiries are usually faster but can still fail
	return &domain.PaymentResponse{
		ExternalID:  "MOCK_BANK_" + tx.ID.String()[:8],
		Status:      domain.StatusSuccess,
		RawResponse: `{"status": "APPROVED", "code": "00", "reconciled": true}`,
	}, nil
}
