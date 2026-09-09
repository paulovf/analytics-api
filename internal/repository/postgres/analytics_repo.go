package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/paulovf/analytics-api/internal/domain"
)

type AnalyticsRepo struct {
	db *pgxpool.Pool
}

func NewAnalyticsRepo(db *pgxpool.Pool) domain.AnalyticsRepository {
	return &AnalyticsRepo{db: db}
}

func (r *AnalyticsRepo) GetDailyStatusStats(ctx context.Context, startDate, endDate time.Time) ([]domain.DailyStatusStat, error) {
	query := `
		SELECT day_bucket, status, total_events 
		FROM payment_status_daily 
		WHERE day_bucket >= $1 AND day_bucket <= $2
		ORDER BY day_bucket DESC, status ASC
	`

	rows, err := r.db.Query(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []domain.DailyStatusStat
	for rows.Next() {
		var s domain.DailyStatusStat
		if err := rows.Scan(&s.DayBucket, &s.Status, &s.TotalEvents); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}

	return stats, rows.Err()
}
