package telegram

import (
	"fmt"
	"mt-bot/internal/domain"
)

func formatAdminApplication(app domain.Application) string {
	return fmt.Sprintf(
		"🚘 НОВАЯ ЗАЯВКА\n\n"+
			"Услуга: %s\n"+
			"Адрес: %s\n\n"+
			"Имя: %s\n"+
			"Телефон: %s\n\n"+
			"Автомобиль: %s\n"+
			"Проблема: %s",
		formatDirection(app.Direction),
		formatAddress(app.Address),
		app.Name,
		app.Phone,
		app.Car,
		app.Problem,
	)
}

func formatDirection(direction domain.Direction) string {
	switch direction {
	case domain.DirectionService:
		return "Сервис"
	case domain.DirectionDetailing:
		return "Детейлинг"
	case domain.DirectionCarWash:
		return "Мойка"
	case domain.DirectionBodyWork:
		return "Кузовные работы"
	default:
		return "Неизвестная услуга"
	}
}

func formatAddress(address domain.Address) string {
	switch address {
	case domain.AddressPoeticheskiy:
		return "Поэтический д.4"
	case domain.AddressSikeirosa:
		return "Сикейроса д.14"
	default:
		return "Неизвестный адрес"
	}
}
