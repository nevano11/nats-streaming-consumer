package entity

import "fmt"

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestId    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

func (p *Payment) String() string {
	return fmt.Sprintf(
		"Payment: {Transaction: %s, RequestId: %s, Currency: %s, Provider: %s, Amount: %d, PaymentDt: %d, Bank: %s, DeliveryCost: %d, GoodsTotal: %d, CustomFee: %d}",
		p.Transaction,
		p.RequestId,
		p.Currency,
		p.Provider,
		p.Amount,
		p.PaymentDt,
		p.Bank,
		p.DeliveryCost,
		p.GoodsTotal,
		p.CustomFee)
}
