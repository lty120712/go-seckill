// Package interfaces
package interfaces

import "go-chat/internal/model"

// WsHandlerInterface  接口
type WsHandlerInterface interface {
	ChatHandler(sendId int64, data interface{})
	HeartBeatHandler(sendId int64, data interface{})
	OnlineStatusNotice(sendId int64, data model.OnlineStatusNotice)
}
