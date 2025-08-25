package model

type OnlineStatusNotice struct {
	UserId       uint         // 用户id
	OnlineStatus OnlineStatus // 在线状态
	ActionType   ActionType   // 操作类型，如：login, logout, status_change
}

type ActionType int

const (
	LoginAction        ActionType = iota // 登录
	LogoutAction                         // 登出
	StatusChangeAction                   // 主动状态变化

	HeartbeatAction // 心跳检测
)
