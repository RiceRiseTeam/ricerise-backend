package query

import "time"

type AdminPageQuery struct {
	StartId   *uint64    `form:"start_id"`
	StartTime *time.Time `form:"start_time"`
	PageSize  int        `form:"page_size,default=10" binding:"omitempty,min=1,max=20"`
}

type AdminReviewQuery struct {
	Pass bool `form:"pass" binding:"required"`
}
