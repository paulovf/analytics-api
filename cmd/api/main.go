package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/paulovf/analytics-api/internal/domain"
	httpHandler "github.com/paulovf/analytics-api/internal/handler/http"
	repoPostgres "github.com/paulovf/analytics-api/internal/repository/postgres"
	"github.com/paulovf/analytics-api/internal/usecase"
	"github.com/paulovf/analytics-api/pkg/postgres"
)

const jwtSecretKey = "analytics_super_secret_jwt_key_2026"

func main() {
	ctx := context.Background()
	dsn := "postgres://analytics_user:analytics_password@localhost:5432/analytics_db?sslmode=disable"

	pool, err := postgres.NewPool(ctx, dsn)
	if err != nil {
		log.Fatalf("Fail on start database: %v", err)
	}
	defer pool.Close()

	clientRepo := repoPostgres.NewClientRepo(pool)
	paymentRepo := repoPostgres.NewPaymentRepo(pool)
	apiClientRepo := repoPostgres.NewApiClientRepo(pool)
	analyticsRepo := repoPostgres.NewAnalyticsRepo(pool)

	// seedTestApiClient(ctx, apiClientRepo)

	clientUC := usecase.NewClientUseCase(clientRepo)
	paymentUC := usecase.NewPaymentUseCase(paymentRepo)
	authUC := usecase.NewAuthUseCase(apiClientRepo, jwtSecretKey)

	clientHandler := httpHandler.NewClientHandler(clientUC)
	paymentHandler := httpHandler.NewPaymentHandler(paymentUC)
	authHandler := httpHandler.NewAuthHandler(authUC)
	analyticsHandler := httpHandler.NewAnalyticsHandler(analyticsRepo)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/token", authHandler.Login)

		r.Group(func(r chi.Router) {
			r.Use(httpHandler.JWTMiddleware(jwtSecretKey))
		
			r.Get("/clients", clientHandler.List)
			r.Get("/clients/{id}", clientHandler.GetByID)

			r.Get("/payments", paymentHandler.List)
			r.Get("/payments/{id}", paymentHandler.GetByID)
			r.Get("/payments/{id}/events", paymentHandler.GetEvents)

			r.Get("/analytics/daily-stats", analyticsHandler.GetDailyStats)
		})
	})

	log.Println("HTTP Server started in port :8080...")
	http.ListenAndServe(":8080", r)
}

func seedTestApiClient(ctx context.Context, repo domain.ApiClientRepository) {
	_, err := repo.FindByClientID(ctx, "external_app_test")
	if err == nil {
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	client := &domain.ApiClient{
		ID:               uuid.New(),
		Name:             "External Integration Test",
		ClientID:         "external_app_test",
		ClientSecretHash: string(hash),
		Type:             "external",
		Status:           "active",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Time{},
	}

	if err := repo.Create(ctx, client); err != nil {
		log.Printf("Info: error on insert ApiClient seed: %v", err)
	} else {
		log.Println("Test apiClient created. [client_id: external_app_test, client_secret: secret123]")
	}
}
