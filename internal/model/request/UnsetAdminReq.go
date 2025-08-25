package model

type UnsetAdminRequest struct {
	MemberID uint `json:"member_id" binding:"required"` // 被取消管理员权限的用户 ID
}
