package model

import (
	"time"

	"gorm.io/gorm"
)

type ParticipantModel struct {
	ID       uint64      `gorm:"primary_key"`
	DinnerId uint64      `gorm:"index;uniqueIndex:idx_dinner_user"`
	Dinner   DinnerModel `gorm:"foreignKey:DinnerId"`
	UserId   uint64      `gorm:"index;uniqueIndex:idx_dinner_user"`
	User     UserModel   `gorm:"foreignKey:UserId"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
