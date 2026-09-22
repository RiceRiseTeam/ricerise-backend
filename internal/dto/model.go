package dto

import (
	"ricerise/internal/model"
	"time"
)

type CommentDto struct {
	ID        uint64
	Location  *LocationDto
	User      *UserDto
	Rating    int
	Content   string
	CreatedAt time.Time
}

func NewCommentDto(comment *model.CommentModel) *CommentDto {
	return &CommentDto{
		ID:        comment.ID,
		Location:  NewLocationDto(&comment.Location),
		User:      NewUserDto(&comment.User),
		Rating:    comment.Rating,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
	}
}

type UserDto struct {
	ID              uint64
	Username        string
	Nickname        string
	PermissionLevel int8
	CreatedAt       time.Time
}

func NewUserDto(user *model.UserModel) *UserDto {
	return &UserDto{
		ID:              user.ID,
		Username:        user.Username,
		Nickname:        user.Nickname,
		PermissionLevel: user.PermissionLevel,
		CreatedAt:       user.CreatedAt,
	}
}

type LocationDto struct {
	ID uint64

	Name        string
	Longitude   float64
	Latitude    float64
	Description string

	CreatedAt time.Time
}

func NewLocationDto(location *model.LocationModel) *LocationDto {
	return &LocationDto{
		ID:          location.ID,
		Name:        location.Name,
		Longitude:   location.Location.Lng,
		Latitude:    location.Location.Lat,
		Description: location.Description,
		CreatedAt:   location.CreatedAt,
	}
}
