package domain

import (
	"github.com/google/uuid"
	"time"
)

type PaymentEvent struct {
	ID         uuid.UUID `json:"id"`
	PaymentID  uuid.UUID `json:"payment_id"`
	Status     string    `json:"status"`
	OccurredAt time.Time `json:"occurred_at"`
}
