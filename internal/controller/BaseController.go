package controller

import (
	"github.com/gin-gonic/gin"
	models "go-chat/internal/model"
	"net/http"
)

type BaseController struct {
}

func (_ *BaseController) Success(c *gin.Context, data ...interface{}) {
	response := models.Response{}
	if len(data) > 1 {
		response.Code = http.StatusInternalServerError
		response.Message = "Success函数只允许传入 0 或 1 个参数"
		response.Data = nil
	}
	if len(data) == 1 {
		response.Code = http.StatusOK
		response.Message = "success"
		response.Data = data[0]
	}
	if len(data) == 0 {
		response.Code = http.StatusOK
		response.Message = "success"
		response.Data = nil
	}
	c.JSON(http.StatusOK, response)
	return
}

// 错误处理 (c,[错误消息][错误码])
func (_ *BaseController) Error(c *gin.Context, data ...interface{}) {
	response := models.Response{}
	switch len(data) {
	case 0:
		break
	case 1:
		response.Code = http.StatusInternalServerError
		if msg, ok := data[0].(string); ok {
			response.Message = msg
		} else {
			response.Message = "Invalid error message"
		}
	case 2:
		if msg, ok := data[0].(string); ok {
			response.Message = msg
		} else {
			response.Message = "Invalid error message"
		}
		if code, ok := data[1].(int); ok {
			response.Code = code
		} else {
			response.Code = http.StatusInternalServerError
		}
	default:
		response.Code = http.StatusInternalServerError
		response.Message = "Invalid number of parameters"
	}
	c.JSON(http.StatusOK, response)
	return
}
