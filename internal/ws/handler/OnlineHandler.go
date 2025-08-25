package wsHandler

import (
	"go-chat/internal/model"
	"go-chat/internal/repository"
	wsClient "go-chat/internal/ws/client"
	wsMessage "go-chat/internal/ws/message"
	"net/http"
	"time"
)

// OnlineStatusNotice 在线状态通知,我的在线状态改变时通知与我相关的朋友或群组
func (ws *WebSocketHandler) OnlineStatusNotice(sendId int64, onlineStatusNotice model.OnlineStatusNotice) {
	memberList, _ := repository.GroupMemberRepositoryInstance.GetRelatedMemberByUserId(uint(sendId))
	var userIdList []int64
	seen := make(map[int64]bool)
	userIdList = append(userIdList, sendId)
	for _, member := range memberList {
		userId := int64(member.UserId)
		if userId != sendId && !seen[userId] {
			seen[userId] = true
			userIdList = append(userIdList, userId)
		}
	}

	friendList, _ := repository.FriendRepositoryInstance.GetFriendsByUserId(uint(sendId))
	for _, friend := range friendList {
		friendId := int64(friend.FriendId)
		if friendId != sendId && !seen[friendId] {
			seen[friendId] = true
			userIdList = append(userIdList, friendId)
		}
	}

	wsClient.WebSocketClient.SendMessageToMultiple(userIdList,
		&model.Response{
			Code:    http.StatusOK,
			Message: "success",
			Data: &wsMessage.Message{
				Type:   wsMessage.OnlineStatus,
				SendId: sendId,
				Data:   onlineStatusNotice,
				Time:   time.Now(),
			},
		})
}
