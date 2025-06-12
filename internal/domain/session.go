package domain

type Session struct {
	ID         int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Token      string `json:"token" gorm:"size:255;uniqueIndex;not null"`
	BusinessID int    `json:"business_id" gorm:"uniqueIndex;not null"`
	JID        string `json:"jid" gorm:"size:255;index"`
	Active     bool   `json:"active" gorm:"default:false;not null"`
}

func (Session) TableName() string { return "sessions" }
