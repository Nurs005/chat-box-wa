package repository

import (
	"github.com/chatbox/whatsapp/internal/repository/chat"
	"github.com/chatbox/whatsapp/internal/repository/session"
	"gorm.io/gorm"
)

type WuzRepo struct {
	chat.IChatRepository
	session.ISessionRepository
}

func NewWuzRepo(db *gorm.DB) *WuzRepo {
	return &WuzRepo{
		chat.NewChatRepository(db),
		session.NewSessionRepository(db),
	}
}
