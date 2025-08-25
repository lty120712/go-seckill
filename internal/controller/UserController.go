package controller

import (
	"github.com/gin-gonic/gin"
	interfacesservice "go-chat/internal/interfaces/service"
	"go-chat/internal/model"
	request "go-chat/internal/model/request"
	"go-chat/internal/service"
	"strconv"
)

// UserController 用户相关控制器
// @Tags User
// @Description 控制用户相关的 API，包含用户注册、登录和 ping 等接口
type UserController struct {
	BaseController
	userService interfacesservice.UserServiceInterface
}

var UserControllerInstance *UserController

func InitUserController(userService interfacesservice.UserServiceInterface) {
	UserControllerInstance = &UserController{
		userService: userService,
	}
}

// Register 用户注册接口
// @Summary 用户注册
// @Description 接口用于用户注册，提供注册所需的字段
// @Tags user
// @Accept json
// @Produce json
// @Param username body string true "用户名"
// @Param password body string true "密码"
// @Success 200 {object} model.Response "成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Router /user/register [post]
func (con UserController) Register(c *gin.Context) {
	registerRequest := &request.RegisterRequest{}
	if err := c.ShouldBindJSON(registerRequest); err != nil {
		con.Error(c, err.Error())
		return
	}
	if err := service.UserServiceInstance.Register(registerRequest.Username, registerRequest.Password, registerRequest.RePassword); err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c)
	return
}

// Login 用户登录接口
// @Summary 用户登录
// @Description 用户登录，返回用户的认证信息
// @Tags user
// @Accept json
// @Produce json
// @Param body body model.LoginRequest true "登录信息"
// @Success 200 {object} model.Response "成功"
// @Failure 401 {object} model.Response "登陆失败"
// @Router /user/login [post]
func (con UserController) Login(c *gin.Context) {
	loginRequest := &request.LoginRequest{}
	if err := c.ShouldBindJSON(loginRequest); err != nil {
		con.Error(c, err.Error())
		return
	}
	token, err := con.userService.Login(loginRequest.Username, loginRequest.Password)
	if err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c, token)
	return
}

// Logout 用户登出接口
// @Summary 用户登出
// @Description 用户登出，清除认证信息
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} model.Response "成功"  // 成功的登出响应
// @Failure 500 {object} model.Response "登出失败"  // 登出失败
// @Router /user/logout [get]
func (con UserController) Logout(c *gin.Context) {
	id := c.GetUint("id")
	con.userService.Logout(id)
	con.Success(c)
	return
}

// OnlineStatusChange 用户在线状态变更接口
// @Summary 用户在线状态变更
// @Description 修改用户的在线状态（例如：在线、离线、忙碌等）
// @Tags user
// @Accept json
// @Produce json
// @Param online_status query int true "在线状态"  // 在线状态，必填，0=离线，1=在线，2=忙碌
// @Success 200 {object} model.Response "成功"  // 在线状态修改成功
// @Failure 400 {object} model.Response "请求参数错误"  // 请求参数错误
// @Failure 500 {object} model.Response "在线状态变更失败"  // 状态变更失败
// @Router /user/online_status_change [get]
func (con UserController) OnlineStatusChange(c *gin.Context) {
	id := c.GetUint("id")
	onlineStatus, _ := strconv.ParseInt(c.Query("online_status"), 10, 64)
	if err := con.userService.OnlineStatusChange(id, model.OnlineStatus(onlineStatus)); err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c)
}

// GetUserInfo 获取用户信息接口
// @Summary 获取用户信息
// @Description 获取当前登录用户的信息
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "用户ID"  // ID，作为URL路径参数
// @Success 200 {object} model.Response "成功"  // 返回当前用户的信息
// @Failure 500 {object} model.Response "获取用户信息失败"  // 获取用户信息失败
// @Router /user/info [get]
func (con UserController) GetUserInfo(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	userInfo, err := con.userService.GetUserInfo(uint(id))
	if err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c, userInfo)
}

// Update 用户更新信息接口
// @Summary 更新用户信息
// @Description 更新用户的个人资料，如昵称、头像等
// @Tags user
// @Accept json
// @Produce json
// @security Bearer
// @Param body body  model.UserUpdateRequest true "用户信息"  // 用户信息，包含昵称、头像等字段
// @Success 200 {string} string "成功"  // 更新成功的响应
// @Failure 401 {string} string "未授权"  // 未授权，必须登录后更新
// @Failure 400 {string} string "请求参数错误"  // 请求参数错误
// @Failure 500 {string} string "更新失败"  // 更新失败
// @Router /user/update [post]
func (con UserController) Update(c *gin.Context) {
	updateRequest := &request.UserUpdateRequest{}
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		con.Error(c, err.Error())
		return
	}
	err := updateRequest.Validate()
	if err != nil {
		con.Error(c, err.Error())
		return
	}
	if err := con.userService.UpdateUser(updateRequest); err != nil {
		con.Error(c, err)
		return
	}
	con.Success(c)
}
