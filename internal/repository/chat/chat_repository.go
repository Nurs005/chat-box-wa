package chat

import (
	"context"
	"github.com/chatbox/whatsapp/internal/domain"
	"gorm.io/gorm"
)

type GormChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) *GormChatRepository {
	db.AutoMigrate(&domain.Chat{}, &domain.Message{})
	return &GormChatRepository{db: db}
}

type IChatRepository interface {
	SaveChat(ctx context.Context, chat *domain.Chat) error
	SaveMessages(ctx context.Context, messages []*domain.Message) error
	GetChats(ctx context.Context, sessionToken string) ([]*domain.Chat, error)
	GetMessages(ctx context.Context, chatJID string) ([]*domain.Message, error)
}
