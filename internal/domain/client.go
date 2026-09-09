package domain

import (
	"github.com/google/uuid"
	"time"
)

type Client struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CPF       string    `json:"cpf"`
	Address   string    `json:"address"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
