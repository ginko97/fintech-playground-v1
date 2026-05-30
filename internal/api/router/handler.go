package router

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ginko97/fintech-playground-v1/internal/domain"
	"github.com/ginko97/fintech-playground-v1/internal/usecase"
	"github.com/google/uuid"
)

// PaymentRequest defines the incoming JSON payload for a transaction
type PaymentRequest struct {
	AccountID      string `json:"account_id" binding:"required,uuid"`
	Amount         int64  `json:"amount" binding:"required,gt=0"`
	Currency       string `json:"currency" binding:"required,len=3"`
	IdempotencyKey string `json:"idempotency_key" binding:"required"`
}

type TransactionHandler struct {
	usecase *usecase.TransactionUsecase
}

func NewTransactionHandler(uc *usecase.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{
		usecase: uc,
	}
}

// CreatePayment handles POST /v1/payments
func (h *TransactionHandler) CreatePayment(c *gin.Context) {
	var req PaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Failed to bind payment request JSON", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	parsedAccountID, err := uuid.Parse(req.AccountID)
	if err != nil {
		slog.Warn("Invalid account ID format received", "account_id", req.AccountID)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid account ID format",
			"details": "The provided account_id is not a valid UUID",
		})
		return
	}

	txID, err := uuid.NewV7() // Or uuid.NewRandom() depending on your preference
	if err != nil {
		slog.Error("Failed to generate transaction UUID", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate transaction identifier",
		})
		return
	}

	tx := &domain.Transaction{
		ID:             txID,
		AccountID:      parsedAccountID,
		Amount:         req.Amount,
		Currency:       req.Currency,
		IdempotencyKey: req.IdempotencyKey,
		Status:         domain.StatusPending,
		Type:           domain.TypeDebit,
	}

	slog.Info("Processing payment request received at API edge",
		"tx_id", tx.ID.String(),
		"account_id", tx.AccountID.String(),
		"amount", tx.Amount,
		"currency", tx.Currency,
		"idempotency_key", tx.IdempotencyKey,
	)

	ctx := c.Request.Context()

	// 3. Pass it cleanly to your Usecase
	processedTx, err := h.usecase.ProcessPayment(ctx, tx)
	if err != nil {
		slog.Error("Payment orchestration pipeline failed",
			"tx_id", tx.ID.String(),
			"idempotency_key", tx.IdempotencyKey,
			"error", err.Error(),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Payment processing failed",
			"details": err.Error(),
		})
		return
	}
	slog.Info("Payment request finalized and recorded successfully",
		"tx_id", processedTx.ID.String(),
		"status", string(processedTx.Status),
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment executed and recorded successfully!",
		"status":  processedTx.Status,
		"id":      processedTx.ID,
	})
}

// GET /v1/accounts/:id/balance
func (h *TransactionHandler) GetBalance(c *gin.Context) {
	idStr := c.Param("id")
	parsedAccountID, err := uuid.Parse(idStr)
	if err != nil {
		slog.Warn("Invalid account ID lookup format", "received_id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID format"})
		return
	}

	ctx := c.Request.Context()
	balance, err := h.usecase.GetAccountBalance(ctx, parsedAccountID)
	if err != nil {
		slog.Error("Failed to calculate account balance", "account_id", idStr, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve balance"})
		return
	}

	slog.Info("Account balance fetched successfully", "account_id", idStr, "balance", balance)

	c.JSON(http.StatusOK, gin.H{
		"account_id": parsedAccountID,
		"balance":    balance,
		"currency":   "USD", // Default currency for tracking
	})
}
