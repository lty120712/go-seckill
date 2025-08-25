package model

type GroupMemberMuteRequest struct {
	MemberID uint  `json:"member_id" binding:"required"` // 要禁言的成员ID
	Duration int64 `json:"duration" binding:"required"`  // 禁言时长（秒）
}
