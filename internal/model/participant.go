package model

import "time"

type ParticipantModel struct {
	ID       uint64      `gorm:"primary_key"`
	DinnerId uint64      `gorm:"index;uniqueIndex:idx_dinner_user"`
	Dinner   DinnerModel `gorm:"foreignKey:DinnerId"`
	UserId   uint64      `gorm:"index;uniqueIndex:idx_dinner_user"`
	User     UserModel   `gorm:"foreignKey:UserId"`
	Status   int8        `gorm:"default:0;index"` // 0: 已报名 1: 已参加 2: 已退出

	CreatedAt time.Time
	UpdatedAt time.Time
}
