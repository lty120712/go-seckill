package model

import "go-chat/internal/model"

type UserVO struct {
	Id            uint               `json:"id"`                               // 用户ID
	Username      string             `json:"username"`                         // 用户名
	Nickname      *string            `json:"nickname"`                         // 昵称
	Desc          *string            `json:"desc"`                             // 简介
	Phone         *string            `json:"phone" validate:"omitempty,phone"` // 用户手机号
	Email         *string            `json:"email" validate:"omitempty,email"` // 用户邮箱
	Avatar        *string            `json:"avatar,omitempty"`                 // 用户头像URL
	ClientIp      string             `json:"client_ip"`                        // 客户端IP地址
	ClientPort    string             `json:"client_port"`                      // 客户端端口号
	LoginTime     int64              `json:"login_time"`                       // 最近一次登录时间
	HeartbeatTime int64              `json:"heartbeat_time"`                   // 最近一次心跳时间
	LogoutTime    int64              `json:"logout_time"`                      // 最近一次登出时间
	Status        model.Status       `json:"status"`                           // 用户状态（如激活、禁用等）
	OnlineStatus  model.OnlineStatus `json:"online_status"`                    // 用户在线状态（如在线、离线、忙碌等）
	DeviceInfo    *string            `json:"device_info"`                      // 客户端设备信息
}
