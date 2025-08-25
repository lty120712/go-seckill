package repository

import (
	"go-chat/internal/db"
	"go-chat/internal/model"
	response "go-chat/internal/model/response"
	"gorm.io/gorm"
)

type FriendRepository struct {
}

var (
	FriendRepositoryInstance *FriendRepository
)

func InitFriendRepository() {
	FriendRepositoryInstance = &FriendRepository{}
}

func (r *FriendRepository) GetFriendsByUserId(userId uint, tx ...*gorm.DB) ([]model.Friend, error) {
	gormDB := db.GetGormDB(tx...)
	var friends []model.Friend
	err := gormDB.Where("user_id = ?", userId).Find(&friends).Error
	return friends, err
}

func (r *FriendRepository) BatchCreate(friends []model.Friend, tx ...*gorm.DB) error {
	gormDB := db.GetGormDB(tx...)
	return gormDB.CreateInBatches(&friends, len(friends)).Error
}

func (r *FriendRepository) CreateFriendRequests(requests []model.FriendRequest, tx ...*gorm.DB) error {
	gormDB := db.GetGormDB(tx...)
	return gormDB.Create(&requests).Error
}

func (r *FriendRepository) BatchDelete(userId uint, friendIds []uint, tx ...*gorm.DB) error {
	gormDB := db.GetGormDB(tx...)

	return gormDB.
		Where("(user_id = ? AND friend_id IN (?)) OR (user_id IN (?) AND friend_id = ?)",
			userId, friendIds, friendIds, userId).
		Delete(&model.Friend{}).Error
}

func (r *FriendRepository) BatchDeleteManyInverse(userId uint, friendIds []uint, tx ...*gorm.DB) error {
	gormDB := db.GetGormDB(tx...)
	return gormDB.Where("user_id IN ? AND friend_id = ?", friendIds, userId).Delete(&model.Friend{}).Error
}
func (r *FriendRepository) GetFriendsWithUserInfo(userId uint, tx ...*gorm.DB) ([]response.FriendVo, error) {
	gormDB := db.GetGormDB(tx...)
	var friends []response.FriendVo
	err := gormDB.Table("friends f").
		Select(`f.friend_id as user_id, u.nickname, u.desc, u.online_status, u.device_info`).
		Joins("LEFT JOIN users u ON u.id = f.friend_id").
		Where("f.user_id = ?", userId).
		Find(&friends).Error
	return friends, err
}
