package domain

import "time"

type Chat struct {
	ID           int    `gorm:"primaryKey"`
	SessionToken string `gorm:"index"`
	JID          string `gorm:"index"`
	Title        string
	UpdatedAt    time.Time
}

type Message struct {
	ID           int    `gorm:"primaryKey"`
	SessionToken string `gorm:"index"`
	ChatJID      string `gorm:"index"`
	FromMe       bool
	Text         string
	Timestamp    time.Time `gorm:"index"`
}

// RawMessage используется на этапе HistorySync до трансформации в Message
// например: когда приходит *waProto.WebMessageInfo из WhatsApp
// и нужно собрать JID, текст, from_me, timestamp

type RawMessage struct {
	JID       string
	FromMe    bool
	Text      string
	Timestamp time.Time
}

type WSMessageDTO struct {
	Type         string    `json:"type"`
	ChatJID      string    `json:"chat"`
	From         string    `json:"from"`
	Text         string    `json:"text"`
	FromMe       bool      `json:"me"`
	Timestamp    time.Time `json:"time"`
	MessageID    int       `json:"message_id"`
	SessionToken string    `json:"session_token"`
}

type WebSocketCommand struct {
	Type    string `json:"type"`
	ChatJID string `json:"chat"`
	Text    string `json:"text"`
}
