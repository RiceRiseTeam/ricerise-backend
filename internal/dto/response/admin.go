package response

import "ricerise/internal/dto"

type AdminStatusResponse struct {
	CurrentDinner int `json:"current_dinner"`
	TotalDinner   int `json:"total_dinner"`
	TotalLocation int `json:"total_location"`
	TotalUser     int `json:"total_user"`
}

type AdminCommentReviewList struct {
	PageSize int               `json:"page_size"`
	HasNext  bool              `json:"has_next"`
	Comments []*dto.CommentDto `json:"comments"`
}

type AdminLocationReviewList struct {
	PageSize  int                `json:"page_size"`
	HasNext   bool               `json:"has_next"`
	Locations []*dto.LocationDto `json:"locations"`
}
