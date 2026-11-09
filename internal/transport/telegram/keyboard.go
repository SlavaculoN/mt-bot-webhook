package telegram

import (
	"mt-bot/internal/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	callbackSubmitApplication   = "application:submit"
	callbackCancelApplication   = "application:cancel"
	callbackDirectionService    = "direction:service"
	callbackDirectionDetailing  = "direction:detailing"
	callbackDirectionCarWash    = "direction:carwash"
	callbackDirectionBodyWork   = "direction:bodywork"
	callbackAddressPoeticheskiy = "address:poeticheskiy"
	callbackAddressSikeirosa    = "address:sikeirosa"
)

func confirmationKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Отправить", callbackSubmitApplication),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❌ Отменить", callbackCancelApplication),
		),
	)
}

func directionKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🔧 Сервис", callbackDirectionService,
			),
			tgbotapi.NewInlineKeyboardButtonData(
				"✨ Детейлинг", callbackDirectionDetailing,
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🚿 Автомойка", callbackDirectionCarWash,
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🔨 Кузовные работы", callbackDirectionBodyWork,
			),
		),
	)
}

func addressKeyboard(direction domain.Direction) tgbotapi.InlineKeyboardMarkup {
	var buttons []tgbotapi.InlineKeyboardButton

	if domain.IsAddressAvailable(direction, domain.AddressPoeticheskiy) {

		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(
			"📍 Поэтический",
			callbackAddressPoeticheskiy,
		))

	}

	if domain.IsAddressAvailable(direction, domain.AddressSikeirosa) {

		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(
			"📍 Сикейроса",
			callbackAddressSikeirosa,
		))

	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(buttons...),
	)
}
