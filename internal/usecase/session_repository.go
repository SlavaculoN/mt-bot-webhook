package usecase

import (
	"context"
	"errors"
	"mt-bot/internal/domain"
)

type SessionRepository interface {
	Get(ctx context.Context, chaID int64) (domain.UserSession, error)
	Save(ctx context.Context, chatID int64, session domain.UserSession) error
	Delete(ctx context.Context, chatID int64) error
}

var (
	ErrSessionNotFound     = errors.New("session not found")
	ErrInvalidDirection    = errors.New("invalid direction")
	ErrInvalidStep         = errors.New("invalid app step")
	ErrAddressNotAvailable = errors.New("address is not available for select direction")
	ErrInvalidName         = errors.New("name cannot be empty")
	ErrInvalidPhone        = errors.New("phone number connot be empty")
	ErrInvalidCar          = errors.New("car cannot be empty")
	ErrInvalidProblem      = errors.New("problem description cannot be empty")
)
