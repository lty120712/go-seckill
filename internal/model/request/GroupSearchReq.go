package model

// GroupSearchRequest 群搜索请求
type GroupSearchRequest struct {
	Code     string `json:"code"`    // 群号，精准查询
	Name     string `json:"name"`    // 群名称，模糊查询
	UserId   uint   `json:"user_id"` // 用户ID(我加入的群)
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}
