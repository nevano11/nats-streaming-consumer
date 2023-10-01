package entity

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Model struct {
	OrderUid          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerId        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmId              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

func (m *Model) String() string {
	result := strings.Builder{}
	result.WriteString("Model:{")

	result.WriteString(
		fmt.Sprintf("OrderUid: %s, TrackNumber: %s, Entry: %s,",
			m.OrderUid, m.TrackNumber, m.Entry))

	result.WriteString(m.Delivery.String())
	result.WriteString(",")

	result.WriteString(m.Payment.String())
	result.WriteString(",")

	result.WriteString("Items: [")
	for _, v := range m.Items {
		result.WriteString(v.String())
		result.WriteString(",")
	}
	result.WriteString("]")

	result.WriteString(",Locale: ")
	result.WriteString(m.Locale)
	result.WriteString(",InternalSignature: ")
	result.WriteString(m.InternalSignature)
	result.WriteString(",CustomerId: ")
	result.WriteString(m.CustomerId)
	result.WriteString(",DeliveryService: ")
	result.WriteString(m.DeliveryService)
	result.WriteString(",Shardkey: ")
	result.WriteString(m.Shardkey)
	result.WriteString(",SmId: ")
	result.WriteString(strconv.Itoa(m.SmId))
	result.WriteString(",DateCreated: ")
	result.WriteString(m.DateCreated.String())
	result.WriteString(",OofShard: ")
	result.WriteString(m.OofShard)
	result.WriteString("}")
	return result.String()
}
