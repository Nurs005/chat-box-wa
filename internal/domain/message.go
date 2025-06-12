package domain

import "time"

// Message represents a message belonging to a chat.
type Message struct {
	ID           int       `gorm:"primaryKey"`
	SessionToken string    `gorm:"size:255;index"`
	ChatJID      string    `gorm:"size:255;index"`
	IsRead       bool      `gorm:"default:false"`
	IsFromMe     bool      `gorm:"default:false"`
	Text         string    `gorm:"type:text"`
	Timestamp    time.Time `gorm:"index"`
}

func (Message) TableName() string { return "messages" }

// WSMessageDTO is a transport representation of Message for websocket clients.
type WSMessageDTO struct {
	Type         string    `json:"type"`
	ChatJID      string    `json:"chat"`
	Text         string    `json:"text"`
	IsFromMe     bool      `json:"isFromMe"`
	Timestamp    time.Time `json:"time"`
	MessageID    int       `json:"message_id"`
	SessionToken string    `json:"session_token"`
}

func (m *Message) ToDTO(t string) *WSMessageDTO {
	return &WSMessageDTO{
		Type:         t,
		ChatJID:      m.ChatJID,
		Text:         m.Text,
		IsFromMe:     m.IsFromMe,
		Timestamp:    m.Timestamp,
		MessageID:    m.ID,
		SessionToken: m.SessionToken,
	}
}
