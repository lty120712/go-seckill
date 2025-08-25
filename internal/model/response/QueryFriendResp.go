package model

import "go-chat/internal/model"

type GroupVo struct {
	GroupId      uint          `json:"group_id"`
	Name         *string       `json:"name"`
	FriendVoList *FriendVoList `json:"friend_vo_list"`
}

type FriendVo struct {
	UserId       uint               `json:"user_id"`       // 用户ID
	Nickname     *string            `json:"nickname"`      // 昵称
	Desc         *string            `json:"desc"`          // 简介
	OnlineStatus model.OnlineStatus `json:"online_status"` // 用户在线状态（如在线、离线、忙碌等）
	DeviceInfo   *string            `json:"device_info"`   // 客户端设备信息
}

type FriendVoList []FriendVo
