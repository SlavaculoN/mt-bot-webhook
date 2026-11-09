package telegram

import (
	"context"
	"fmt"
	"log"
	"mt-bot/internal/domain"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ApplicationService interface {
	Start(ctx context.Context, chatID int64) (domain.UserSession, error)

	SelectDirection(
		ctx context.Context,
		chatID int64,
		direction domain.Direction,
	) (domain.UserSession, error)

	SelectAddress(
		ctx context.Context,
		chatID int64,
		address domain.Address,
	) (domain.UserSession, error)

	SetName(
		ctx context.Context,
		chatID int64,
		name string,
	) (domain.UserSession, error)

	SetPhone(
		ctx context.Context,
		chatID int64,
		phone string,
	) (domain.UserSession, error)

	SetCar(
		ctx context.Context,
		chatID int64,
		car string,
	) (domain.UserSession, error)

	SetProblem(
		ctx context.Context,
		chatID int64,
		problem string,
	) (domain.UserSession, error)

	CancelApplication(ctx context.Context, chatID int64) error

	SubmitApplication(ctx context.Context, chatID int64) error

	GetSession(ctx context.Context, chatID int64) (domain.UserSession, error)
}

type TelegramHandler struct {
	app        ApplicationService
	bot        *tgbotapi.BotAPI
	messageIDs MessageIDRepository
}

func NewTelegramHandler(app ApplicationService, bot *tgbotapi.BotAPI, messageIDs MessageIDRepository) *TelegramHandler {
	return &TelegramHandler{
		app:        app,
		bot:        bot,
		messageIDs: messageIDs,
	}
}

var _ MessageIDRepository = (*MemoryMessageIDRepository)(nil)

func (h *TelegramHandler) HandleUpdate(ctx context.Context, update tgbotapi.Update) error {
	if update.Message != nil {
		return h.handleMessage(ctx, update.Message)
	}

	if update.CallbackQuery != nil {
		return h.handleCallback(ctx, update.CallbackQuery)
	}

	return nil
}

func (h *TelegramHandler) handleMessage(ctx context.Context, message *tgbotapi.Message) error {
	chatID := message.Chat.ID

	if message.IsCommand() {
		switch message.Command() {
		case "start":
			_, err := h.app.Start(ctx, chatID)
			if err != nil {
				return err
			}

			msg := tgbotapi.NewMessage(chatID, "Добро пожаловать!\n\nВыберите направление обслуживания:")

			msg.ReplyMarkup = directionKeyboard()

			sentMessage, err := h.bot.Send(msg)
			if err != nil {
				return err
			}

			if err := h.messageIDs.Save(chatID, sentMessage.MessageID); err != nil {
				return err
			}
		default:
			return nil
		}
		return nil
	}

	return h.handleTextMessage(ctx, message)
}

func (h *TelegramHandler) handleCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) error {
	if err := h.answerCallback(callback.ID); err != nil {
		return err
	}

	if callback.Message == nil {
		return nil
	}

	switch callback.Data {
	case callbackDirectionService:
		return h.handleDirectionSelection(ctx, callback, domain.DirectionService)

	case callbackDirectionDetailing:
		return h.handleDirectionSelection(ctx, callback, domain.DirectionDetailing)

	case callbackDirectionCarWash:
		return h.handleDirectionSelection(ctx, callback, domain.DirectionCarWash)

	case callbackDirectionBodyWork:
		return h.handleDirectionSelection(ctx, callback, domain.DirectionBodyWork)

	case callbackAddressPoeticheskiy:
		return h.handleAddressSelection(ctx, callback, domain.AddressPoeticheskiy)

	case callbackAddressSikeirosa:
		return h.handleAddressSelection(ctx, callback, domain.AddressSikeirosa)

	case callbackSubmitApplication:
		return h.handleSubmitApplication(ctx, callback)

	case callbackCancelApplication:
		return h.handleCancelApplication(ctx, callback)
	default:
		return nil
	}
}

func (h *TelegramHandler) answerCallback(callbackID string) error {
	config := tgbotapi.NewCallback(callbackID, "")

	_, err := h.bot.Request(config)

	return err
}

func (h *TelegramHandler) handleDirectionSelection(ctx context.Context, callback *tgbotapi.CallbackQuery, direction domain.Direction) error {
	chatID := callback.Message.Chat.ID

	messageID, err := h.messageIDs.Get(chatID)
	if err != nil {
		return err
	}

	if callback.Message.MessageID != messageID {
		return nil
	}

	updatedSession, err := h.app.SelectDirection(ctx, chatID, direction)
	if err != nil {
		return err
	}

	text := formatApplicationProgress(updatedSession)

	keyboard := addressKeyboard(updatedSession.Application.Direction)

	edit := tgbotapi.NewEditMessageText(chatID, messageID, text)

	edit.ReplyMarkup = &keyboard

	if _, err := h.bot.Send(edit); err != nil {
		return err
	}

	return nil
}

func (h *TelegramHandler) handleAddressSelection(ctx context.Context, callback *tgbotapi.CallbackQuery, address domain.Address) error {
	chatID := callback.Message.Chat.ID

	messageID, err := h.messageIDs.Get(chatID)
	if err != nil {
		return err
	}

	if callback.Message.MessageID != messageID {
		return nil
	}

	updatedSession, err := h.app.SelectAddress(ctx, chatID, address)
	if err != nil {
		return err
	}

	text := formatApplicationProgress(updatedSession)

	edit := tgbotapi.NewEditMessageText(chatID, messageID, text)

	edit.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{},
	}

	if _, err := h.bot.Send(edit); err != nil {
		return err
	}

	return nil
}

func (h *TelegramHandler) handleTextMessage(ctx context.Context, message *tgbotapi.Message) error {
	chatID := message.Chat.ID

	session, err := h.app.GetSession(ctx, chatID)
	if err != nil {
		return err
	}

	switch session.Step {
	case domain.StepName:
		messageID, err := h.messageIDs.Get(chatID)
		if err != nil {
			return err
		}

		updatedSession, err := h.app.SetName(ctx, chatID, message.Text)
		if err != nil {
			return err
		}

		text := formatApplicationProgress(updatedSession)

		edit := tgbotapi.NewEditMessageText(chatID, messageID, text)

		if _, err := h.bot.Send(edit); err != nil {
			return err
		}

		confirmation := tgbotapi.NewMessage(chatID, fmt.Sprintf("✅ Записал имя: %s", updatedSession.Application.Name))

		confirmationMessage, err := h.bot.Send(confirmation)
		if err != nil {
			return err
		}

		h.deleteMessagesAfter(chatID, 1*time.Second, confirmationMessage.MessageID, message.MessageID)
	case domain.StepPhone:
		messageID, err := h.messageIDs.Get(chatID)
		if err != nil {
			return err
		}

		updatedSession, err := h.app.SetPhone(ctx, chatID, message.Text)
		if err != nil {
			return err
		}

		text := formatApplicationProgress(updatedSession)

		edit := tgbotapi.NewEditMessageText(chatID, messageID, text)

		if _, err := h.bot.Send(edit); err != nil {
			return err
		}

		confirmation := tgbotapi.NewMessage(chatID, fmt.Sprintf("✅ Записал телефон: %s", updatedSession.Application.Phone))

		confirmationMessage, err := h.bot.Send(confirmation)
		if err != nil {
			return err
		}

		h.deleteMessagesAfter(chatID, 1*time.Second, confirmationMessage.MessageID, message.MessageID)
	case domain.StepCar:
		messageID, err := h.messageIDs.Get(chatID)
		if err != nil {
			return err
		}

		updatedSession, err := h.app.SetCar(ctx, chatID, message.Text)
		if err != nil {
			return err
		}

		text := formatApplicationProgress(updatedSession)

		edit := tgbotapi.NewEditMessageText(chatID, messageID, text)

		if _, err := h.bot.Send(edit); err != nil {
			return err
		}

		confirmation := tgbotapi.NewMessage(chatID, fmt.Sprintf("✅ Записал автомобиль: %s", updatedSession.Application.Car))

		confirmationMessage, err := h.bot.Send(confirmation)
		if err != nil {
			return err
		}

		h.deleteMessagesAfter(chatID, 1*time.Second, confirmationMessage.MessageID, message.MessageID)
	case domain.StepProblem:
		messageID, err := h.messageIDs.Get(chatID)
		if err != nil {
			return err
		}

		updatedSession, err := h.app.SetProblem(ctx, chatID, message.Text)
		if err != nil {
			return err
		}

		text := formatApplicationPreview(updatedSession)

		keyboard := confirmationKeyboard()

		edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
		edit.ReplyMarkup = &keyboard

		if _, err := h.bot.Send(edit); err != nil {
			return err
		}

		confirmation := tgbotapi.NewMessage(chatID, "✅ Описание проблемы записано!")

		confirmationMessage, err := h.bot.Send(confirmation)
		if err != nil {
			return err
		}

		h.deleteMessagesAfter(chatID, 1*time.Second, confirmationMessage.MessageID, message.MessageID)

	default:
		return nil
	}

	return nil
}

func (h *TelegramHandler) deleteTelegramMessage(chatID int64, messageID int) error {
	config := tgbotapi.NewDeleteMessage(chatID, messageID)

	_, err := h.bot.Request(config)
	if err != nil {
		return err
	}

	return nil
}

func (h *TelegramHandler) deleteMessagesAfter(chatID int64, delay time.Duration, messageIDs ...int) {
	time.AfterFunc(delay, func() {
		for _, messageID := range messageIDs {
			if err := h.deleteTelegramMessage(chatID, messageID); err != nil {
				log.Printf("failed to delete message %d in chat %d: %v", messageID, chatID, err)
			}
		}
	})
}

func (h *TelegramHandler) handleSubmitApplication(ctx context.Context, callback *tgbotapi.CallbackQuery) error {
	chatID := callback.Message.Chat.ID

	messageID, err := h.messageIDs.Get(chatID)
	if err != nil {
		return err
	}

	if callback.Message.MessageID != messageID {
		return nil
	}

	if err := h.app.SubmitApplication(ctx, chatID); err != nil {
		return err
	}

	text := "Спасибо за заявку!\n\n" +
		"Ваша заявка успешно отправлена администратору.\n\n" +
		"В ближайшее время администратор свяжется с Вами.\n\n" +
		"Чтобы оставить новую заявку, отправьте /start"

	edit := tgbotapi.NewEditMessageText(chatID, messageID, text)

	edit.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{},
	}

	if _, err := h.bot.Send(edit); err != nil {
		return err
	}

	if err := h.messageIDs.Delete(chatID); err != nil {
		return err
	}

	return nil
}

func (h *TelegramHandler) handleCancelApplication(ctx context.Context, callback *tgbotapi.CallbackQuery) error {
	chatID := callback.Message.Chat.ID

	messageID, err := h.messageIDs.Get(chatID)
	if err != nil {
		return err
	}

	if callback.Message.MessageID != messageID {
		return nil
	}

	if err := h.app.CancelApplication(ctx, chatID); err != nil {
		return err
	}

	text := "❌ Заявка отменена.\n\nВы можете начать оформление заново, отправив /start."

	edit := tgbotapi.NewEditMessageText(chatID, messageID, text)

	edit.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{},
	}

	if _, err := h.bot.Send(edit); err != nil {
		return err
	}

	if err := h.messageIDs.Delete(chatID); err != nil {
		return err
	}

	return nil
}
