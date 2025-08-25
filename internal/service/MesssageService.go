package service

import (
	"errors"
	"fmt"
	interfacerepository "go-chat/internal/interfaces/repository"
	"go-chat/internal/model"
	request "go-chat/internal/model/request"
	response "go-chat/internal/model/response"
	"go-chat/internal/utils"
	"sync"
)

type MessageService struct {
	messageRepository     interfacerepository.MessageRepositoryInterface
	userRepository        interfacerepository.UserRepositoryInterface
	groupMemberRepository interfacerepository.GroupMemberRepositoryInterface
}

var (
	MessageServiceInstance *MessageService
	messageOnce            sync.Once
)

func InitMessageService(messageRepository interfacerepository.MessageRepositoryInterface,
	userRepository interfacerepository.UserRepositoryInterface,
	groupMemberRepository interfacerepository.GroupMemberRepositoryInterface) {
	messageOnce.Do(func() {
		MessageServiceInstance = &MessageService{
			messageRepository:     messageRepository,
			userRepository:        userRepository,
			groupMemberRepository: groupMemberRepository,
		}
	})
}

// SendMessage 发送消息（支持私聊和群聊）
// msg 是已经构造好的 message 对象（建议外部构建 content 等）
func (s *MessageService) SendMessage(msg *model.Message) (*response.MessageVo, error) {
	if msg == nil {
		return nil, errors.New("消息不能为空")
	}
	if msg.SenderId == 0 {
		return nil, errors.New("发送者 Id 不能为空")
	}
	if *msg.TargetType == model.PrivateTarget && (msg.ReceiverId == nil) {
		return nil, errors.New("私聊消息必须有接收者 Id")
	}
	if *msg.TargetType == model.GroupTarget && (msg.GroupId == nil) {
		return nil, errors.New("群聊消息必须有群组 Id")
	}
	if msg.Type == nil {
		return nil, errors.New("消息类型不能为空")
	}
	if len(*msg.Content) == 0 {
		return nil, errors.New("消息内容不能为空")
	}
	if err := s.messageRepository.Save(msg); err != nil {
		return nil, err
	}
	vo, err := s.GetMessageById(msg.ID)
	if err != nil {
		return nil, err
	}
	return vo, nil
}

// GetMessageById  获取消息
func (s *MessageService) GetMessageById(id uint) (*response.MessageVo, error) {
	message, err := s.messageRepository.GetById(id)
	if err != nil {
		return nil, err
	}
	var messageVo = &response.MessageVo{}
	messageVo.GetFieldsFromMessage(message)
	//获取发送者信息
	sender, _ := s.userRepository.GetById(uint(message.SenderId))
	messageVo.SenderNickName = new(string)
	messageVo.SenderNickName = sender.Nickname
	messageVo.SenderAvatar = new(string)
	messageVo.SenderAvatar = sender.Avatar
	messageVo.SenderOnlineStatus = new(model.OnlineStatus)
	*messageVo.SenderOnlineStatus = sender.OnlineStatus
	return messageVo, nil
}

func (s *MessageService) ReadMessage(messageId uint, userId uint) error {
	//1.消息是否存在
	message, err := s.messageRepository.GetById(messageId)
	if err != nil {
		return err
	}
	if message == nil {
		return errors.New("消息不存在")
	}
	//2.将自己插入
	if message.ReaderIdList == nil {
		// 初始化为空切片
		message.ReaderIdList = &model.ReaderIdList{}
	}
	if !utils.Contains(*message.ReaderIdList, userId) {
		*message.ReaderIdList = append(*message.ReaderIdList, userId)
	}
	//3.更新数据库
	updateFields := map[string]interface{}{
		"reader_id_list": message.ReaderIdList,
	}
	err = s.messageRepository.UpdateFields(messageId, updateFields)
	if err != nil {
		return err
	}
	return nil
}

func (s *MessageService) QueryMessages(userId uint, req *request.QueryMessagesRequest) (*response.QueryMessagesResponse, error) {
	// 查询历史消息
	messages, err := s.messageRepository.QueryHistoryMessages(userId, req)
	if err != nil {
		return nil, err
	}

	// 是否有更多消息
	hasMore := false
	if len(messages) > req.Limit {
		hasMore = true
		messages = messages[:req.Limit]
	}

	// 获取发送者ID列表
	senderIds := make([]uint, len(messages))
	for i, msg := range messages {
		senderIds[i] = uint(msg.SenderId)
	}

	// 获取发送者信息
	var idToUserMap = make(map[uint]*model.User)
	userList, _ := s.userRepository.GetByIdList(senderIds)
	for _, user := range userList {
		user1 := user
		idToUserMap[user.ID] = &user1
	}

	// 获取群组成员
	if *req.TargetType == model.GroupTarget {
		memberList, _ := s.groupMemberRepository.GetMemberListByGroupId(req.TargetId)
		if memberList != nil {
			for _, member := range memberList {
				if user, exists := idToUserMap[member.UserId]; exists {
					*user.Nickname = member.Nickname
				}
			}
		}
	}

	// 构建 MessageVo 列表
	var list []*response.MessageVo
	for _, msg := range messages {
		// 填充基本的消息信息
		sender := idToUserMap[uint(msg.SenderId)]
		messageVo := &response.MessageVo{
			ID:                 msg.ID,
			CreatedAt:          msg.CreatedAt,
			UpdatedAt:          msg.UpdatedAt,
			SenderId:           msg.SenderId,
			ReceiverId:         msg.ReceiverId,
			GroupId:            msg.GroupId,
			ReplyId:            msg.ReplyId,
			ReaderIdList:       msg.ReaderIdList,
			TargetType:         msg.TargetType,
			Content:            msg.Content,
			Type:               msg.Type,
			Status:             msg.Status,
			ExtraData:          msg.ExtraData,
			IsRead:             msg.ReaderIdList != nil && utils.Contains(*msg.ReaderIdList, userId),
			SenderNickName:     sender.Nickname,
			SenderAvatar:       sender.Avatar,
			SenderOnlineStatus: &sender.OnlineStatus,
		}

		// 如果有 ReplyId，查询并填充被引用的消息
		if msg.ReplyId != nil {
			replyMessage, err := s.messageRepository.GetById(uint(*msg.ReplyId))
			if err == nil && replyMessage != nil {
				messageVo.Reply = &response.MessageVo{
					ID:                 replyMessage.ID,
					CreatedAt:          replyMessage.CreatedAt,
					UpdatedAt:          replyMessage.UpdatedAt,
					SenderId:           replyMessage.SenderId,
					ReceiverId:         replyMessage.ReceiverId,
					GroupId:            replyMessage.GroupId,
					ReplyId:            replyMessage.ReplyId,
					ReaderIdList:       replyMessage.ReaderIdList,
					TargetType:         replyMessage.TargetType,
					Content:            replyMessage.Content,
					Type:               replyMessage.Type,
					Status:             replyMessage.Status,
					ExtraData:          replyMessage.ExtraData,
					SenderNickName:     sender.Nickname,
					SenderAvatar:       sender.Avatar,
					SenderOnlineStatus: &sender.OnlineStatus,
				}
			}
		}

		list = append(list, messageVo)
	}

	// 计算游标
	var cursor int64 = 0
	if len(list) > 0 {
		cursor = int64(list[len(list)-1].ID)
	}

	return &response.QueryMessagesResponse{
		List:    list,
		Cursor:  cursor,
		HasMore: hasMore,
	}, nil
}

func (s *MessageService) Revoke(userId uint, messageId uint) error {

	message, err := s.messageRepository.GetById(messageId)
	if err != nil {
		return fmt.Errorf("消息未找到: %w", err)
	}

	if *message.TargetType == model.PrivateTarget {
		if message.SenderId != int64(userId) {
			return fmt.Errorf("没有权限撤回此消息")
		}
	}

	if *message.TargetType == model.GroupTarget {

		isAdminOrOwner := s.groupMemberRepository.IsOwnerOrAdmin(userId, uint(*message.GroupId))
		if message.SenderId != int64(userId) && !isAdminOrOwner {
			return fmt.Errorf("没有权限撤回他人消息")
		}
	}

	fields := map[string]interface{}{
		"status": model.Disable,
	}

	err = s.messageRepository.UpdateFields(messageId, fields)
	if err != nil {
		return fmt.Errorf("撤回消息失败: %w", err)
	}

	return nil
}
