package model

import "time"

type CommentModel struct {
	ID uint64 `gorm:"primary_key"`

	LocationId uint64        `gorm:"index,not null"`
	Location   LocationModel `gorm:"foreignkey:LocationId"`
	DinnerId   uint64        `gorm:"not null"`
	Dinner     DinnerModel   `gorm:"foreignkey:DinnerId"`
	UserId     uint64        `gorm:"index;not null"`
	User       UserModel     `gorm:"foreignkey:UserId"`
	Rating     int           `gorm:"not null"`
	Content    string        `gorm:"type:text;not null"`
	Reviewed   *bool         `gorm:"not null;default:false"`
	CreatedAt  time.Time
}
