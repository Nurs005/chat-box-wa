package domain

import "time"

type Chat struct {
	ID           int       `gorm:"primaryKey"`
	SessionToken string    `gorm:"size:255;index"`
	JID          string    `gorm:"size:255;index"`
	Title        string    `gorm:"size:255"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func (Chat) TableName() string { return "chats" }

type Message struct {
	ID           int    `gorm:"primaryKey"`
	SessionToken string `gorm:"size:255;index"`
	ChatJID      string `gorm:"size:255;index"`
	FromMe       bool
	Text         string    `gorm:"type:text"`
	Timestamp    time.Time `gorm:"index"`
}

func (Message) TableName() string { return "messages" }

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
