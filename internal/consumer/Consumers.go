package consumer

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"go-chat/internal/db"
	"go-chat/internal/model"
	dto "go-chat/internal/model/dto"
	"go-chat/internal/repository"
	"go-chat/internal/service"
	"go-chat/internal/utils"
	"go-chat/internal/utils/redisUtil"
	wsClient "go-chat/internal/ws/client"
	wsMessage "go-chat/internal/ws/message"
	"net/http"
	"strconv"
	"time"
)

// HandleChatConsumer 处理 chat 队列的消息
func HandleChatConsumer(msg []byte) {
	var messageDto dto.MessageDto
	err := json.Unmarshal(msg, &messageDto)
	if err != nil {
		logrus.Printf("消息反序列化失败: %s", err)
		return
	}

	vo, err := service.MessageServiceInstance.SendMessage(&messageDto.Message)
	if err != nil {
		wsClient.WebSocketClient.SendMessageToOne(messageDto.SendId, &model.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	wsClient.WebSocketClient.SendMessageToOne(messageDto.SendId, &model.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data: &wsMessage.Message{
			SendId: messageDto.SendId,
			Type:   wsMessage.ChatAck,
			Data:   vo,
			Time:   time.Now(),
		},
	})
	if *messageDto.Message.TargetType == model.PrivateTarget {
		wsClient.WebSocketClient.SendMessageToOne(*messageDto.Message.ReceiverId, &model.Response{
			Code:    http.StatusOK,
			Message: "success",
			Data: &wsMessage.Message{
				SendId: messageDto.SendId,
				Type:   wsMessage.Chat,
				Data:   vo,
				Time:   time.Now(),
			},
		})
	} else if *messageDto.Message.TargetType == model.GroupTarget {
		//群组限流
		group, err := repository.GroupRepositoryInstance.GetByID(uint(*messageDto.Message.GroupId))
		if err != nil {
			wsClient.WebSocketClient.SendMessageToOne(messageDto.SendId, &model.Response{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}
		if group.LimitCount != nil && group.LimitInterval != nil && !utils.IsZero(group.LimitCount) && !utils.IsZero(group.LimitInterval) {
			if redisUtil.IsRateLimited(db.Redis, "group_"+strconv.FormatInt(*messageDto.Message.GroupId, 10), *group.LimitCount, *group.LimitInterval) {
				wsClient.WebSocketClient.SendMessageToOne(messageDto.SendId, &model.Response{
					Code:    http.StatusTooManyRequests,
					Message: "消息发送频繁,请稍后再试",
					Data:    nil,
				})
			}
		}
		memberList, err := service.GroupServiceInstance.Member(uint(*messageDto.Message.GroupId))
		if err != nil {
			wsClient.WebSocketClient.SendMessageToOne(messageDto.SendId, &model.Response{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}
		onlineUserIds := wsClient.WebSocketClient.GetOnlineUserIds()
		groupOnlineUserIds := make([]int64, 0)
		for _, member := range memberList {
			if utils.Contains(onlineUserIds, int64(member.UserId)) && member.OnlineStatus == model.Online {
				groupOnlineUserIds = append(groupOnlineUserIds, int64(member.UserId))
			}
		}
		wsClient.WebSocketClient.SendMessageToMultiple(groupOnlineUserIds, &model.Response{
			Code:    http.StatusOK,
			Message: "success",
			Data: &wsMessage.Message{
				SendId: messageDto.SendId,
				Type:   wsMessage.Chat,
				Data:   vo,
				Time:   time.Now(),
			},
		})
	}

}
