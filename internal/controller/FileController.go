package controller

import (
	"github.com/gin-gonic/gin"
	interfacesservice "go-chat/internal/interfaces/service"
)

// FileController 文件相关控制器
// @Tags File
// @Description 控制文件相关的 API
type FileController struct {
	BaseController
	fileService interfacesservice.FileServiceInterface
}

var FileControllerInstance *FileController

func InitFileController(fileService interfacesservice.FileServiceInterface) {
	FileControllerInstance = &FileController{
		fileService: fileService,
	}
}

func (con FileController) Upload(c *gin.Context) {
	id := c.GetUint("id")
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		con.Error(c, "上传文件读取失败: "+err.Error())
		return
	}
	url, err := con.fileService.Upload(id, file)
	if err != nil {
		con.Error(c, err.Error())
		return
	}
	if err != nil {
		return
	}
	con.Success(c, url)
}
