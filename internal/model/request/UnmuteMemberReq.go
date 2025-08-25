package model

type UnmuteMemberRequest struct {
	MemberID uint `json:"member_id" binding:"required"`
}
