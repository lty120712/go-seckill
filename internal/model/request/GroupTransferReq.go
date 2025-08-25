package model

type GroupTransferRequest struct {
	GroupID    uint `json:"group_id" binding:"required"`     // 群ID
	NewOwnerID uint `json:"new_owner_id" binding:"required"` // 新群主的用户ID
}
