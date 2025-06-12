package session

import (
	"context"
	"github.com/chatbox/whatsapp/internal/domain"
)

func (r *SessionRepository) GetActiveSessions(ctx context.Context) ([]domain.Session, error) {
	var models []domain.Session
	if err := r.db.WithContext(ctx).Where("active = ?", true).Find(&models).Error; err != nil {
		return nil, err
	}

	sessions := make([]domain.Session, len(models))
	for i, m := range models {
		sessions[i] = m
	}
	return sessions, nil
}

func (r *SessionRepository) Save(ctx context.Context, model *domain.Session) error {
	return r.db.WithContext(ctx).Save(&model).Error
}

func (r *SessionRepository) GetByToken(ctx context.Context, token string) (model *domain.Session, err error) {
	err = r.db.WithContext(ctx).Where("token = ?", token).Find(model).Error
	if err != nil {
		return nil, err
	}
	return model, err
}
