package repository

import (
	"errors"
	"go-chat/internal/db"
	"go-chat/internal/model"
	response "go-chat/internal/model/response"
	"gorm.io/gorm"
)

type FriendRequestRepository struct {
}

var (
	FriendRequestRepositoryInstance *FriendRequestRepository
)

func InitFriendRequestRepository() {
	FriendRequestRepositoryInstance = &FriendRequestRepository{}
}

func (r *FriendRequestRepository) GetFriendRequestsByUser(userId uint, tx ...*gorm.DB) ([]model.FriendRequest, error) {
	gormDB := db.GetGormDB(tx...)
	var requests []model.FriendRequest
	err := gormDB.Where("user_id = ?", userId).Find(&requests).Error
	return requests, err
}

func (r *FriendRequestRepository) GetSentFriendRequests(userId uint, tx ...*gorm.DB) ([]response.FriendRequestVo, error) {
	gormDB := db.GetGormDB(tx...)
	var sentRequests []response.FriendRequestVo
	err := gormDB.Table("friend_requests fr").
		Select("fr.id,fr.friend_id as target_id, u.nickname, u.avatar, fr.status, fr.created_at as time").
		Joins("LEFT JOIN users u ON u.id = fr.friend_id").
		Where("fr.user_id = ?", userId).
		Order("fr.status DESC").
		Scan(&sentRequests).Error
	for i := range sentRequests {
		sentRequests[i].Sent = 1
	}
	return sentRequests, err
}

func (r *FriendRequestRepository) GetReceivedFriendRequests(userId uint, tx ...*gorm.DB) ([]response.FriendRequestVo, error) {
	gormDB := db.GetGormDB(tx...)
	var receivedRequests []response.FriendRequestVo
	err := gormDB.Table("friend_requests fr").
		Select("fr.id,fr.user_id as target_id, u.nickname, u.avatar, fr.status, fr.created_at as time").
		Joins("LEFT JOIN users u ON u.id = fr.user_id").
		Where("fr.friend_id = ?", userId).
		Order("fr.status DESC").
		Scan(&receivedRequests).Error
	for i := range receivedRequests {
		receivedRequests[i].Sent = 0
	}
	return receivedRequests, err
}

func (r *FriendRequestRepository) GetById(id int64, tx ...*gorm.DB) (*model.FriendRequest, error) {
	gormDB := db.GetGormDB(tx...)
	var req model.FriendRequest
	err := gormDB.First(&req, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &req, nil
}

func (r *FriendRequestRepository) UpdateStatus(id int64, status model.Status, tx ...*gorm.DB) error {
	gormDB := db.GetGormDB(tx...)
	return gormDB.Model(&model.FriendRequest{}).Where("id = ?", id).Update("status", status).Error
}
