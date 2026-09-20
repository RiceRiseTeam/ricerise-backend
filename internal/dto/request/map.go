package request

type UploadLocationRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=128"`
	Longitude   float64 `json:"longitude" binding:"required"`
	Latitude    float64 `json:"latitude" binding:"required"`
	Description string  `json:"description" binding:"required,min=2,max=1024"`
}

type UploadCommentRequest struct {
	Rating     int    `json:"rating" binding:"required"`
	Content    string `json:"content" binding:"required,min=1,max=1024"`
	LocationId uint64 `json:"location_id" binding:"required"`
}
