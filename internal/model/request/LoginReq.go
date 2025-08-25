package model

// 登录请求体
type LoginRequest struct {
	Username *string `json:"username" binding:"required"` // 用户名
	Password *string `json:"password" binding:"required"` // 密码
}
