package domain

type Session struct {
	ID         int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Token      string `json:"token" gorm:"unique;not null"`
	BusinessID int    `json:"business_id" gorm:"unique;index"`
	JID        string `json:"jid" gorm:"index"`
	Active     bool   `json:"active" gorm:"default:false"`
}

func (Session) TableName() string { return "sessions" }
