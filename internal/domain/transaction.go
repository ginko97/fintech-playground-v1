package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrDuplicateRequest = errors.New("duplicate request")
)

// Define types for better type safety in business logic
type TransactionStatus string
type TransactionType string

const (
	StatusPending   TransactionStatus = "pending"
	StatusSuccess   TransactionStatus = "success"
	StatusFailed    TransactionStatus = "failed"
	StatusUncertain TransactionStatus = "uncertain"
	StatusReversed  TransactionStatus = "reversed"

	TypeDebit  TransactionType = "debit"
	TypeCredit TransactionType = "credit"
)

// Metadata is a key-value map for extensible payment data
type Metadata map[string]any

// MarshalJSON is a custom "Security Shield."
// It ensures that sensitive keys are never written to the DB or Logs in plain text.
func (m Metadata) MarshalJSON() ([]byte, error) {
	type alias Metadata
	copy := make(alias)

	// List of sensitive keys we want to hide
	sensitiveKeys := map[string]bool{
		"card_number": true,
		"cvv":         true,
		"otp":         true,
		"pin":         true,
	}

	for k, v := range m {
		if sensitiveKeys[strings.ToLower(k)] {
			copy[k] = "[REDACTED]"
		} else {
			copy[k] = v
		}
	}

	return json.Marshal(copy)
}

type Transaction struct {
	ID             uuid.UUID         `json:"id"`
	AccountID      uuid.UUID         `json:"account_id"`
	Amount         int64             `json:"amount"`
	Currency       string            `json:"currency"`
	Type           TransactionType   `json:"type"`
	FundingSource  string            `json:"funding_source"`
	Metadata       Metadata          `json:"metadata"`
	Status         TransactionStatus `json:"status"`
	IdempotencyKey string            `json:"idempotency_key"`
	SettledAt      *time.Time        `json:"settled_at,omitempty"` // Use pointer for nullability
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

func NewTransaction(accountID uuid.UUID, amount int64, currency string, txType TransactionType, fundingSource string, idempotencyKey string) (*Transaction, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return &Transaction{
		ID:             id,
		AccountID:      accountID,
		Amount:         amount,
		Currency:       currency,
		Type:           txType,
		FundingSource:  fundingSource,
		Status:         StatusPending,
		IdempotencyKey: idempotencyKey,
	}, nil
}
