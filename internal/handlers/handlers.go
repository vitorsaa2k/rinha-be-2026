package handlers

import (
	"gin-test/internal/ivf"
	"gin-test/models"
	"gin-test/pkg/utils"
	"math"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

type ApproveTransactionResponse struct {
	Approved   bool    `json:"approved"`
	FraudScore float64 `json:"fraud_score"`
}

func ApproveTransaction(c fiber.Ctx) error {
	transaction := models.ApproveTransactionDTO{}
	if err := c.Bind().Body(&transaction); err != nil {
		return err
	}
	normalizedTransaction := utils.NormalizeTransaction(transaction)
	dataset := ivf.SearchIVF(normalizedTransaction[:], 5, 5)
	result, err := utils.SearchInVector(normalizedTransaction[:], 14, dataset)
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}
	score := math.Round(result.Score*100) / 100
	switch result.IsPossibleFraud {
	case true:
		return c.JSON(ApproveTransactionResponse{Approved: false, FraudScore: score})
	case false:
		return c.JSON(ApproveTransactionResponse{Approved: true, FraudScore: score})
	default:
		return c.JSON(ApproveTransactionResponse{Approved: false, FraudScore: score})
	}
}

func Ready(c fiber.Ctx) error {
	return c.SendStatus(http.StatusOK)
}
