package usecase

import (
	"mt-bot/internal/domain"
)

type AdminRouter interface {
	SelectAdmin(app domain.Application) (int64, error)
}
