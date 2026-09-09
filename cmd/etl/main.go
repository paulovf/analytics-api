package main

import (
	"context"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/paulovf/analytics-api/pkg/postgres"
)

func main() {
	ctx := context.Background()

	dsnPG := "postgres://analytics_user:analytics_password@localhost:5432/analytics_db?sslmode=disable"
	pgPool, err := postgres.NewPool(ctx, dsnPG)
	if err != nil {
		log.Fatalf("Error during connect on postgres: %v", err)
	}
	defer pgPool.Close()

	chConn, err := connectClickHouse()
	if err != nil {
		log.Fatalf("Error during connect on ClickHouse: %v", err)
	}
	defer chConn.Close()

	if err := createClickHouseTable(ctx, chConn); err != nil {
		log.Fatalf("Error creating table in ClickHouse: %v", err)
	}

	if err := runETL(ctx, pgPool, chConn); err != nil {
		log.Fatalf("Error during ETL: %v", err)
	}
}

func runETL(ctx context.Context, pgPool *pgxpool.Pool, chConn driver.Conn) error {
	log.Println("Starting ETL events: TimescaleDB -> ClickHouse...")
	startETL := time.Now()

	var totalEvents int
	if err := pgPool.QueryRow(ctx, "SELECT COUNT(*) FROM payment_events").Scan(&totalEvents); err != nil {
		return err
	}
	log.Printf("Total of records to be migrated: %d", totalEvents)

	rows, err := pgPool.Query(ctx, "SELECT id, payment_id, status, occurred_at FROM payment_events")
	if err != nil {
		return err
	}
	defer rows.Close()

	batch, err := chConn.PrepareBatch(ctx, "INSERT INTO payment_events_ch (id, payment_id, status, occurred_at)")
	if err != nil {
		return err
	}

	count, err := migratePaymentEvents(ctx, chConn, rows, batch, totalEvents)
	if err != nil {
		return err
	}

	log.Printf("ETL Completed! %d records copied in %v", count, time.Since(startETL))
	return nil
}

func migratePaymentEvents(
	ctx context.Context,
	chConn driver.Conn,
	rows pgx.Rows,
	batch driver.Batch,
	totalEvents int,
) (int, error) {
	count := 0
	for rows.Next() {
		var id, paymentID string
		var status string
		var occurredAt time.Time

		if err := rows.Scan(&id, &paymentID, &status, &occurredAt); err != nil {
			return count, err
		}

		if err := batch.Append(id, paymentID, status, occurredAt); err != nil {
			return count, err
		}

		count++
		if count%100000 == 0 {
			if err := batch.Send(); err != nil {
				return count, err
			}
			log.Printf("Migrated %d/%d records...", count, totalEvents)
			var err error
			batch, err = chConn.PrepareBatch(ctx, "INSERT INTO payment_events_ch (id, payment_id, status, occurred_at)")
			if err != nil {
				return count, err
			}
		}
	}

	if err := batch.Send(); err != nil {
		return count, err
	}

	return count, nil
}

func connectClickHouse() (driver.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"localhost:9000"},
		Auth: clickhouse.Auth{
			Database: "analytics_db",
			Username: "analytics_user",
			Password: "analytics_password",
		},
		Debug: false,
	})
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}
	return conn, nil
}

func createClickHouseTable(ctx context.Context, conn driver.Conn) error {
	query := `
		CREATE TABLE IF NOT EXISTS payment_events_ch (
			id UUID,
			payment_id UUID,
			status String,
			occurred_at DateTime
		) ENGINE = MergeTree()
		ORDER BY (status, occurred_at)
	`
	return conn.Exec(ctx, query)
}
