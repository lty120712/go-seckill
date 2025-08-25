package model

type GroupAnnouncementUpdateRequest struct {
	AnnouncementId int64  `json:"announcement_id"` // 公告ID
	Content        string `json:"content"`         // 公告内容
}
