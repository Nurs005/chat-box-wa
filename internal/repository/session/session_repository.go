package session

import (
	"context"
	"github.com/chatbox/whatsapp/internal/domain"
	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) ISessionRepository {
	return &SessionRepository{db: db}
}

type ISessionRepository interface {
	GetActiveSessions(ctx context.Context) ([]domain.Session, error)
	Save(ctx context.Context, session *domain.Session) error
	GetByToken(ctx context.Context, token string) (model *domain.Session, err error)
}
