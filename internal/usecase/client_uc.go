package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/paulovf/analytics-api/internal/domain"
)

type ClientUseCase struct {
	repo domain.ClientRepository
}

func NewClientUseCase(repo domain.ClientRepository) *ClientUseCase {
	return &ClientUseCase{repo: repo}
}

func (uc *ClientUseCase) ListClients(ctx context.Context, limit, page int) ([]domain.Client, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	return uc.repo.FindAll(ctx, limit, offset)
}

func (uc *ClientUseCase) GetClientByID(ctx context.Context, id uuid.UUID) (domain.Client, error) {
	return uc.repo.FindByID(ctx, id)
}
