package model

type GroupUpdateRequest struct {
	GroupId uint    `json:"group_id" binding:"required"`
	Name    *string `json:"name"`
	Avatar  *string `json:"avatar"`
	Desc    *string `json:"desc"`
	MaxNum  *int    `json:"max_num"`
}
