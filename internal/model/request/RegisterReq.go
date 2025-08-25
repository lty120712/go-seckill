package model

// RegisterRequest 注册请求体
type RegisterRequest struct {
	Username   *string `json:"username" binding:"required"`    // 用户名
	Password   *string `json:"password" binding:"required"`    // 密码
	RePassword *string `json:"re_password" binding:"required"` // 确认密码
}
