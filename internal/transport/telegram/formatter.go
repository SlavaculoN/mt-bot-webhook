package telegram

import (
	"fmt"
	"mt-bot/internal/domain"
	"strings"
)

func formatApplicationPreview(session domain.UserSession) string {
	return fmt.Sprintf(
		"🚘 ВАША ЗАЯВКА\n\n"+
			"Услуга: %s\n"+
			"Адрес: %s\n"+
			"Имя: %s\n"+
			"Телефон: %s\n"+
			"Автомобиль: %s\n"+
			"Проблема: %s\n\n"+
			"Проверьте данные перед отправкой.",
		formatDirection(session.Application.Direction),
		formatAddress(session.Application.Address),
		session.Application.Name,
		session.Application.Phone,
		session.Application.Car,
		session.Application.Problem,
	)
}

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

func formatApplicationProgress(session domain.UserSession) string {
	app := session.Application

	direction := "-"
	if app.Direction != "" {
		direction = formatDirection(app.Direction)
	}

	address := "-"
	if app.Address != "" {
		address = formatAddress(app.Address)
	}

	nextStep := ""

	switch session.Step {
	case domain.StepDirection:
		nextStep = "Выберете направление обслуживания:"
	case domain.StepAddress:
		nextStep = "Выберете адрес обслуживания:"
	case domain.StepName:
		nextStep = "Введите ваше имя:"
	case domain.StepPhone:
		nextStep = "Введите номер телефона:"
	case domain.StepCar:
		nextStep = "Введите марку и модель автомобиля:"
	case domain.StepProblem:
		nextStep = "Опишите проблему, с которой обратились:"
	case domain.StepConfirm:
		nextStep = "Проверьте данные перед отправкой."
	}

	return fmt.Sprintf(
		"🚘 ВАША ЗАЯВКА\n\n"+
			"Услуга: %s\n"+
			"Адрес: %s\n"+
			"Имя: %s\n"+
			"Телефон: %s\n"+
			"Автомобиль: %s\n"+
			"Проблема: %s\n\n"+
			"%s",
		direction,
		address,
		valueOrDash(app.Name),
		valueOrDash(app.Phone),
		valueOrDash(app.Car),
		valueOrDash(app.Problem),
		nextStep,
	)
}

func valueOrDash(value string) string {
	if strings.TrimSpace(value) == "" {
		return "-"
	}

	return value
}
