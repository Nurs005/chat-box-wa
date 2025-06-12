package domain

import "time"

// Chat represents a chat conversation.
type Chat struct {
	ID           int       `gorm:"primaryKey"`
	SessionToken string    `gorm:"size:255;index"`
	JID          string    `gorm:"size:255;index"`
	Title        string    `gorm:"size:255"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

// ChatDTO is a transport representation of Chat.
type ChatDTO struct {
	ID           int    `json:"id"`
	SessionToken string `json:"session_token"`
	JID          string `json:"jid"`
	Title        string `json:"title"`
	UnreadCount  int    `json:"unread_count"`
}

func (*Chat) TableName() string {
	return "chats"
}

// ToDTO converts a Chat to its DTO form.
func (c *Chat) ToDTO(count int) ChatDTO {
	return ChatDTO{
		ID:           c.ID,
		SessionToken: c.SessionToken,
		JID:          c.JID,
		Title:        c.Title,
		UnreadCount:  count,
	}
}
