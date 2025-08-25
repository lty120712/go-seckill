package model

// Status 用户状态
type Status int

const (
	Disable Status = iota // 0 禁用
	Enable                // 1 启用
)

const (
	Todo   Status = iota // 0 待处理
	Accept               // 1 接受
	Reject               // 2 拒绝
)

// OnlineStatus 在线状态
type OnlineStatus int

const (
	Offline OnlineStatus = iota // 0 离线
	Online                      // 1 在线
	Busy                        // 2 忙碌
	Away                        // 3 离开
)

// Role 群成员角色

type Role int

const (
	Member Role = iota // 0 成员
	Admin              // 1 管理员
	Owner              // 2 群主
)
