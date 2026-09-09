package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ClientRepository interface {
	FindAll(ctx context.Context, limit, offset int) ([]Client, error)
	FindByID(ctx context.Context, id uuid.UUID) (Client, error)
}

type PaymentRepository interface {
	FindAll(ctx context.Context, limit, offset int) ([]Payment, error)
	FindByID(ctx context.Context, id uuid.UUID) (Payment, error)
	FindEventsByPaymentID(ctx context.Context, paymentID uuid.UUID) ([]PaymentEvent, error)
}

type ApiClientRepository interface {
	FindByClientID(ctx context.Context, clientID string) (*ApiClient, error)
	Create(ctx context.Context, client *ApiClient) error
}

type AnalyticsRepository interface {
	GetDailyStatusStats(ctx context.Context, startDate, endDate time.Time) ([]DailyStatusStat, error)
}
