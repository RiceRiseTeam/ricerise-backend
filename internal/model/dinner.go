package model

import "time"

type DinnerModel struct {
	ID uint64 `gorm:"primary_key"`

	LocationId uint64        `gorm:"not null"`
	Location   LocationModel `gorm:"foreignkey:LocationId"`
	HostId     uint64        `gorm:"not null"`
	Host       UserModel     `gorm:"foreignkey:HostId"`
	MeetTime   time.Time     `gorm:"not null"`
	MaxPeople  int           `gorm:"not null,default:2"`
	CreatedAt  time.Time
}
