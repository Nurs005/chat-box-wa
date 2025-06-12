package usecase

import (
	"context"
	"github.com/chatbox/whatsapp/internal/domain"
	"github.com/chatbox/whatsapp/internal/infrastructure"
	"github.com/chatbox/whatsapp/internal/repository/session"
	"github.com/chatbox/whatsapp/pkg/errors"
	"github.com/google/uuid"
)

type SessionService struct {
	repo      session.ISessionRepository
	waClient  infrastructure.WhatsAppClient
	logger    *infrastructure.Logger
	eventServ IEventsService
}

type ISessionService interface {
	ReconnectActiveSessions(ctx context.Context) error
	SendMessage(ctx context.Context, session *domain.Session, jid string, text string) error
	NewSession(ctx context.Context, session *domain.Session) error
	GetQRLogin(ctx context.Context, token string) (string, error)
	GetSession(ctx context.Context, token string) (*domain.Session, error)
}

func NewSessionService(repo session.ISessionRepository, wa infrastructure.WhatsAppClient, log *infrastructure.Logger, eventSrv IEventsService) ISessionService {
	return &SessionService{
		repo:      repo,
		waClient:  wa,
		logger:    log,
		eventServ: eventSrv,
	}
}

func (s *SessionService) NewSession(ctx context.Context, session *domain.Session) error {
	if session.BusinessID == 0 {
		return errors.HaveNoBussinessId
	}
	token, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	session.Token = token.String()
	if err := s.repo.Save(ctx, session); err != nil {
		return err
	}
	return s.waClient.ConnectWithHandler(session, func(i interface{}) {
		s.eventServ.HandleEvent(session, i)
	})
}

func (s *SessionService) GetSession(ctx context.Context, token string) (*domain.Session, error) {
	return s.repo.GetByToken(ctx, token)
}

func (s *SessionService) GetQRLogin(ctx context.Context, token string) (string, error) {
	return s.waClient.GenerateQR(ctx, &domain.Session{Token: token})
}

func (s *SessionService) ReconnectActiveSessions(ctx context.Context) error {
	sessions, err := s.repo.GetActiveSessions(ctx)
	if err != nil {
		return err
	}

	for _, session := range sessions {
		s.logger.Info().Str("token", session.Token).Msg("Reconnecting session")
		err := s.waClient.ConnectWithHandler(&session, func(i interface{}) {
			s.eventServ.HandleEvent(&session, i)
		})
		if err != nil {
			s.logger.Error().Err(err).Int("id", session.ID).Msg("Failed to reconnect session")
			continue
		}
	}

	return nil
}

func (s *SessionService) SendMessage(ctx context.Context, session *domain.Session, jid string, text string) error {
	// 1. Отправка через WhatsApp клиента
	if err := s.waClient.Send(ctx, session, jid, text); err != nil {
		s.logger.Error().Err(err).Int("id", session.ID).Msg("Failed to send message")
		return err
	}
	return nil
}
