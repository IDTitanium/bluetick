package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID         uint64         `gorm:"primaryKey" json:"id"`
	SenderID   uint64         `json:"sender_id"`
	ReceiverID uint64         `json:"receiver_id"`
	Content    string         `json:"content"`
	IsRead     bool           `gorm:"default:false" json:"is_read"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
