package domain

import "context"

type PaymentResponse struct {
	ExternalID  string
	Status      TransactionStatus
	RawResponse string // Crucial for debugging "uncertain" cases
}

type PaymentGateway interface {
	Authorize(ctx context.Context, tx Transaction) (*PaymentResponse, error)

	Inquiry(ctx context.Context, tx Transaction) (*PaymentResponse, error)
}
