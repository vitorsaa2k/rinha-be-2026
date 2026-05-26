package handlers

import (
	"fmt"
	"gin-test/internal/store"
	"gin-test/models"
	"gin-test/pkg/utils"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

type ApproveTransactionResponse struct {
	Approved   bool    `json:"approved"`
	FraudScore float64 `json:"fraud_score"`
}

func Ready(c fiber.Ctx) error {
	return c.SendStatus(http.StatusOK)
}

func ApproveTransaction(c fiber.Ctx) error {
	transaction := models.ApproveTransactionDTO{}
	if err := c.Bind().Body(&transaction); err != nil {
		return err
	}
	fmt.Println(transaction)
	first := utils.GetLimit(transaction.Amount, store.Normalizer.Max_amount)
	second := utils.GetLimit(transaction.Hour, store.Normalizer.Max_hour)
	third := utils.GetLimit(transaction.Customer_avg_amount, store.Normalizer.Max_avg)
	vector := []float64{first, second, third}
	result, err := utils.SearchInVector(vector, 3, store.Dataset)
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}
	switch result.IsPossibleFraud {
	case true:
		return c.JSON(ApproveTransactionResponse{Approved: false, FraudScore: 1})
	default:
		return c.JSON(ApproveTransactionResponse{Approved: false, FraudScore: 0})
	}
}
