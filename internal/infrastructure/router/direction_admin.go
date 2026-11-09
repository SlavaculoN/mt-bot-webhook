package router

import (
	"errors"
	"mt-bot/internal/domain"
	"mt-bot/internal/usecase"
)

type DirectionAdminRouter struct {
	generalAdminChatID        int64
	servicePoeticheskiyChatID int64
	serviceSikeirosaChatID    int64
}

var (
	ErrAdminNotFound      = errors.New("admin not found")
	ErrAdminNotConfigured = errors.New("admin chatID not loaded")
)

func NewDirectionAdminRouter(
	generalAdminChatID int64,
	servicePoeticheskiyChatID int64,
	serviceSikeirosaChatID int64,
) *DirectionAdminRouter {
	return &DirectionAdminRouter{
		generalAdminChatID:        generalAdminChatID,
		servicePoeticheskiyChatID: servicePoeticheskiyChatID,
		serviceSikeirosaChatID:    serviceSikeirosaChatID,
	}
}

func (r *DirectionAdminRouter) SelectAdmin(app domain.Application) (int64, error) {

	var adminChatID int64

	switch app.Direction {
	case domain.DirectionDetailing, domain.DirectionCarWash, domain.DirectionBodyWork:
		adminChatID = r.generalAdminChatID
	case domain.DirectionService:
		switch app.Address {
		case domain.AddressPoeticheskiy:
			adminChatID = r.servicePoeticheskiyChatID
		case domain.AddressSikeirosa:
			adminChatID = r.serviceSikeirosaChatID
		default:
			return 0, ErrAdminNotFound
		}
	default:
		return 0, ErrAdminNotFound
	}

	if adminChatID == 0 {
		return 0, ErrAdminNotConfigured
	}

	return adminChatID, nil
}

var _ usecase.AdminRouter = (*DirectionAdminRouter)(nil)
