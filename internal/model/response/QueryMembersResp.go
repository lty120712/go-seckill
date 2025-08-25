package model

import (
	"go-chat/internal/model"
	"time"
)

type MemberVo struct {
	GroupId      uint               ` json:"group_id"`
	UserId       uint               `json:"user_id"`
	Nickname     string             `json:"nickname"`
	Role         model.Role         `json:"role"`
	MuteEnd      *time.Time         `json:"mute_end"`
	Avatar       *string            `json:"avatar"`
	OnlineStatus model.OnlineStatus `json:"online_status"`
}
