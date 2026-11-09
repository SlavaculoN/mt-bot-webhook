package usecase

import (
	"context"
	"mt-bot/internal/domain"
)

type ApplicationSender interface {
	SendToAdmin(ctx context.Context, adminChatID int64, application domain.Application) error
}
