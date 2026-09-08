package domain

import (
	"time"
	"github.com/google/uuid"
)

type ApiClient struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	ClientID         string    `json:"client_id"`
	ClientSecretHash string    `json:"-"`
	Type             string    `json:"type"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
