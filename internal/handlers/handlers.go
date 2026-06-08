package handlers

import (
	"encoding/json"
	"gin-test/internal/ivf"
	"gin-test/models"
	"gin-test/pkg/utils"

	"github.com/valyala/fasthttp"
)

type ApproveTransactionResponse struct {
	Approved   bool    `json:"approved"`
	FraudScore float64 `json:"fraud_score"`
}

/* func ApproveTransaction(c fiber.Ctx) error {
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
} */

func FastHTTPHandler(ctx *fasthttp.RequestCtx) {
	// 1. Get the requested path
	path := string(ctx.Path())

	// 2. Simple route handling
	switch path {
	case "/ready":
		ctx.SetStatusCode(fasthttp.StatusOK)

	case "/json":
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(fasthttp.StatusOK)

		response := ApproveTransactionResponse{Approved: false, FraudScore: 1}

		// Encode and write the response
		if err := json.NewEncoder(ctx).Encode(response); err != nil {
			ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		}

	case "/fraud-score":
		transaction := models.ApproveTransactionDTO{}
		response := ApproveTransactionResponse{Approved: false, FraudScore: 1}
		err := json.Unmarshal(ctx.Request.Body(), &transaction)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.SetBodyString("Invalid JSON structure")
			return
		}
		normalizedTransaction := utils.NormalizeTransaction(transaction)
		dataset := ivf.SearchIVF(normalizedTransaction[:], 5)
		result, err := utils.SearchInVector(normalizedTransaction[:], 14, dataset)
		if err != nil {
			ctx.Error("Error when searching in vector", fasthttp.StatusBadRequest)
			return
		}
		if result.IsPossibleFraud {
			ctx.SetContentType("application/json")
			ctx.SetStatusCode(fasthttp.StatusOK)
			response = ApproveTransactionResponse{Approved: false, FraudScore: result.Score}
			if err := json.NewEncoder(ctx).Encode(response); err != nil {
				ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
			}
		} else {
			ctx.SetContentType("application/json")
			ctx.SetStatusCode(fasthttp.StatusOK)
			response = ApproveTransactionResponse{Approved: true, FraudScore: result.Score}
			if err := json.NewEncoder(ctx).Encode(response); err != nil {
				ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
			}

		}

	default:
		// Handle 404 Not Found
		ctx.Error("Page Not Found", fasthttp.StatusNotFound)
	}
}
