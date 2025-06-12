package infrastructure

import (
	"fmt"
	"github.com/chatbox/whatsapp/internal/config"
	"github.com/chatbox/whatsapp/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DBConfig.DatabaseUrl), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}
	if err := db.AutoMigrate(&domain.Session{}, &domain.Chat{}, &domain.Message{}); err != nil {
		return nil, fmt.Errorf("failed to migrate DB: %w", err)
	}
	return db, nil
}
