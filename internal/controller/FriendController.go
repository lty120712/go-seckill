package controller

import (
	"github.com/gin-gonic/gin"
	interfacesservice "go-chat/internal/interfaces/service"
	request "go-chat/internal/model/request"
	"strconv"
)

// FriendController 用户相关控制器
// @Tags Friend
// @Description 控制好友相关的 API
type FriendController struct {
	BaseController
	friendService interfacesservice.FriendServiceInterface
}

var FriendControllerInstance *FriendController

func InitFriendController(friendService interfacesservice.FriendServiceInterface) {
	FriendControllerInstance = &FriendController{
		friendService: friendService,
	}
}

// Add 添加好友申请
// @Summary 添加好友申请
// @Description 当前用户向其他用户发送好友申请
// @Tags Friend
// @Accept json
// @Produce json
// @Param data body []uint true "好友用户ID数组"
// @Success 200 {object} model.Response
// @Router /friend/add [post]
func (con *FriendController) Add(c *gin.Context) {
	id := c.GetUint("id")
	var req []uint
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}
	if err := con.friendService.Add(id, req); err != nil {
		con.Error(c, err.Error())
		return
	}

	con.Success(c)
}

// ListReq 获取收到的好友申请
// @Summary 获取好友申请列表
// @Description 获取当前用户收到的所有好友申请
// @Tags Friend
// @Produce json
// @Success 200 {object} model.Response{data=[]model.FriendRequestVo}
// @Router /friend/list_req [get]
func (con *FriendController) ListReq(c *gin.Context) {
	id := c.GetUint("id")
	data, err := con.friendService.ListReq(id)
	if err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c, data)
}

// HandleReq 处理好友申请
// @Summary 同意或拒绝好友申请
// @Description 同意或拒绝某条好友请求
// @Tags Friend
// @Accept json
// @Produce json
// @Param data body request.FriendHandlerReq true "处理请求信息"
// @Success 200 {object} model.Response
// @Router /friend/handle_req [post]
func (con *FriendController) HandleReq(c *gin.Context) {

	var req request.FriendHandlerReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		con.Error(c, err.Error())
		return
	}
	if err := con.friendService.HandleReq(req.Id, req.Status); err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c)
}

// Remove 删除好友
// @Summary 删除好友
// @Description 删除当前用户的好友
// @Tags Friend
// @Accept json
// @Produce json
// @Param data body []int64 true "好友ID数组"
// @Success 200 {object} model.Response
// @Router /friend/remove [post]
func (con *FriendController) Remove(c *gin.Context) {
	id := c.GetUint("id")
	var req []int64
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}
	if err := con.friendService.Remove(id, req); err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c)
}

// GroupCreate 创建好友分组
// @Summary 创建好友分组
// @Description 给好友添加分组标签
// @Tags Friend
// @Accept json
// @Produce json
// @Param data body request.FriendGroupCreateRequest true "分组创建请求"
// @Success 200 {object} model.Response
// @Router /friend/group_create [post]
func (con *FriendController) GroupCreate(c *gin.Context) {
	id := c.GetUint("id")
	var req request.FriendGroupCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}
	if err := con.friendService.GroupCreate(id, req); err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c)
}

// GroupDelete 删除好友分组
// @Summary 删除好友分组
// @Description 根据分组 ID 删除好友分组
// @Tags Friend
// @Produce json
// @Param id query int true "分组ID"
// @Success 200 {object} model.Response
// @Router /friend/group_delete [get]
func (con *FriendController) GroupDelete(c *gin.Context) {
	groupIdStr := c.Query("id")
	groupId, _ := strconv.ParseInt(groupIdStr, 10, 64)
	if err := con.friendService.GroupDelete(groupId); err != nil {
		con.Error(c, err.Error())
		return
	}
}

// GroupUpdate 修改好友分组
// @Summary 修改好友分组信息
// @Description 修改分组的名称等信息
// @Tags Friend
// @Accept json
// @Produce json
// @Param data body request.FriendGroupUpdateRequest true "分组更新请求"
// @Success 200 {object} model.Response
// @Router /friend/group_update [post]
func (con *FriendController) GroupUpdate(c *gin.Context) {
	var req request.FriendGroupUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}

	if err := con.friendService.GroupUpdate(req); err != nil {
		con.Error(c, err.Error())
		return
	}

	con.Success(c)
}

// GroupList 查询好友分组
// @Summary 查询当前用户的所有好友分组
// @Description 返回当前用户的分组及其成员情况
// @Tags Friend
// @Produce json
// @Success 200 {object} model.Response{data=[]model.GroupVo}
// @Router /friend/group_list [get]
func (con *FriendController) GroupList(c *gin.Context) {
	userId := c.GetUint("id")
	data, err := con.friendService.GroupList(userId)
	if err != nil {
		con.Error(c, err.Error())
		return
	}

	con.Success(c, data)
}
