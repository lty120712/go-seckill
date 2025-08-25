package timer

import (
	"github.com/robfig/cron/v3"
)

var Timer = cron.New(cron.WithSeconds())

func InitTimer() {
	HeartBeatTimer()
	Timer.Start()
}
