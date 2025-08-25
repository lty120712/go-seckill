package repository

import (
	"errors"
	"go-chat/internal/db"
	"go-chat/internal/model"
	request "go-chat/internal/model/request"
	"gorm.io/gorm"
	"sync"
)

type MessageRepository struct {
}

var (
	MessageRepositoryInstance *MessageRepository
	messageOnce               sync.Once
)

func InitMessageRepository() {
	messageOnce.Do(func() {
		MessageRepositoryInstance = &MessageRepository{}
	})
}

func (r *MessageRepository) Save(message *model.Message) (err error) {
	err = db.Mysql.Create(message).Error
	return
}

func (r *MessageRepository) GetById(id uint) (message *model.Message, err error) {
	err = db.Mysql.Where("id = ?", id).First(&message).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}

func (r *MessageRepository) UpdateFields(id uint, fields map[string]interface{}) (err error) {
	err = db.Mysql.Model(&model.Message{}).Where("id = ?", id).Updates(fields).Error
	return
}

func (r *MessageRepository) QueryHistoryMessages(userId uint, req *request.QueryMessagesRequest) ([]*model.Message, error) {
	tx := db.Mysql.Model(&model.Message{})
	switch *req.TargetType {
	case model.PrivateTarget:
		tx = tx.Where("target_type = ?", model.PrivateTarget).
			Where(
				tx.Where("sender_id = ? AND receiver_id = ?", userId, req.TargetId).
					Or("sender_id = ? AND receiver_id = ?", req.TargetId, userId),
			)
	case model.GroupTarget:
		tx = tx.Where("target_type = ?", model.GroupTarget).
			Where("group_id = ?", req.TargetId)
	default:
		return nil, errors.New("非法的 target_type")
	}

	if req.Cursor > 0 {
		tx = tx.Where("id < ?", req.Cursor)
	}

	if req.MessageTypes != nil {
		tx = tx.Where("type = ?", *req.MessageTypes)
	}
	if req.Keyword != nil {
		tx = tx.Where("JSON_EXTRACT(content, '$[*].text') LIKE ?", "%"+*req.Keyword+"%")
	}
	if !req.StartTime.IsZero() {
		tx = tx.Where("created_at >= ?", req.StartTime)
	}
	if !req.EndTime.IsZero() {
		tx = tx.Where("created_at <= ?", req.EndTime)
	}

	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	var messages []*model.Message
	err := tx.Order("id DESC").Limit(limit + 1).Find(&messages).Error
	return messages, err
}
