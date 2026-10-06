package model

import "time"

type ChatMessageModel struct {
	ID        uint64 `gorm:"primaryKey"`
	SessionId uint64 `gorm:"index"`
	Role      string `gorm:"not null"`
	Content   string `gorm:"type:text"`
	CreatedAt time.Time
}
