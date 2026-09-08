package domain

import (
	"time"
	"github.com/google/uuid"
)

type Payment struct {
	ID            uuid.UUID `json:"id"`
	ClientID      uuid.UUID `json:"client_id"`
	Quantity      int       `json:"quantity"`
	UnitPrice     float64   `json:"unit_price"`
	TotalAmount   float64   `json:"total_amount"`
	PaymentMethod string    `json:"payment_method"`
	CreatedAt     time.Time `json:"created_at"`
}
