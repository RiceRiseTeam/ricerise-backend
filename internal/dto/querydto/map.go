package querydto

import "time"

type MapLocationCommentsQuery struct {
	StartId   *uint64    `form:"start_id"`
	StartTime *time.Time `form:"start_time"`
	StartRank *int       `form:"start_rank"`
	OrderedBy string     `form:"ordered_by,default=time" binding:"omitempty,oneof=time rank"`
	PageSize  int        `form:"page_size,default=10"`
}
