// Package model models.response
package model

// Response 通用的 API 响应体结构体
// @Description 通用响应格式，用于包装所有 API 响应
type Response struct {
	Code    int         `json:"code"`    // 状态码
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 数据部分
}
