package model

import (
	"database/sql/driver"
	"go-chat/internal/utils/jsonUtil"
	"gorm.io/gorm"
)

type FriendGroup struct {
	gorm.Model
	UserId       uint          `json:"user_id"` // 用户id
	Name         *string       `json:"name"`    // 群组名称
	FriendIdList *FriendIdList `json:"friend_id_list"`
}

// 群组成员id列表
type FriendIdList []uint

func (f *FriendIdList) Scan(value interface{}) error {
	return jsonUtil.UnmarshalValue(value, f)
}
func (f *FriendIdList) Value() (driver.Value, error) {
	return jsonUtil.MarshalValue(f)
}
