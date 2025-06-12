package internal

import (
	"github.com/chatbox/whatsapp/internal/infrastructure"
	"github.com/chatbox/whatsapp/internal/repository"
	"github.com/chatbox/whatsapp/internal/usecase"
)

type WuzApi struct {
	usecase.ISessionService
	usecase.IHub
	usecase.IEventsService
	usecase.IWebSocketCommandResolver
	usecase.IChatService
}

func NewWuzApi(
	repo *repository.WuzRepo,
	client infrastructure.WhatsAppClient,
	logger *infrastructure.Logger,
) *WuzApi {
	chatSrv := usecase.NewChatService(repo.IChatRepository)
	hub := usecase.NewHub()
	eventSrv := usecase.NewEventsService(hub, logger, chatSrv)
	sessionSrv := usecase.NewSessionService(repo.ISessionRepository, client, logger, eventSrv)
	return &WuzApi{
		ISessionService: sessionSrv,
		IHub:            hub,
		IEventsService:  eventSrv,
		IChatService:    chatSrv,
	}
}
