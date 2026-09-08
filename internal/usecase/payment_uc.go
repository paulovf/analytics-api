package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/paulovf/analytics-api/internal/domain"
)

type PaymentUseCase struct {
	repo domain.PaymentRepository
}

func NewPaymentUseCase(repo domain.PaymentRepository) *PaymentUseCase {
	return &PaymentUseCase{repo: repo}
}

func (uc *PaymentUseCase) ListPayments(ctx context.Context, limit, page int) ([]domain.Payment, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	return uc.repo.FindAll(ctx, limit, offset)
}

func (uc *PaymentUseCase) GetPaymentByID(ctx context.Context, id uuid.UUID) (domain.Payment, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *PaymentUseCase) GetPaymentEvents(ctx context.Context, paymentID uuid.UUID) ([]domain.PaymentEvent, error) {
	return uc.repo.FindEventsByPaymentID(ctx, paymentID)
}
