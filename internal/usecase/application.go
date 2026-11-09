package usecase

import (
	"context"
	"mt-bot/internal/domain"
	"strings"
)

type ApplicationUseCase struct {
	sessions SessionRepository
	sender   ApplicationSender
	router   AdminRouter
}

func NewApplicationUseCase(sessions SessionRepository, sender ApplicationSender, router AdminRouter) *ApplicationUseCase {
	return &ApplicationUseCase{
		sessions: sessions,
		sender:   sender,
		router:   router,
	}
}

func (u *ApplicationUseCase) Start(ctx context.Context, chatID int64) (domain.UserSession, error) {
	if err := ctx.Err(); err != nil {
		return domain.UserSession{}, err
	}

	session := domain.UserSession{
		Step:        domain.StepDirection,
		Application: domain.Application{},
	}

	if err := u.sessions.Save(ctx, chatID, session); err != nil {
		return domain.UserSession{}, err
	}

	return session, nil
}

func (u *ApplicationUseCase) SelectDirection(ctx context.Context, chatID int64, direction domain.Direction) (domain.UserSession, error) {

	session, err := u.sessions.Get(ctx, chatID)
	if err != nil {
		return domain.UserSession{}, err
	}

	if session.Step != domain.StepDirection {
		return domain.UserSession{}, ErrInvalidStep
	}

	switch direction {
	case domain.DirectionService,
		domain.DirectionDetailing,
		domain.DirectionCarWash,
		domain.DirectionBodyWork:
	default:
		return domain.UserSession{}, ErrInvalidDirection
	}

	session.Application.Direction = direction
	session.Step = domain.StepAddress

	if err := u.sessions.Save(ctx, chatID, session); err != nil {
		return domain.UserSession{}, err
	}

	return session, nil
}

func (u *ApplicationUseCase) SelectAddress(ctx context.Context, chatID int64, address domain.Address) (domain.UserSession, error) {
	session, err := u.sessions.Get(ctx, chatID)
	if err != nil {
		return domain.UserSession{}, err
	}

	if session.Step != domain.StepAddress {
		return domain.UserSession{}, ErrInvalidStep
	}

	ok := domain.IsAddressAvailable(session.Application.Direction, address)
	if !ok {
		return domain.UserSession{}, ErrAddressNotAvailable
	}

	session.Application.Address = address
	session.Step = domain.StepName

	if err := u.sessions.Save(ctx, chatID, session); err != nil {
		return domain.UserSession{}, err
	}

	return session, nil
}

func (u *ApplicationUseCase) SetName(ctx context.Context, chatID int64, name string) (domain.UserSession, error) {
	session, err := u.sessions.Get(ctx, chatID)
	if err != nil {
		return domain.UserSession{}, err
	}

	if session.Step != domain.StepName {
		return domain.UserSession{}, ErrInvalidStep
	}

	name = strings.TrimSpace(name)

	if name == "" {
		return domain.UserSession{}, ErrInvalidName
	}

	session.Application.Name = name
	session.Step = domain.StepPhone

	if err := u.sessions.Save(ctx, chatID, session); err != nil {
		return domain.UserSession{}, err
	}

	return session, nil
}

func (u *ApplicationUseCase) SetPhone(ctx context.Context, chatID int64, phone string) (domain.UserSession, error) {
	session, err := u.sessions.Get(ctx, chatID)
	if err != nil {
		return domain.UserSession{}, err
	}

	if session.Step != domain.StepPhone {
		return domain.UserSession{}, ErrInvalidStep
	}

	phone = strings.TrimSpace(phone)

	if phone == "" {
		return domain.UserSession{}, ErrInvalidPhone
	}

	session.Application.Phone = phone
	session.Step = domain.StepCar

	if err := u.sessions.Save(ctx, chatID, session); err != nil {
		return domain.UserSession{}, err
	}

	return session, nil
}

func (u *ApplicationUseCase) SetCar(ctx context.Context, chatID int64, car string) (domain.UserSession, error) {
	session, err := u.sessions.Get(ctx, chatID)
	if err != nil {
		return domain.UserSession{}, err
	}

	if session.Step != domain.StepCar {
		return domain.UserSession{}, ErrInvalidStep
	}

	carValid := strings.TrimSpace(car)

	if carValid == "" {
		return domain.UserSession{}, ErrInvalidCar
	}

	session.Application.Car = car
	session.Step = domain.StepProblem

	if err := u.sessions.Save(ctx, chatID, session); err != nil {
		return domain.UserSession{}, err
	}

	return session, nil
}

func (u *ApplicationUseCase) SetProblem(ctx context.Context, chatID int64, problem string) (domain.UserSession, error) {
	session, err := u.sessions.Get(ctx, chatID)
	if err != nil {
		return domain.UserSession{}, err
	}

	if session.Step != domain.StepProblem {
		return domain.UserSession{}, ErrInvalidStep
	}

	problemValid := strings.TrimSpace(problem)

	if problemValid == "" {
		return domain.UserSession{}, ErrInvalidProblem
	}

	session.Application.Problem = problem
	session.Step = domain.StepConfirm

	if err := u.sessions.Save(ctx, chatID, session); err != nil {
		return domain.UserSession{}, err
	}

	return session, nil
}

func (u *ApplicationUseCase) CancelApplication(ctx context.Context, chatID int64) error {
	session, err := u.sessions.Get(ctx, chatID)
	if err != nil {
		return err
	}

	if session.Step != domain.StepConfirm {
		return ErrInvalidStep
	}

	if err := u.sessions.Delete(ctx, chatID); err != nil {
		return err
	}

	return nil
}

func (u *ApplicationUseCase) SubmitApplication(ctx context.Context, chatID int64) error {
	session, err := u.sessions.Get(ctx, chatID)
	if err != nil {
		return err
	}

	if session.Step != domain.StepConfirm {
		return ErrInvalidStep
	}

	adminChatID, err := u.router.SelectAdmin(session.Application)
	if err != nil {
		return err
	}

	if err := u.sender.SendToAdmin(ctx, adminChatID, session.Application); err != nil {
		return err
	}

	if err := u.sessions.Delete(ctx, chatID); err != nil {
		return err
	}

	return nil
}

func (u *ApplicationUseCase) GetSession(ctx context.Context, chatID int64) (domain.UserSession, error) {
	return u.sessions.Get(ctx, chatID)
}
