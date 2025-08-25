package controller

import (
	"github.com/gin-gonic/gin"
	interfacesservice "go-chat/internal/interfaces/service"
	request "go-chat/internal/model/request"
	"strconv"
)

// MessageController 消息相关控制器
// @Tags Message
// @Description 消息相关控制器
type MessageController struct {
	BaseController
	messageService interfacesservice.MessageServiceInterface
}

var MessageControllerInstance *MessageController

func InitMessageController(messageService interfacesservice.MessageServiceInterface) {
	MessageControllerInstance = &MessageController{
		messageService: messageService,
	}
}

// SendJson 已读消息接口
// @Summary 已读消息
// @Description 对该条消息已读
// @Tags Message
// @Accept json
// @Produce json
// @Param msg body model.ReadMessageReq true "消息内容"
// @Success 200 {object} model.Response "成功"
// @Failure 500 {object} model.Response "发送消息失败"
// @Router /message/read [post]
func (con MessageController) Read(c *gin.Context) {
	var req request.ReadMessageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, "参数错误")
		return
	}
	if err := con.messageService.ReadMessage(req.MessageId, req.UserId); err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c)
}

// Query godoc
// @Summary 查询历史消息（分页，支持游标分页）
// @Description 根据目标ID和目标类型查询聊天消息历史，支持分页、时间范围等过滤
// @Tags Message
// @Accept application/json
// @Produce application/json
// @Param query body model.QueryMessagesRequest true "查询参数"
// @Success 200 {object} model.QueryMessagesResponse "查询成功，返回消息列表及分页信息"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /message/query [post]
func (con MessageController) Query(c *gin.Context) {
	var req request.QueryMessagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}
	idStr, _ := c.Get("id")
	id := idStr.(uint)
	if data, err := con.messageService.QueryMessages(id, &req); err != nil {
		con.Error(c, err.Error())
		return
	} else {
		con.Success(c, data)
	}
}

// Revoke 撤回消息接口
// @Summary 撤回消息
// @Description 根据消息ID撤回指定消息，只有发送者或管理员才能撤回消息
// @Tags Message
// @Accept json
// @Produce json
// @Param id path int true "消息ID"  // 消息ID，作为URL路径参数
// @Success 200 {object} model.Response "成功"  // 成功的撤回响应
// @Failure 500 {object} model.Response "撤回失败"  // 撤回消息失败
// @Router /message/{id}/revoke [get]
func (con MessageController) Revoke(c *gin.Context) {
	messageIdStr := c.Param("id")
	messageId, _ := strconv.Atoi(messageIdStr)
	userId := c.GetUint("id")
	if err := con.messageService.Revoke(userId, uint(messageId)); err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c)
}
