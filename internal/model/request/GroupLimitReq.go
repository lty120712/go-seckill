package model

type GroupLimitRequest struct {
	GroupId       uint `json:"group_id" binding:"required"`
	LimitInterval *int `json:"limit_interval"`
	LimitCount    *int `json:"limit_count"`
}
