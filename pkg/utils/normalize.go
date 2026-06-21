package utils

import (
	"fmt"
	"gin-test/internal/store"
	"gin-test/models"
	"slices"
	"time"
)

func NormalizeTransaction(transaction models.ApproveTransactionDTO) [14]float32 {
	var transactionVector [14]float32
	transactionTime, err := time.Parse(time.RFC3339, transaction.Transaction.RequestedAt)
	if err != nil {
		panic(err)
	}
	dayOfTheWeek := int(transactionTime.Weekday() - 1)
	if dayOfTheWeek == -1 {
		dayOfTheWeek = 6
	}
	transactionVector[0] = GetLimit(transaction.Transaction.Amount, float64(store.Normalizer.MaxAmount))
	transactionVector[1] = GetLimit(float64(transaction.Transaction.Installments), float64(store.Normalizer.MaxInstallments))
	if transaction.Customer.AvgAmount > 0 {
		transactionVector[2] = GetLimit(float64(transaction.Transaction.Amount/transaction.Customer.AvgAmount), float64(store.Normalizer.AmountVsAvgRatio))
	} else {
		transactionVector[2] = 0
	}
	transactionVector[3] = GetLimit(float64(transactionTime.Hour()), float64(23))
	transactionVector[4] = GetLimit(float64(dayOfTheWeek), float64(6))

	if transaction.LastTransaction == nil || transaction.LastTransaction.Timestamp == "" {
		transactionVector[5] = -1
		transactionVector[6] = -1
	} else {
		lastTransactionTime, err := time.Parse(time.RFC3339, transaction.LastTransaction.Timestamp)
		if err != nil {
			fmt.Println("Error parsing last transaction time", err)
		}
		minutesDiff := transactionTime.Sub(lastTransactionTime).Minutes()
		transactionVector[5] = GetLimit(minutesDiff, float64(store.Normalizer.MaxMinutes))
		transactionVector[6] = GetLimit(transaction.LastTransaction.KmFromCurrent, float64(store.Normalizer.MaxKm))
	}

	transactionVector[7] = GetLimit(transaction.Terminal.KmFromHome, float64(store.Normalizer.MaxKm))
	transactionVector[8] = GetLimit(transaction.Customer.TxCount24h, float64(store.Normalizer.MaxTxCount24h))

	transactionVector[9] = 0
	if transaction.Terminal.IsOnline {
		transactionVector[9] = 1
	}

	transactionVector[10] = 0
	if transaction.Terminal.CardPresent {
		transactionVector[10] = 1
	}

	/* transactionVector[11] should be 1 if the merchant is unkown */
	transactionVector[11] = 1
	if slices.Contains(transaction.Customer.KnowMerchants, transaction.Merchant.Id) {
		transactionVector[11] = 0
	}

	/* default value is 0.5 */
	transactionVector[12] = 0.5
	if mcc, ok := store.MccMap[transaction.Merchant.Mcc]; ok {
		transactionVector[12] = float32(mcc)
	}

	transactionVector[13] = GetLimit(transaction.Merchant.AvgAmount, float64(store.Normalizer.MaxMerchantAvgAmount))

	return transactionVector
}
