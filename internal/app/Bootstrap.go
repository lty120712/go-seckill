package app

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-chat/configs"
	_ "go-chat/docs"
	apiv1 "go-chat/internal/api/v1"
	"go-chat/internal/db"
	"go-chat/internal/manager"
	"go-chat/internal/timer"
	"time"
)

func Start() {
	//加载配置文件
	err := configs.LoadConfig()
	if err != nil {
		logrus.Error("Error loading config: %v", err)
		return
	}
	//配置数据库
	db.InitMysql()
	db.InitRedis()
	//配置logrus
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	//配置路由基本信息
	router := gin.Default()
	// 配置swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//配置Minio
	manager.InitMinIO()
	//配置rabbitmq
	manager.InitRabbitMQ()
	//配置WebSocket
	manager.InitWebSocket()
	//配置定时任务
	timer.InitTimer()
	// 配置 CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"}, // 允许的前端地址
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-User-ID"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,           // 是否允许带 Cookie
		MaxAge:           12 * time.Hour, // 预检请求缓存
	}))
	//配置依赖注入 要在倒数第二步
	doWire()
	//配置路由 要在最后一步
	apiv1.InitRouter(router)
	//启动服务
	router.Run(fmt.Sprintf(":%d", configs.AppConfig.Server.Port))
}
