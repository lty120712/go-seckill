package logUtil

import (
	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"
)

// getCallerInfo 获取调用者的信息
func getCallerInfo(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "未知位置"
	}
	fn := runtime.FuncForPC(pc)
	return fmt.Sprintf("%s:%d [%s]", file, line, fn.Name())
}

// Errorf 打印错误日志（带调用信息）
func Errorf(format string, args ...interface{}) {
	logrus.Errorf("[%s] %s", getCallerInfo(2), fmt.Sprintf(format, args...))
}

// Warnf 打印警告日志（带调用信息）
func Warnf(format string, args ...interface{}) {
	logrus.Warnf("[%s] %s", getCallerInfo(2), fmt.Sprintf(format, args...))
}

// Infof 打印信息日志（带调用信息）
func Infof(format string, args ...interface{}) {
	logrus.Infof("[%s] %s", getCallerInfo(2), fmt.Sprintf(format, args...))
}

// Debugf 打印调试日志（带调用信息）
func Debugf(format string, args ...interface{}) {
	logrus.Debugf("[%s] %s", getCallerInfo(2), fmt.Sprintf(format, args...))
}
