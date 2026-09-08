package domain

import (
	"time"
	"github.com/google/uuid"
)

type PaymentEvent struct {
	ID         uuid.UUID `json:"id"`
	PaymentID  uuid.UUID `json:"payment_id"`
	Status     string    `json:"status"`
	OccurredAt time.Time `json:"occurred_at"`
}
