package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/paulovf/analytics-api/pkg/postgres"
)

const runs = 20

func main() {
	ctx := context.Background()

	dsnPG := "postgres://analytics_user:analytics_password@localhost:5432/analytics_db?sslmode=disable"
	pgPool, err := postgres.NewPool(ctx, dsnPG)
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}
	defer pgPool.Close()

	chConn, err := connectClickHouse()
	if err != nil {
		log.Fatalf("Failed to connect to ClickHouse: %v", err)
	}
	defer chConn.Close()

	fmt.Println("=====================================================")
	fmt.Printf("🚀 STARTING ANALYTICAL BENCHMARK (%d runs)\n", runs)
	fmt.Println("=====================================================")

	// =====================================================================
	// TEST 1: Simple Aggregation (Total events by Status)
	// =====================================================================
	fmt.Println("\n[TEST 1] Event Count Grouped by Status")

	pgQuery1 := "SELECT status, COUNT(*) FROM payment_events GROUP BY status"
	runPostgresBenchmark(ctx, pgPool, "TimescaleDB", pgQuery1)

	chQuery1 := "SELECT status, COUNT(*) FROM payment_events_ch GROUP BY status"
	runClickHouseBenchmark(ctx, chConn, "ClickHouse", chQuery1)

	// =====================================================================
	// TEST 2: Time Aggregation (Total events by Day and Status)
	// =====================================================================
	fmt.Println("\n[TEST 2] Time Aggregation (Grouped by Day and Status)")

	pgQuery2 := `
		SELECT date_trunc('day', occurred_at) as day, status, COUNT(*) 
		FROM payment_events 
		GROUP BY day, status 
		ORDER BY day DESC`
	runPostgresBenchmark(ctx, pgPool, "TimescaleDB", pgQuery2)

	chQuery2 := `
		SELECT toStartOfDay(occurred_at) as day, status, COUNT(*) 
		FROM payment_events_ch 
		GROUP BY day, status 
		ORDER BY day DESC`
	runClickHouseBenchmark(ctx, chConn, "ClickHouse", chQuery2)

	fmt.Println("\n=====================================================")
	fmt.Println("🏁 BENCHMARK FINISHED")
	fmt.Println("=====================================================")
}

func runPostgresBenchmark(ctx context.Context, pool *pgxpool.Pool, name, query string) {
	var totalDuration time.Duration

	for i := 0; i < runs; i++ {
		start := time.Now()
		rows, err := pool.Query(ctx, query)
		if err != nil {
			log.Fatalf("Error in %s: %v", name, err)
		}

		for rows.Next() {
		}
		rows.Close()
		totalDuration += time.Since(start)
	}

	avg := totalDuration / runs
	fmt.Printf("▶ %-12s : %v (avg)\n", name, avg)
}

func runClickHouseBenchmark(ctx context.Context, conn driver.Conn, name, query string) {
	var totalDuration time.Duration

	for i := 0; i < runs; i++ {
		start := time.Now()
		rows, err := conn.Query(ctx, query)
		if err != nil {
			log.Fatalf("Error in %s: %v", name, err)
		}
		for rows.Next() {
		}
		rows.Close()
		totalDuration += time.Since(start)
	}

	avg := totalDuration / runs
	fmt.Printf("▶ %-12s : %v (avg)\n", name, avg)
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
	return conn, conn.Ping(context.Background())
}
