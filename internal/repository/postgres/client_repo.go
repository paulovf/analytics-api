package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/paulovf/analytics-api/internal/domain"
)

type ClientRepo struct {
	db *pgxpool.Pool
}

func NewClientRepo(db *pgxpool.Pool) domain.ClientRepository {
	return &ClientRepo{db: db}
}

func (r *ClientRepo) FindAll(ctx context.Context, limit, offset int) ([]domain.Client, error) {
	query := `SELECT id, name, cpf, address, status, created_at, updated_at 
	          FROM clients ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []domain.Client
	for rows.Next() {
		var c domain.Client
		if err := rows.Scan(&c.ID, &c.Name, &c.CPF, &c.Address, &c.Status, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}
	return clients, rows.Err()
}

func (r *ClientRepo) FindByID(ctx context.Context, id uuid.UUID) (domain.Client, error) {
	query := `SELECT id, name, cpf, address, status, created_at, updated_at 
	          FROM clients WHERE id = $1`
	
	var c domain.Client
	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.CPF, &c.Address, &c.Status, &c.CreatedAt, &c.UpdatedAt,
	)
	return c, err
}
