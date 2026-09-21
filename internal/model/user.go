package model

import (
	"time"
)

type UserModel struct {
	ID              uint64             `gorm:"primary_key"`
	Username        string             `gorm:"type:varchar(64);not null;uniqueIndex"`
	Nickname        string             `gorm:"type:varchar(64);not null"`
	Password        string             `gorm:"type:varchar(256);not null"`
	Email           string             `gorm:"type:varchar(256);not null"`
	PermissionLevel int8               `gorm:"default:0;index"`
	Comments        []CommentModel     `gorm:"foreignkey:UserId"`
	Participants    []ParticipantModel `gorm:"foreignkey:UserId"`
	CreatedAt       time.Time
}
