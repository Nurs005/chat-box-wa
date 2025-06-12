package domain

import "time"

type Chat struct {
	ID           int       `gorm:"primaryKey"`
	SessionToken string    `gorm:"size:255;index"`
	JID          string    `gorm:"size:255;index"`
	Title        string    `gorm:"size:255"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}


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

func (c *Chat) ToDTO(count int) (chatDto ChatDTO) {
	return ChatDTO{
		ID:           c.ID,
		SessionToken: c.SessionToken,
		JID:          c.JID,
		Title:        c.Title,
		UnreadCount:  count,
	}
}
