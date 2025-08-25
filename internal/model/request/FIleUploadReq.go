package model

type FileUploadRequest struct {
	Type string `form:"type" binding:"required,oneof=image audio video file"` // 文件类型，限定四类
}
