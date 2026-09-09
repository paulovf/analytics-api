package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/paulovf/analytics-api/pkg/postgres"
)

const (
	numClients  = 10000
	numPayments = 50000
)

func main() {
	ctx := context.Background()
	dsn := "postgres://analytics_user:analytics_password@localhost:5432/analytics_db?sslmode=disable"

	pool, err := postgres.NewPool(ctx, dsn)
	if err != nil {
		log.Fatalf("Fail in connect: %v", err)
	}
	defer pool.Close()

	gofakeit.Seed(0)
	log.Println("Starting massive data generator...")

	start := time.Now()
	var clientIDs []uuid.UUID
	var clientRows [][]any

	for i := 0; i < numClients; i++ {
		id := uuid.New()
		clientIDs = append(clientIDs, id)

		status := "active"
		if rand.Float64() > 0.88 {
			status = "inactive"
		}

		clientRows = append(clientRows, []any{
			id,
			gofakeit.Name(),
			gofakeit.SSN(),
			gofakeit.Address().Address,
			status,
			time.Now(),
			time.Now(),
		})
	}

	_, err = pool.CopyFrom(ctx, pgx.Identifier{"clients"},
		[]string{"id", "name", "cpf", "address", "status", "created_at", "updated_at"},
		pgx.CopyFromRows(clientRows))
	if err != nil {
		log.Fatalf("Error on insert clients: %v", err)
	}
	log.Printf("%d Clients inserted in %v", numClients, time.Since(start))

	start = time.Now()
	var paymentRows [][]any
	var eventRows [][]any
	methods := []string{"credit_card", "pix", "boleto", "debit_card"}

	for i := 0; i < numPayments; i++ {
		payID := uuid.New()
		clientID := clientIDs[rand.Intn(len(clientIDs))]
		quantity := rand.Intn(5) + 1
		unitPrice := gofakeit.Price(10, 500)
		totalAmount := float64(quantity) * unitPrice
		method := methods[rand.Intn(len(methods))]

		createdAt := time.Now().Add(-time.Duration(rand.Intn(365*24)) * time.Hour)

		paymentRows = append(paymentRows, []any{
			payID, clientID, quantity, unitPrice, totalAmount, method, createdAt,
		})

		var statusSequence []string
		flowChance := rand.Float64()

		switch {
		case flowChance < 0.70:
			statusSequence = []string{"created", "completed"}
		case flowChance < 0.80:
			statusSequence = []string{"created", "cancelled"}
		case flowChance < 0.90:
			statusSequence = []string{"created", "completed", "refunded"}
		case flowChance < 0.97:
			statusSequence = []string{"created", "failed"}
		default:
			statusSequence = []string{"created", "completed", "refunded", "failed"}
		}

		eventTime := createdAt
		for _, status := range statusSequence {
			eventRows = append(eventRows, []any{
				uuid.New(), payID, status, eventTime,
			})

			eventTime = eventTime.Add(time.Duration(rand.Intn(10)+1) * time.Second)
		}
	}

	_, err = pool.CopyFrom(ctx, pgx.Identifier{"payments"},
		[]string{"id", "client_id", "quantity", "unit_price", "total_amount", "payment_method", "created_at"},
		pgx.CopyFromRows(paymentRows))
	if err != nil {
		log.Fatalf("Error on insert payments: %v", err)
	}

	_, err = pool.CopyFrom(ctx, pgx.Identifier{"payment_events"},
		[]string{"id", "payment_id", "status", "occurred_at"},
		pgx.CopyFromRows(eventRows))
	if err != nil {
		log.Fatalf("Error on insert events: %v", err)
	}

	log.Printf("%d Payments and %d Events inserted in %v", numPayments, len(eventRows), time.Since(start))
	log.Println("Seed finished succesfully!")
}
