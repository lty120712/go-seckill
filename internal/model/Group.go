package model

import (
	"gorm.io/gorm"
	"time"
)

type Group struct {
	gorm.Model
	Code          string     `json:"code"`           //群号
	Name          string     `json:"name"`           //群组名称
	Avatar        string     `json:"avatar"`         //群组头像
	Desc          string     `json:"description"`    //群组简介
	OwnerId       uint       `json:"owner_id"`       //群主ID
	MaxNum        int        `json:"max_num"`        //群组最大人数
	Status        Status     `json:"status" `        //群组状态（1=正常，0=关闭）
	MuteEnd       *time.Time `json:"mute_end"`       // 禁言截至时间
	LimitInterval *int       `json:"limit_interval"` // 限制发送消息的间隔时间
	LimitCount    *int       `json:"limit_count"`    // 限制发送消息的次数
}

func (m *Group) TableName() string {
	return "groups"
}
