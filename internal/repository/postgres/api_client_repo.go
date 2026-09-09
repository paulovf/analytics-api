package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/paulovf/analytics-api/internal/domain"
)

type ApiClientRepo struct {
	db *pgxpool.Pool
}

func NewApiClientRepo(db *pgxpool.Pool) domain.ApiClientRepository {
	return &ApiClientRepo{db: db}
}

func (r *ApiClientRepo) FindByClientID(ctx context.Context, clientID string) (*domain.ApiClient, error) {
	query := `SELECT id, name, client_id, client_secret_hash, type, status, created_at, updated_at 
	          FROM api_clients WHERE client_id = $1 AND status = 'active'`

	var c domain.ApiClient
	err := r.db.QueryRow(ctx, query, clientID).Scan(
		&c.ID, &c.Name, &c.ClientID, &c.ClientSecretHash, &c.Type, &c.Status, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ApiClientRepo) Create(ctx context.Context, client *domain.ApiClient) error {
	query := `INSERT INTO api_clients (id, name, client_id, client_secret_hash, type, status, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.Exec(ctx, query,
		client.ID, client.Name, client.ClientID, client.ClientSecretHash, client.Type, client.Status, client.CreatedAt, client.UpdatedAt,
	)
	return err
}
