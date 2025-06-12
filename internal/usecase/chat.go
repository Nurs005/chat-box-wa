package usecase

import (
	"context"
	"github.com/chatbox/whatsapp/internal/domain"
	"github.com/chatbox/whatsapp/internal/repository/chat"
	"time"
)

type ChatService struct {
	repo chat.IChatRepository
}

type IChatService interface {
	SaveOrUpdate(ctx context.Context, chat *domain.Chat, message *domain.Message) error
}

func NewChatService(repo chat.IChatRepository) IChatService {
	return &ChatService{repo: repo}
}

func (s *ChatService) SaveOrUpdate(ctx context.Context, chat *domain.Chat, message *domain.Message) error {
	// Обновим чат (если он уже есть — обновим updated_at, если нет — создадим)
	chat.UpdatedAt = time.Now()
	if err := s.repo.SaveChat(ctx, chat); err != nil {
		return err
	}

	// Сохраняем сообщение
	return s.repo.SaveMessages(ctx, []*domain.Message{message})
}
