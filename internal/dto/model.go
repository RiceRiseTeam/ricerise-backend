package dto

import (
	"ricerise/internal/model"
	"time"
)

// CommentDto 评论结构
// @Description 评论的数据结构
type CommentDto struct {
	ID        uint64       `json:"id"`
	Location  *LocationDto `json:"location"`
	User      *UserDto     `json:"user"`
	Rating    int          `json:"rating"`
	Content   string       `json:"content"`
	CreatedAt time.Time    `json:"createdAt"`
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

// UserDto 用户结构
// @Description 用户的数据结构
type UserDto struct {
	ID              uint64    `json:"id"`
	Username        string    `json:"username"`
	Nickname        string    `json:"nickname"`
	PermissionLevel int8      `json:"permissionLevel"`
	CreatedAt       time.Time `json:"createdAt"`
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

// LocationDto 地点结构
// @Description 地点的数据结构
type LocationDto struct {
	ID uint64 `json:"id"`

	Name        string  `json:"name"`
	Address     string  `json:"address"`
	Longitude   float64 `json:"longitude"`
	Latitude    float64 `json:"latitude"`
	Description string  `json:"description"`

	CreatedAt time.Time `json:"createdAt"`
}

func NewLocationDto(location *model.LocationModel) *LocationDto {
	return &LocationDto{
		ID:          location.ID,
		Name:        location.Name,
		Address:     location.Address,
		Longitude:   location.Location.Lng,
		Latitude:    location.Location.Lat,
		Description: location.Description,
		CreatedAt:   location.CreatedAt,
	}
}

// DinnerDto 饭局结构
// @Description 饭局的数据结构
type DinnerDto struct {
	ID           uint64       `json:"id"`
	Location     *LocationDto `json:"location"`
	Host         *UserDto     `json:"host"`
	Participants []*UserDto   `json:"participants"`
	MaxPeople    int          `json:"max"`
	MeetTime     time.Time    `json:"meet_time"`
	Status       int8         `json:"status"`
	CreatedAt    time.Time    `json:"createdAt"`
}

func NewDinnerDto(dinner *model.DinnerModel) *DinnerDto {
	return &DinnerDto{
		ID:        dinner.ID,
		Location:  NewLocationDto(&dinner.Location),
		Host:      NewUserDto(&dinner.Host),
		MaxPeople: dinner.MaxPeople,
		CreatedAt: dinner.CreatedAt,
		MeetTime:  dinner.MeetTime,
		Status:    dinner.Status,
		Participants: Map(dinner.Participants, func(m model.ParticipantModel) *UserDto {
			return NewUserDto(&m.User)
		}),
	}
}

// ChatMessageDto AI聊天消息结构
// @Description AI聊天消息的数据结构
type ChatMessageDto struct {
	ID        uint64 `json:"id"`
	SessionId uint64 `json:"session_id"`
	Role      string `json:"role"`
	Content   string `json:"content"`
}

func NewChatMessageDto(chatMessage *model.ChatMessageModel) *ChatMessageDto {
	return &ChatMessageDto{
		ID:        chatMessage.ID,
		SessionId: chatMessage.SessionId,
		Role:      chatMessage.Role,
		Content:   chatMessage.Content,
	}
}
