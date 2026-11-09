package telegram

import (
	"context"
	"mt-bot/internal/domain"
	"mt-bot/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramApplicationSender struct {
	bot *tgbotapi.BotAPI
}

func NewTelegramApplicationSender(bot *tgbotapi.BotAPI) *TelegramApplicationSender {
	return &TelegramApplicationSender{
		bot: bot,
	}
}

func (s *TelegramApplicationSender) SendToAdmin(ctx context.Context, adminChatID int64, application domain.Application) error {

	if err := ctx.Err(); err != nil {
		return err
	}

	text := formatAdminApplication(application)

	message := tgbotapi.NewMessage(adminChatID, text)

	if _, err := s.bot.Send(message); err != nil {
		return err
	}

	return nil
}

var _ usecase.ApplicationSender = (*TelegramApplicationSender)(nil)
