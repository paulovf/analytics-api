package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/paulovf/analytics-api/internal/domain"
)

type PaymentRepo struct {
	db *pgxpool.Pool
}

func NewPaymentRepo(db *pgxpool.Pool) domain.PaymentRepository {
	return &PaymentRepo{db: db}
}

func (r *PaymentRepo) FindAll(ctx context.Context, limit, offset int) ([]domain.Payment, error) {
	query := `SELECT id, client_id, quantity, unit_price, total_amount, payment_method, created_at 
	          FROM payments ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []domain.Payment
	for rows.Next() {
		var p domain.Payment
		if err := rows.Scan(&p.ID, &p.ClientID, &p.Quantity, &p.UnitPrice, &p.TotalAmount, &p.PaymentMethod, &p.CreatedAt); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, rows.Err()
}

func (r *PaymentRepo) FindByID(ctx context.Context, id uuid.UUID) (domain.Payment, error) {
	query := `SELECT id, client_id, quantity, unit_price, total_amount, payment_method, created_at 
	          FROM payments WHERE id = $1`
	
	var p domain.Payment
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.ClientID, &p.Quantity, &p.UnitPrice, &p.TotalAmount, &p.PaymentMethod, &p.CreatedAt,
	)
	return p, err
}

func (r *PaymentRepo) FindEventsByPaymentID(ctx context.Context, paymentID uuid.UUID) ([]domain.PaymentEvent, error) {
	query := `SELECT id, payment_id, status, occurred_at 
	          FROM payment_events WHERE payment_id = $1 ORDER BY occurred_at ASC`
	
	rows, err := r.db.Query(ctx, query, paymentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []domain.PaymentEvent
	for rows.Next() {
		var e domain.PaymentEvent
		if err := rows.Scan(&e.ID, &e.PaymentID, &e.Status, &e.OccurredAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}
