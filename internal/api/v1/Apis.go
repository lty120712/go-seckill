package apiv1

import (
	"github.com/gin-gonic/gin"
	"go-chat/configs"
	controllers "go-chat/internal/controller"
	"go-chat/internal/middleware"
)

func InitRouter(r *gin.Engine) {
	//配置路由中间件
	RegisterMiddlewares(r)
	//配置控制器的路由
	UserApi(r)
	MessageApi(r)
	GroupApi(r)
	FriendApi(r)
	FileApi(r)
}

func UserApi(r *gin.Engine) {
	userApi := r.Group(configs.AppConfig.Api.Prefix + "/user")
	{
		userApi.POST("/register", controllers.UserControllerInstance.Register)
		userApi.POST("/login", controllers.UserControllerInstance.Login)
		userApi.GET("/logout", middleware.AuthMiddleware(), controllers.UserControllerInstance.Logout)
		userApi.GET("/online_status_change", middleware.AuthMiddleware(), controllers.UserControllerInstance.OnlineStatusChange)
		userApi.GET("/info", controllers.UserControllerInstance.GetUserInfo)
		userApi.POST("/update", middleware.AuthMiddleware(), controllers.UserControllerInstance.Update)
	}
}

func MessageApi(r *gin.Engine) {
	messageApi := r.Group(configs.AppConfig.Api.Prefix+"/message", middleware.AuthMiddleware())
	{
		messageApi.POST("/read", controllers.MessageControllerInstance.Read)
		messageApi.POST("/query", controllers.MessageControllerInstance.Query)
		messageApi.GET("/:id/revoke", controllers.MessageControllerInstance.Revoke)
	}
}

func GroupApi(r *gin.Engine) {
	groupApi := r.Group(configs.AppConfig.Api.Prefix+"/group", middleware.AuthMiddleware())
	{
		// 群聊相关
		groupApi.POST("/create", controllers.GroupControllerInstance.Create)
		groupApi.POST("/update", controllers.GroupControllerInstance.Update)
		groupApi.GET("/join", controllers.GroupControllerInstance.Join)
		groupApi.GET("/quit", controllers.GroupControllerInstance.Quit)
		groupApi.POST("/search", controllers.GroupControllerInstance.Search)
		groupApi.GET("/member", controllers.GroupControllerInstance.Member)
		groupApi.POST("/mute", controllers.GroupControllerInstance.Mute)
		groupApi.POST("/limit", controllers.GroupControllerInstance.Limit)
		// 群公告相关
		groupApi.POST("/:group_id/announcement/create", controllers.GroupControllerInstance.CreateAnnouncement)
		groupApi.POST("/:group_id/announcement/update", controllers.GroupControllerInstance.UpdateAnnouncement)
		groupApi.GET("/:group_id/announcement/delete", controllers.GroupControllerInstance.DeleteAnnouncement)
		groupApi.GET("/:group_id/announcement", controllers.GroupControllerInstance.GetAnnouncement)
		groupApi.GET("/:group_id/announcement_list", controllers.GroupControllerInstance.GetAnnouncementList)

		// 群聊权限相关
		groupApi.POST("/:group_id/kick", controllers.GroupControllerInstance.KickMember)
		groupApi.POST("/:group_id/set_admin", controllers.GroupControllerInstance.SetAdmin)
		groupApi.POST("/:group_id/unset_admin", controllers.GroupControllerInstance.UnsetAdmin)
		groupApi.POST("/:group_id/mute", controllers.GroupControllerInstance.MuteMember)
		groupApi.POST("/:group_id/unmute", controllers.GroupControllerInstance.UnmuteMember)
		groupApi.POST("/:group_id/dissolve", controllers.GroupControllerInstance.Dissolve)
		groupApi.POST("/:group_id/transfer", controllers.GroupControllerInstance.Transfer)
	}
}

func FriendApi(r *gin.Engine) {
	friendApi := r.Group(configs.AppConfig.Api.Prefix+"/friend", middleware.AuthMiddleware())
	{
		friendApi.POST("/add", controllers.FriendControllerInstance.Add)                  //好友申请
		friendApi.GET("/list_req", controllers.FriendControllerInstance.ListReq)          //获取好友申请列表
		friendApi.POST("/handle_req", controllers.FriendControllerInstance.HandleReq)     //处理好友申请
		friendApi.POST("/remove", controllers.FriendControllerInstance.Remove)            //删除好友
		friendApi.POST("/group_create", controllers.FriendControllerInstance.GroupCreate) //创建好友分组
		friendApi.GET("/group_delete", controllers.FriendControllerInstance.GroupDelete)  //删除好友分组
		friendApi.POST("/group_update", controllers.FriendControllerInstance.GroupUpdate) //修改好友分组
		friendApi.GET("/group_list", controllers.FriendControllerInstance.GroupList)      //查询好友分组列表
	}
}

func FileApi(r *gin.Engine) {
	fileApi := r.Group(configs.AppConfig.Api.Prefix+"/file", middleware.AuthMiddleware())
	{
		fileApi.POST("/upload", controllers.FileControllerInstance.Upload)
	}
}
