package models

type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type NormalizerStruct struct {
	Max_amount float64
	Max_hour   float64
	Max_avg    float64
}

type TransactionSctruct struct {
	Amount              float64
	Hour                float64
	Customer_avg_amount float64
}

type ApproveTransactionDTO struct {
	Amount              float64 `json:"amount"`
	Hour                float64 `json:"hour"`
	Customer_avg_amount float64 `json:"customer_avg_amount"`
}
