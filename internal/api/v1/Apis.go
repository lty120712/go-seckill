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

func FileApi(r *gin.Engine) {
	fileApi := r.Group(configs.AppConfig.Api.Prefix+"/file", middleware.AuthMiddleware())
	{
		fileApi.POST("/upload", controllers.FileControllerInstance.Upload)
	}
}
