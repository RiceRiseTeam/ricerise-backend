package request

// UploadLocationRequest 上传地点请求
// @Description 上传地点请求结构体
type UploadLocationRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=128"`
	Address     string  `json:"address" binding:"required,min=2,max=1024"`
	Longitude   float64 `json:"longitude" binding:"required"`
	Latitude    float64 `json:"latitude" binding:"required"`
	Description string  `json:"description" binding:"required,min=2,max=1024"`
}

// UploadCommentRequest 上传评论请求
// @Description 上传评论请求结构体
type UploadCommentRequest struct {
	Rating  int    `json:"rating" binding:"required"`
	Content string `json:"content" binding:"required,min=1,max=1024"`
}

// GetLocationsRequest 获取地点请求
// @Description 获取地点请求结构体
type GetLocationsRequest struct {
	MinLng float64 `json:"min_lng" binding:"required"`
	MaxLng float64 `json:"max_lng" binding:"required"`
	MinLat float64 `json:"min_lat" binding:"required"`
	MaxLat float64 `json:"max_lat" binding:"required"`
}
