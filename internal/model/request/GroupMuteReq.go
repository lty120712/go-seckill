package model

type GroupMuteRequest struct {
	GroupId uint  `json:"group_id" binding:"required"`
	MuteEnd int64 `json:"mute_end" binding:"required"`
}
