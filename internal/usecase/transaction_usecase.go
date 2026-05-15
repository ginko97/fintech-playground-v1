package usecase

import (
	"context"
	"errors"
	"log"

	"github.com/ginko97/fintech-playground-v1/internal/domain"
)

type TransactionUsecase struct {
	repo    domain.TransactionRepository
	gateway domain.PaymentGateway
}

func NewTransactionUsecase(repo domain.TransactionRepository, gateway domain.PaymentGateway) *TransactionUsecase {
	return &TransactionUsecase{
		repo:    repo,
		gateway: gateway,
	}
}

func (u *TransactionUsecase) ProcessPayment(ctx context.Context, tx *domain.Transaction) (*domain.Transaction, error) {
	// 1. Save to Ledger as Pending (using your Idempotent Repo)
	err := u.repo.Create(ctx, tx)
	if err != nil {
		// If it's a duplicate, your repo returns the saved result.
		// If it's a real error, we stop here.
		return nil, err
	}

	// 2. Call External Gateway
	resp, err := u.gateway.Authorize(ctx, *tx)
	if err != nil {
		// Check for TIMEOUT
		if errors.Is(err, context.DeadlineExceeded) {
			log.Printf("Transaction %s timed out. Marking as UNCERTAIN.", tx.ID)
			u.repo.UpdateStatus(ctx, tx.ID, domain.StatusUncertain)
			tx.Status = domain.StatusUncertain
			return tx, nil
		}

		// Other errors mark as FAILED
		u.repo.UpdateStatus(ctx, tx.ID, domain.StatusFailed)
		tx.Status = domain.StatusFailed
		return tx, err
	}

	// 3. Update Ledger with Success
	err = u.repo.UpdateStatus(ctx, tx.ID, resp.Status)
	tx.Status = resp.Status

	return tx, err
}
