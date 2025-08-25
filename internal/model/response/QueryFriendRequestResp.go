package model

import (
	"go-chat/internal/model"
	"time"
)

type FriendRequestVo struct {
	Id       int64        `json:"id"`        // 申请Id
	TargetId int64        `json:"target_id"` //对方的Id
	NickName string       `json:"nick_name"` //对方昵称
	Avatar   *string      `json:"avatar"`    //对方头像
	Sent     int          `json:"sent"`      // 是否由自己发出 1 是 0否
	Status   model.Status `json:"status"`    //申请状态
	Time     time.Time    `json:"time"`      //申请时间
}
