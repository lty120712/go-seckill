package timer

import (
	"github.com/sirupsen/logrus"
	"go-chat/internal/service"
)

func HeartBeatTimer() {
	_, err := Timer.AddFunc("*/30 * * * * *", func() {
		err := service.UserServiceInstance.CheckOfflineUsers() // 你心跳检测的函数
		if err != nil {
			logrus.Errorf("心跳检测执行失败: %v", err)
		}
	})
	if err != nil {
		logrus.Error("定时任务(%v)添加失败:", "HeartBeatTimer", err)
		return
	}
}
