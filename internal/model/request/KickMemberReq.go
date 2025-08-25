package model

type KickMemberRequest struct {
	MemberID uint `json:"member_id" binding:"required"` // 被踢的成员ID
}
