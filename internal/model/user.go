package model

type UserModel struct {
	ID       uint64 `gorm:"primary_key"`
	Username string `gorm:"type:varchar(64);not null;uniqueIndex"`
	Nickname string `gorm:"type:varchar(64);not null"`
	Password string `gorm:"type:varchar(256);not null"`
	Email    string `gorm:"type:varchar(256);not null"`
}
