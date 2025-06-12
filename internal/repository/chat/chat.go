package chat

import (
	"context"
	"github.com/chatbox/whatsapp/internal/domain"
	"gorm.io/gorm/clause"
)

func (r *GormChatRepository) SaveChat(ctx context.Context, chat *domain.Chat) error {
	return r.db.WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "jid"}},
			DoUpdates: clause.AssignmentColumns([]string{"title", "updated_at"}),
		},
	).Create(chat).Error
}

func (r *GormChatRepository) SaveMessages(ctx context.Context, messages []*domain.Message) error {
	return r.db.WithContext(ctx).Create(&messages).Error
}

func (r *GormChatRepository) GetChats(ctx context.Context, sessionToken string) ([]*domain.Chat, error) {
	var chats []*domain.Chat
	err := r.db.WithContext(ctx).
		Where("session_token = ?", sessionToken).
		Order("updated_at DESC").
		Find(&chats).Error
	return chats, err
}

func (r *GormChatRepository) GetMessages(ctx context.Context, chatJID string) ([]*domain.Message, error) {
	var messages []*domain.Message
	err := r.db.WithContext(ctx).
		Where("chat_jid = ?", chatJID).
		Order("timestamp ASC").
		Find(&messages).Error
	return messages, err
}
