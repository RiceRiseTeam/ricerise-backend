package request

type DinnerLikeFindRequest struct {
	LocationName string `json:"location_name"`
}

type NewParticipateRequest struct {
	DinnerID uint64 `json:"dinner_id" binding:"required"`
}
