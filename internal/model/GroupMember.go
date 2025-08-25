package model

import (
	"gorm.io/gorm"
	"time"
)

type GroupMember struct {
	gorm.Model
	GroupId   uint       `json:"group_id"`    //群ID
	MemberId  uint       `json:"member_id"`   //成员ID
	GNickName string     `json:"g_nick_name"` //群昵称
	Role      Role       `json:"role"`        //成员角色（0=普通成员,1=群主，2=管理员）
	MuteEnd   *time.Time `json:"mute_end"`    //禁言截至时间(null未禁言)
}

func (gm *GroupMember) TableName() string {
	return "group_members"
}

func (gm *GroupMember) IsOwner() bool {
	return gm.Role == 1
}

func (gm *GroupMember) IsAdmin() bool {
	return gm.Role == 2
}

func (gm *GroupMember) IsMuted() bool {
	return gm.MuteEnd != nil && gm.MuteEnd.After(time.Now())
}
