package timer

import (
	"github.com/robfig/cron/v3"
)

var Timer = cron.New(cron.WithSeconds())

func InitTimer() {
	Timer.Start()
}

func heartBeatTimer() {
	//_, err := Timer.AddFunc("*/30 * * * * *", func() {
	//	err := service.UserServiceInstance.CheckOfflineUsers() // 你心跳检测的函数
	//	if err != nil {
	//		logrus.Errorf("心跳检测执行失败: %v", err)
	//	}
	//})
	//if err != nil {
	//	logrus.Error("定时任务(%v)添加失败:", "HeartBeatTimer", err)
	//	return
	//}
}
