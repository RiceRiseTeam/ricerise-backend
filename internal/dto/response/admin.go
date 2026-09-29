package response

import "ricerise/internal/dto"

// AdminStatusResponse 管理员状态响应
// @Description 管理员状态响应结构体
type AdminStatusResponse struct {
	CurrentDinner int `json:"current_dinner"`
	TotalDinner   int `json:"total_dinner"`
	TotalLocation int `json:"total_location"`
	TotalUser     int `json:"total_user"`
}

// AdminCommentReviewList 管理员评论审核列表响应
// @Description 管理员评论审核列表响应结构体
type AdminCommentReviewList struct {
	PageSize int               `json:"page_size"`
	HasNext  bool              `json:"has_next"`
	Comments []*dto.CommentDto `json:"comments"`
}

// AdminLocationReviewList 管理员地点审核列表响应
// @Description 管理员地点审核列表响应结构体
type AdminLocationReviewList struct {
	PageSize  int                `json:"page_size"`
	HasNext   bool               `json:"has_next"`
	Locations []*dto.LocationDto `json:"locations"`
}
