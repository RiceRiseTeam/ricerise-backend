package model

import "time"

type ChatMessageModel struct {
	ID        uint64 `gorm:"primaryKey"`
	SessionId uint64 `gorm:"index:session_user"`
	UserId    uint64 `gorm:"index:session_user"`
	Role      string `gorm:"not null"`
	Content   string `gorm:"type:text"`
	CreatedAt time.Time
}
