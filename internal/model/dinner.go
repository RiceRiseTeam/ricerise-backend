package model

import "time"

const (
	DINNER_HIRING = iota
	DINNER_FULL
	DINNER_ONGOING
	DINNER_FINISHED
	DINNER_CANCELLED
)

type DinnerModel struct {
	ID         uint64        `gorm:"primary_key"`
	LocationId uint64        `gorm:"not null"`
	Location   LocationModel `gorm:"foreignkey:LocationId"`
	HostId     uint64        `gorm:"not null"`
	Host       UserModel     `gorm:"foreignkey:HostId"`
	MeetTime   time.Time     `gorm:"not null"`
	MaxPeople  int           `gorm:"not null,default:2"`
	Status     int8          `gorm:"default:0"` // 0: 招募中 1: 已满员 2: 进行中 3: 已结束 4： 已取消
	CreatedAt  time.Time
}
