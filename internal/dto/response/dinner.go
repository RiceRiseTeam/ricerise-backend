package response

import (
	"ricerise/internal/model"
)

type DinnerFindResponse struct {
	Result []model.DinnerModel `json:"result"`
}

type ParticipantFindResponse struct {
	Result []model.ParticipantModel `json:"result"`
}
