package model

import "gorm.io/gorm"

type File struct {
	gorm.Model
	UserID   uint     `json:"user_id"` // 上传者ID
	Type     string   `gorm:"type:enum('image','audio','video','file');not null" json:"type"`
	Name     string   `gorm:"size:255;not null" json:"name"` // 原始文件名
	Ext      string   `gorm:"size:20" json:"ext"`            // 文件扩展名
	Mime     string   `gorm:"size:100" json:"mime"`          // MIME 类型
	Size     uint64   `json:"size"`                          // 文件大小（字节）
	Url      string   `gorm:"type:text;not null" json:"url"` // 访问地址
	Width    *uint    `json:"width,omitempty"`               // 图像/视频宽度
	Height   *uint    `json:"height,omitempty"`              // 图像/视频高度
	Duration *float64 `json:"duration,omitempty"`            // 音/视频时长（秒）
}
