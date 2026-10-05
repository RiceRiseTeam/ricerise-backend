package request

import "time"

type JoinDinnerRequest struct {
	Code string `json:"code"`
}
type CreateDinnerRequest struct {
	LocationId uint64    `json:"location_id" binding:"required"`
	MeetTime   time.Time `json:"meet_time" binding:"required"`
	MaxPeople  int       `json:"max_people" binding:"required,min=1,max=64"`
}

type DinnerRoomChatRequest struct {
	Message string `json:"message" binding:"required,min=1,max=128"`
}

type UpdateDinnerStatusRequest struct {
	Status int8 `json:"status" binding:"required,oneof=2 3 4"`
}

type DinnerLikeFindRequest struct {
	LocationId uint64 `json:"dinner_id" binding:"required"`
}

type NewParticipateRequest struct {
	DinnerID uint64 `json:"dinner_id" binding:"required"`
}
