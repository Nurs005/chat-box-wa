package usecase

import (
	"context"
	"encoding/json"
	"github.com/chatbox/whatsapp/internal/domain"
	"github.com/chatbox/whatsapp/pkg/errors"

	"github.com/chatbox/whatsapp/internal/repository/session"
	"github.com/rs/zerolog"
)

type WebSocketCommandHandler struct {
	sessionRepo    session.ISessionRepository
	sessionService ISessionService
	logger         zerolog.Logger
}

func NewWebSocketCommandHandler(
	sessionRepo session.ISessionRepository,
	sessionService ISessionService,
	logger zerolog.Logger,
) IWebSocketCommandResolver {
	return &WebSocketCommandHandler{
		sessionRepo:    sessionRepo,
		sessionService: sessionService,
		logger:         logger,
	}
}

type IWebSocketCommandResolver interface {
	Handle(session *domain.Session, raw []byte) error
}

func (h *WebSocketCommandHandler) Handle(session *domain.Session, raw []byte) error {
	var cmd domain.WebSocketCommand
	if err := json.Unmarshal(raw, &cmd); err != nil {
		h.logger.Warn().Err(err).Msg("Invalid WS command format")
		return errors.InvalidBodyForCmd
	}

	switch cmd.Type {
	case "send":
		if err := h.sessionService.SendMessage(context.Background(), session, cmd.ChatJID, cmd.Text); err != nil {
			h.logger.Error().Err(err).Msg("Failed to send message")
			return err
		}
	default:
		h.logger.Warn().Str("type", cmd.Type).Msg("Unknown command")
		return errors.InvalidCommand
	}
	return nil
}
