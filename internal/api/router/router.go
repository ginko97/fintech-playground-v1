package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ginko97/fintech-playground-v1/internal/usecase"
)

func NewRouter(uc *usecase.TransactionUsecase) *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery())

	handler := NewTransactionHandler(uc)

	v1 := r.Group("/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "UP"})
		})
		v1.POST("/payments", handler.CreatePayment)

		v1.GET("/accounts/:id/balance", handler.GetBalance)
	}

	return r
}
