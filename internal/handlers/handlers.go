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

func ApproveTransaction(c fiber.Ctx) error {
	transaction := models.ApproveTransactionDTO{}
	if err := c.Bind().Body(&transaction); err != nil {
		return err
	}
	normalizedTransaction := utils.NormalizeTransaction(transaction)
	fmt.Println(normalizedTransaction)

	/* transaction := models.ApproveTransactionDTO{}
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
	} */
	/* TODO Add nearest neighboor votes for checking if it is a fraud or not */
	/* switch result.IsPossibleFraud {
	case true:
		return c.JSON(ApproveTransactionResponse{Approved: false, FraudScore: 1})
	default:
		return c.JSON(ApproveTransactionResponse{Approved: false, FraudScore: 0})
	} */
	return c.JSON(normalizedTransaction)
}

func LoadJson(c fiber.Ctx) error {
	/* var references []models.DatasetStruct
	path := filepath.Join("public", "references.json")
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return fiber.NewError(http.StatusInternalServerError, "Internal server error")
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&references); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return fiber.NewError(http.StatusInternalServerError, "Internal server error")
	}
	fmt.Printf("Loaded %d references", len(references)) */
	/* if len(store.References) < 1 {
		store.LoadReferencesStreamed()
	} */
	return c.JSON(store.References[12])
}

func Ready(c fiber.Ctx) error {
	return c.SendStatus(http.StatusOK)
}
