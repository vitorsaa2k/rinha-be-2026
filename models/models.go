package models

type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type NormalizerStruct struct {
	MaxAmount            uint16 `json:"max_amount"`
	MaxInstallments      uint8  `json:"max_installments"`
	AmountVsAvgRatio     uint8  `json:"amount_vs_avg_ratio"`
	MaxMinutes           uint16 `json:"max_minutes"`
	MaxKm                uint16 `json:"max_km"`
	MaxTxCount24h        uint8  `json:"max_tx_count_24h"`
	MaxMerchantAvgAmount uint16 `json:"max_merchant_avg_amount"`
}

type MccRisk struct {
	Amount              float64
	Hour                float64
	Customer_avg_amount float64
}
type ApproveTransactionDTO struct {
	Id              string           `json:"id"`
	Transaction     Transaction      `json:"transaction"`
	Customer        Customer         `json:"customer"`
	Merchant        Merchant         `json:"merchant"`
	Terminal        Terminal         `json:"terminal"`
	LastTransaction *LastTransaction `json:"last_transaction"`
}

type Transaction struct {
	Amount       float64 `json:"amount"`
	Installments uint8   `json:"installments"`
	RequestedAt  string  `json:"requested_at"`
}

type Customer struct {
	AvgAmount     float64  `json:"avg_amount"`
	TxCount24h    float64  `json:"tx_count_24h"`
	KnowMerchants []string `json:"known_merchants"`
}

type Merchant struct {
	Id        string  `json:"id"`
	Mcc       string  `json:"mcc"`
	AvgAmount float64 `json:"avg_amount"`
}

type Terminal struct {
	IsOnline    bool    `json:"is_online"`
	CardPresent bool    `json:"card_present"`
	KmFromHome  float64 `json:"km_from_home"`
}

type LastTransaction struct {
	Timestamp     string  `json:"timestamp"`
	KmFromCurrent float64 `json:"km_from_current"`
}

type DatasetStruct struct {
	Vector []float32 `json:"vector"`
	Label  string    `json:"label"`
}

type QuantizedData struct {
	Vector []int16
	Label  uint8
}
