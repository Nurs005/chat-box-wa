package session

import (
	"context"
	"errors"

	"github.com/chatbox/whatsapp/internal/domain"
	"gorm.io/gorm"
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
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *SessionRepository) GetByToken(ctx context.Context, token string) (model *domain.Session, err error) {
	model = &domain.Session{}
	err = r.db.WithContext(ctx).Where("token = ?", token).First(model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
  
	return model, nil
}
