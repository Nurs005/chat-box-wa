package usecase

import (
	"context"
	"encoding/json"
	"github.com/chatbox/whatsapp/internal/domain"
	"github.com/chatbox/whatsapp/internal/infrastructure"
	"go.mau.fi/whatsmeow/types/events"
)

type EventsService struct {
	webHub  IHub
	logger  *infrastructure.Logger
	chatSrv IChatService
}

type IEventsService interface {
	HandleEvent(session *domain.Session, evt interface{})
}

func NewEventsService(webHub IHub, logger *infrastructure.Logger, chatSrv IChatService) IEventsService {
	return &EventsService{
		webHub:  webHub,
		logger:  logger,
		chatSrv: chatSrv,
	}
}

func (e *EventsService) HandleEvent(session *domain.Session, evt interface{}) {
	switch ev := evt.(type) {
	case *events.Message:
		if ev.Message == nil {
			return
		}

		// 1. Подготовка данных
		chatJID := ev.Info.Chat.String()
		timestamp := ev.Info.Timestamp

		chat := &domain.Chat{
			SessionToken: session.Token,
			JID:          chatJID,
			Title:        ev.Info.Chat.User,
			UpdatedAt:    timestamp,
		}

		message := &domain.Message{
			SessionToken: session.Token,
			ChatJID:      chatJID,
			FromMe:       ev.Info.IsFromMe,
			Text:         ev.Message.GetConversation(),
			Timestamp:    timestamp,
		}

		// 2. Сохраняем в БД
		if err := e.chatSrv.SaveOrUpdate(context.Background(), chat, message); err != nil {
			e.logger.Error().Err(err).Msg("failed to save chat/message")
			return
		}

		// 3. Отправка через WebSocket
		ws := domain.WSMessageDTO{
			Type:         "message",
			ChatJID:      chat.JID,
			From:         ev.Info.SourceString(),
			Text:         message.Text,
			FromMe:       message.FromMe,
			Timestamp:    message.Timestamp,
			MessageID:    message.ID,
			SessionToken: session.Token,
		}

		jsonBytes, _ := json.Marshal(ws)
		if err := e.webHub.Send(session.Token, jsonBytes); err != nil {
			e.logger.Error().Err(err).Msg("Failed to send WS message")
		}
	}
}
