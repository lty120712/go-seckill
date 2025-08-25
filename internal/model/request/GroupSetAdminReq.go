package model

type SetAdminRequest struct {
	MemberID uint `json:"member_id" binding:"required"`
}
