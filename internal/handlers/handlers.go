package handlers

import (
	"gin-test/internal/ivf"
	"gin-test/models"
	"gin-test/pkg/utils"

	"github.com/goccy/go-json"

	"github.com/valyala/fasthttp"
)

type ApproveTransactionResponse struct {
	Approved   bool    `json:"approved"`
	FraudScore float64 `json:"fraud_score"`
}

var ResponseRefuse1, _ = json.Marshal(ApproveTransactionResponse{Approved: false, FraudScore: 1})
var ResponseRefuse08, _ = json.Marshal(ApproveTransactionResponse{Approved: false, FraudScore: 0.8})
var ResponseRefuse06, _ = json.Marshal(ApproveTransactionResponse{Approved: false, FraudScore: 0.6})
var ResponseApproved0, _ = json.Marshal(ApproveTransactionResponse{Approved: true, FraudScore: 0})
var ResponseApproved02, _ = json.Marshal(ApproveTransactionResponse{Approved: true, FraudScore: 0.2})
var ResponseApproved04, _ = json.Marshal(ApproveTransactionResponse{Approved: true, FraudScore: 0.4})

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
		err := json.Unmarshal(ctx.Request.Body(), &transaction)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.SetBodyString("Invalid JSON structure")
			return
		}
		normalizedTransaction := utils.NormalizeTransaction(transaction)
		//closestCentroids := ivf.SearchIVF(normalizedTransaction[:], 5)
		quantizedQuery := ivf.QuantizeQuery(normalizedTransaction)
		closestCentroids := ivf.SearchIVFQuantized(quantizedQuery, 5)
		result, err := utils.SearchInVectorQuantized(quantizedQuery, 14, closestCentroids)
		//result, err := utils.SearchInVector(normalizedTransaction[:], 14, closestCentroids)
		if err != nil {
			ctx.Error("Error when searching in vector", fasthttp.StatusBadRequest)
			return
		}
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(fasthttp.StatusOK)
		switch result.Score {
		case 0:
			ctx.SetBody(ResponseApproved0)
		case 0.2:
			ctx.SetBody(ResponseApproved02)
		case 0.4:
			ctx.SetBody(ResponseApproved04)
		case 0.6:
			ctx.SetBody(ResponseRefuse06)
		case 0.8:
			ctx.SetBody(ResponseRefuse08)
		default:
			ctx.SetBody(ResponseRefuse1)
		}

	default:
		// Handle 404 Not Found
		ctx.Error("Page Not Found", fasthttp.StatusNotFound)
	}
}
