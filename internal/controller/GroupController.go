package controller

import (
	"github.com/gin-gonic/gin"
	interfacesservice "go-chat/internal/interfaces/service"
	request "go-chat/internal/model/request"
	"strconv"
)

// GroupController 群组相关控制器
// @Tags Group
// @Description 群组相关控制器
type GroupController struct {
	BaseController
	groupService interfacesservice.GroupServiceInterface
}

var GroupControllerInstance *GroupController

func InitGroupController(groupService interfacesservice.GroupServiceInterface) {
	GroupControllerInstance = &GroupController{
		groupService: groupService,
	}
}

// Create 创建群组
// @Summary 创建群组
// @Description 创建一个新的群组，并将指定的成员添加到群组中
// @Tags Group
// @Accept json
// @Produce json
// @Param group body model.GroupCreateRequest true "创建群组请求参数"
// @Success 200 {object} model.Response "群组创建成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 500 {object} model.Response "内部服务器错误"
// @Router /group/create [post]
func (con GroupController) Create(c *gin.Context) {
	var req *request.GroupCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}
	err := req.Validate()
	if err != nil {
		con.Error(c, err.Error())
		return
	}

	if err := con.groupService.Create(req); err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c)
}

// Update 更新群组
// @Summary 更新群组信息
// @Description 更新群组的名称、描述等信息
// @Tags Group
// @Accept json
// @Produce json
// @Param group body model.GroupUpdateRequest true "更新群组信息"  // 更新请求数据
// @Success 200 {object} model.Response "群组更新成功"  // 更新成功响应
// @Failure 400 {object} model.Response "请求参数错误"  // 参数错误
// @Failure 500 {object} model.Response "内部服务器错误"  // 服务器错误
// @Router /group/update [post]
func (con GroupController) Update(c *gin.Context) {
	var req *request.GroupUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}
	if err := con.groupService.Update(req); err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c)
}

// Join 加入群组
// @Summary 加入群组
// @Description 加入群组
// @Tags Group
// @Accept json
// @Produce json
// @Param group_id query uint true "群组ID"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /group/join [Get]
func (con GroupController) Join(c *gin.Context) {
	groupIdStr, _ := c.GetQuery("group_id")
	groupId, _ := strconv.ParseUint(groupIdStr, 10, 64)
	userId, ok := c.Get("id")
	if !ok {
		con.Error(c, "need user_id")
		return
	}

	if err := con.groupService.Join(uint(groupId), userId.(uint)); err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c)
}

// Quit 退出群组
// @Summary 退出群组
// @Description 退出群组
// @Tags Group
// @Accept json
// @Produce json
// @Param group_id query uint true "群组ID"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /group/quit [Get]
func (con GroupController) Quit(c *gin.Context) {
	groupIdStr, _ := c.GetQuery("group_id")
	groupId, _ := strconv.ParseUint(groupIdStr, 10, 64)
	userId, ok := c.Get("id")
	if !ok {
		con.Error(c, "need user_id")
		return
	}
	if err := con.groupService.Quit(uint(groupId), userId.(uint)); err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c)
}

// Search 搜索群组
// @Summary 搜索群组
// @Description 搜索群组并返回搜索结果
// @Tags Group
// @Accept json
// @Produce json
// @Param data body model.GroupSearchRequest true "搜索群组请求参数"  // 搜索请求参数
// @Success 200 {object} model.Response{data=[]model.GroupVo} "群组搜索结果"  // 返回搜索到的群组列表
// @Failure 400 {object} model.Response "请求参数错误"  // 参数错误
// @Failure 500 {object} model.Response "内部服务器错误"  // 服务器错误
// @Router /group/search [post]
func (con GroupController) Search(c *gin.Context) {
	var req request.GroupSearchRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		con.Error(c, err.Error())
		return
	}

	resp, err := con.groupService.Search(req)
	if err != nil {
		con.Error(c, err.Error())
		return
	}

	con.Success(c, resp)
}

// Member 获取群组成员
// @Summary 获取群组成员
// @Description 获取指定群组的所有成员
// @Tags Group
// @Accept json
// @Produce json
// @Param group_id query uint true "群组ID"  // 群组ID
// @Success 200 {object} model.Response{data=[]model.MemberVo} "群组成员列表"  // 返回群组成员列表
// @Failure 400 {object} model.Response "请求参数错误"  // 参数错误
// @Failure 500 {object} model.Response "内部服务器错误"  // 服务器错误
// @Router /group/member [get]
func (con GroupController) Member(c *gin.Context) {
	groupIdStr := c.Query("group_id")
	groupId, _ := strconv.ParseUint(groupIdStr, 10, 64)
	members, err := con.groupService.Member(uint(groupId))
	if err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c, members)
}

// CreateAnnouncement 创建群组公告
// @Summary 创建群组公告
// @Description 创建一个新的群组公告
// @Tags Group
// @Accept json
// @Produce json
// @Param group_id path int true "群组ID"
// @Param data body request.GroupAnnouncementCreateRequest true "公告内容"
// @Success 200 {object} model.Response
// @Router /group/{group_id}/announcement/create [post]
func (con GroupController) CreateAnnouncement(c *gin.Context) {
	groupIdStr := c.Param("group_id")
	groupId, _ := strconv.Atoi(groupIdStr)

	var req request.GroupAnnouncementCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}

	if err := con.groupService.CreateAnnouncement(uint(groupId), &req); err != nil {
		con.Error(c, err.Error())
		return
	}

	con.Success(c)
}

// UpdateAnnouncement 更新群组公告
// @Summary 更新群组公告
// @Description 更新指定的群组公告
// @Tags Group
// @Accept json
// @Produce json
// @Param group_id path int true "群组ID"
// @Param data body request.GroupAnnouncementUpdateRequest true "公告更新内容"
// @Success 200 {object} model.Response
// @Router /group/{group_id}/announcement/update [post]
func (con GroupController) UpdateAnnouncement(c *gin.Context) {
	groupIdStr := c.Param("group_id")
	groupId, _ := strconv.Atoi(groupIdStr)

	var req request.GroupAnnouncementUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}

	if err := con.groupService.UpdateAnnouncement(uint(groupId), &req); err != nil {
		con.Error(c, err.Error())
		return
	}

	con.Success(c)
}

// DeleteAnnouncement 删除群组公告
// @Summary 删除群组公告
// @Description 删除指定群组的公告
// @Tags Group
// @Accept json
// @Produce json
// @Param group_id path int true "群组ID"
// @Param announcement_id query int true "公告ID"
// @Success 200 {object} model.Response
// @Router /group/{group_id}/announcement/delete [get]
func (con GroupController) DeleteAnnouncement(c *gin.Context) {
	groupIdStr := c.Param("group_id")
	groupId, _ := strconv.Atoi(groupIdStr)

	announcementIdStr := c.DefaultQuery("announcement_id", "")
	announcementId, _ := strconv.Atoi(announcementIdStr)

	if err := con.groupService.DeleteAnnouncement(uint(groupId), uint(announcementId)); err != nil {
		con.Error(c, err.Error())
		return
	}

	con.Success(c)
}

// GetAnnouncement 获取群组公告
// @Summary 查询群组公告
// @Description 获取指定群组的公告内容
// @Tags Group
// @Produce json
// @Param group_id path int true "群组ID"
// @Success 200 {object} model.Response
// @Router /group/{group_id}/announcement [get]
func (con GroupController) GetAnnouncement(c *gin.Context) {
	groupIdStr := c.Param("group_id")
	groupId, _ := strconv.Atoi(groupIdStr)

	announcement, err := con.groupService.GetAnnouncement(uint(groupId))
	if err != nil {
		con.Error(c, err.Error())
		return
	}

	con.Success(c, announcement)
}

// GetAnnouncementList 获取群组公告列表
// @Summary 获取群组公告列表
// @Description 获取指定群组的所有公告
// @Tags Group
// @Produce json
// @Param group_id path int true "群组ID"
// @Success 200 {object} model.Response
// @Router /group/{group_id}/announcement_list [get]
func (con GroupController) GetAnnouncementList(c *gin.Context) {
	groupIdStr := c.Param("group_id")
	groupId, _ := strconv.Atoi(groupIdStr)

	announcements, err := con.groupService.GetAnnouncementList(uint(groupId))
	if err != nil {
		con.Error(c, err.Error())
		return
	}

	con.Success(c, announcements)
}

// KickMember 踢出群成员
// @Summary 踢出群成员
// @Description 从群组中踢出指定成员
// @Tags Group
// @Accept json
// @Produce json
// @Param group_id path int true "群组ID"  // 群组ID
// @Param data body model.KickMemberRequest true "踢出成员请求数据"  // 踢出成员的数据（成员ID）
// @Success 200 {object} model.Response "成员成功被踢出"  // 成员踢出成功
// @Failure 400 {object} model.Response "请求参数错误"  // 参数错误
// @Failure 500 {object} model.Response "内部服务器错误"  // 服务器错误
// @Router /group/{group_id}/kick_member [post]
func (con GroupController) KickMember(c *gin.Context) {
	var req request.KickMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}

	userId := c.GetUint("id")
	groupIdStr := c.Param("group_id")
	groupId, _ := strconv.ParseUint(groupIdStr, 10, 64)
	if err := con.groupService.KickMember(userId, uint(groupId), req.MemberID); err != nil {
		con.Error(c, err.Error())
		return
	}

	con.Success(c)
}

// SetAdmin 设置群组管理员
// @Summary 设置群组管理员
// @Description 设置指定成员为群组管理员
// @Tags Group
// @Accept json
// @Produce json
// @Param group_id path int true "群组ID"  // 群组ID
// @Param data body model.SetAdminRequest true "设置管理员请求数据"  // 设置管理员的数据（成员ID）
// @Success 200 {object} model.Response "管理员设置成功"  // 管理员设置成功
// @Failure 400 {object} model.Response "请求参数错误"  // 参数错误
// @Failure 500 {object} model.Response "内部服务器错误"  // 服务器错误
// @Router /group/{group_id}/set_admin [post]
func (con GroupController) SetAdmin(c *gin.Context) {
	var req request.SetAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}

	userId := c.GetUint("id")
	groupIdStr := c.Param("group_id")
	groupId, _ := strconv.ParseUint(groupIdStr, 10, 64)
	if err := con.groupService.SetAdmin(userId, uint(groupId), req.MemberID); err != nil {
		con.Error(c, err.Error())
		return
	}

	con.Success(c)
}

// UnsetAdmin 取消群组管理员
// @Summary 取消群组管理员
// @Description 取消指定成员的群组管理员身份
// @Tags Group
// @Accept json
// @Produce json
// @Param group_id path int true "群组ID"  // 群组ID
// @Param data body model.UnsetAdminRequest true "取消管理员请求数据"  // 取消管理员的数据（成员ID）
// @Success 200 {object} model.Response "管理员取消成功"  // 管理员取消成功
// @Failure 400 {object} model.Response "请求参数错误"  // 参数错误
// @Failure 500 {object} model.Response "内部服务器错误"  // 服务器错误
// @Router /group/{group_id}/unset_admin [post]
func (con GroupController) UnsetAdmin(c *gin.Context) {
	var req request.UnsetAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}

	operatorId := c.GetUint("id")
	groupIdStr := c.Param("group_id")
	groupId, _ := strconv.ParseUint(groupIdStr, 10, 64)
	if err := con.groupService.UnsetAdmin(operatorId, uint(groupId), req.MemberID); err != nil {
		con.Error(c, err.Error())
		return
	}

	con.Success(c)
}

// MuteMember 禁言群成员
// @Summary 禁言群成员
// @Description 将指定成员禁言指定时间
// @Tags Group
// @Accept json
// @Produce json
// @Param group_id path int true "群组ID"  // 群组ID
// @Param data body model.GroupMemberMuteRequest true "禁言请求数据"  // 禁言请求的数据（成员ID和时长）
// @Success 200 {object} model.Response "成员禁言成功"  // 禁言成功
// @Failure 400 {object} model.Response "请求参数错误"  // 参数错误
// @Failure 500 {object} model.Response "内部服务器错误"  // 服务器错误
// @Router /group/{group_id}/mute [post]
func (con GroupController) MuteMember(c *gin.Context) {
	var req request.GroupMemberMuteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}

	operatorId := c.GetUint("id")
	groupIdStr := c.Param("group_id")
	groupId, _ := strconv.ParseUint(groupIdStr, 10, 64)

	if err := con.groupService.MuteMember(operatorId, uint(groupId), req.MemberID, req.Duration); err != nil {
		con.Error(c, err.Error())
		return
	}

	con.Success(c)
}

// UnmuteMember 解除禁言
// @Summary 解除禁言
// @Description 解除指定成员的禁言
// @Tags Group
// @Accept json
// @Produce json
// @Param group_id path int true "群组ID"  // 群组ID
// @Param data body model.UnmuteMemberRequest true "解除禁言请求数据"  // 解除禁言请求的数据（成员ID）
// @Success 200 {object} model.Response "成员解除禁言成功"  // 解除禁言成功
// @Failure 400 {object} model.Response "请求参数错误"  // 参数错误
// @Failure 500 {object} model.Response "内部服务器错误"  // 服务器错误
// @Router /group/{group_id}/unmute [post]
func (con GroupController) UnmuteMember(c *gin.Context) {
	var req request.UnmuteMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}
	userId := c.GetUint("id")
	groupIdStr := c.Param("group_id")
	groupId, _ := strconv.Atoi(groupIdStr)

	if err := con.groupService.UnmuteMember(userId, uint(groupId), req.MemberID); err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c)
}

// Dissolve 解散群组
// @Summary 解散群组
// @Description 解散指定的群组，群主或管理员可以执行此操作
// @Tags Group
// @Accept json
// @Produce json
// @Param group_id path int true "群组ID"  // 群组ID
// @Success 200 {object} model.Response "群组解散成功"  // 群组解散成功
// @Failure 400 {object} model.Response "请求参数错误"  // 参数错误
// @Failure 500 {object} model.Response "内部服务器错误"  // 服务器错误
// @Router /group/{group_id}/dissolve [post]
func (con GroupController) Dissolve(c *gin.Context) {
	groupIdStr := c.Param("group_id")
	groupId, _ := strconv.Atoi(groupIdStr)
	userId := c.GetUint("id")
	if err := con.groupService.Dissolve(userId, uint(groupId)); err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c)
}

// Transfer 转让群组所有权
// @Summary 转让群组所有权
// @Description 将群组所有权转让给其他成员
// @Tags Group
// @Accept json
// @Produce json
// @Param group_id path int true "群组ID"  // 群组ID
// @Param data body request.GroupTransferRequest true "转让请求数据"  // 转让请求的数据（新群主ID）
// @Success 200 {object} model.Response "群组所有权转让成功"  // 转让成功
// @Failure 400 {object} model.Response "请求参数错误"  // 参数错误
// @Failure 500 {object} model.Response "内部服务器错误"  // 服务器错误
// @Router /group/{group_id}/transfer [post]
func (con GroupController) Transfer(c *gin.Context) {
	var req request.GroupTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}
	userId := c.GetUint("id")
	err := con.groupService.TransferOwnership(userId, req)
	if err != nil {
		con.Error(c, err.Error())
		return
	}

	con.Success(c)
}

// Mute 群组禁言
// @Summary 群组禁言
// @Description 将群组成员禁言指定时间
// @Tags Group
// @Accept json
// @Produce json
// @Param data body model.GroupMuteRequest true "群组禁言请求数据"  // 群组禁言请求的数据（禁言时长等）
// @Success 200 {object} model.Response "群组禁言成功"  // 群组禁言成功
// @Failure 400 {object} model.Response "请求参数错误"  // 参数错误
// @Failure 500 {object} model.Response "内部服务器错误"  // 服务器错误
// @Router /group/mute [post]
func (con GroupController) Mute(c *gin.Context) {
	var req request.GroupMuteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}
	userId := c.GetUint("id")
	if err := con.groupService.Mute(userId, req); err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c)
}

// Limit 群组限流
// @Summary 群组限流
// @Description 群组限制发送消息速度
// @Tags Group
// @Accept json
// @Produce json
// @Param data body model.GroupLimitRequest true "群组限流请求数据"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response "请求参数错误"  // 参数错误
// @Failure 500 {object} model.Response "内部服务器错误"  // 服务器错误
// @Router /group/limit [post]
func (con GroupController) Limit(c *gin.Context) {
	var req request.GroupLimitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}
	userId := c.GetUint("id")
	if err := con.groupService.Limit(userId, req); err != nil {
		con.Error(c, err.Error())
		return
	}
}
