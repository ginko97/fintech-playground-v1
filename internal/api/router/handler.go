package router

import (
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	// 🚀 1. Safely parse the incoming string into a uuid.UUID type
	parsedAccountID, err := uuid.Parse(req.AccountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid account ID format",
			"details": "The provided account_id is not a valid UUID",
		})
		return
	}

	txID, err := uuid.NewV7() // Or uuid.NewRandom() depending on your preference
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate transaction identifier",
		})
		return
	}

	// 🚀 2. Now construct the Domain Object with the correctly typed variable
	tx := &domain.Transaction{
		ID:             txID,            // 🚀 ASSIGN THE FRESH REAL UUID HERE!
		AccountID:      parsedAccountID, // Successfully matched type!
		Amount:         req.Amount,
		Currency:       req.Currency,
		IdempotencyKey: req.IdempotencyKey,
		Status:         domain.StatusPending,
		Type:           domain.TypeDebit,
	}

	ctx := c.Request.Context()

	// 3. Pass it cleanly to your Usecase
	processedTx, err := h.usecase.ProcessPayment(ctx, tx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Payment processing failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment executed and recorded successfully!",
		"status":  processedTx.Status,
		"id":      processedTx.ID,
	})
}
