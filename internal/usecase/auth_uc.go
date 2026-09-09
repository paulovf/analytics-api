package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/paulovf/analytics-api/internal/domain"
)

type AuthUseCase struct {
	repo      domain.ApiClientRepository
	jwtSecret string
}

func NewAuthUseCase(repo domain.ApiClientRepository, jwtSecret string) *AuthUseCase {
	return &AuthUseCase{repo: repo, jwtSecret: jwtSecret}
}

func (uc *AuthUseCase) Authenticate(ctx context.Context, clientID, clientSecret string) (string, error) {
	client, err := uc.repo.FindByClientID(ctx, clientID)
	if err != nil {
		return "", errors.New("Invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(client.ClientSecretHash), []byte(clientSecret)); err != nil {
		return "", errors.New("Invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":       client.ID.String(),
		"client_id": client.ClientID,
		"type":      client.Type,
		"exp":       time.Now().Add(time.Hour * 1).Unix(),
		"iat":       time.Now().Unix(),
	})

	return token.SignedString([]byte(uc.jwtSecret))
}
