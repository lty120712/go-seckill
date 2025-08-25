package repository

import (
	"errors"
	"go-chat/internal/db"
	"go-chat/internal/model"
	"gorm.io/gorm"
)

type FriendGroupRepository struct {
}

var (
	FriendGroupRepositoryInstance *FriendGroupRepository
)

func InitFriendGroupRepository() {
	FriendGroupRepositoryInstance = &FriendGroupRepository{}
}

func (r *FriendGroupRepository) CreateGroup(group *model.FriendGroup, tx ...*gorm.DB) error {
	gormDB := db.GetGormDB(tx...)
	return gormDB.Create(group).Error
}

func (r *FriendGroupRepository) GetGroupsByUserId(userId uint, tx ...*gorm.DB) ([]model.FriendGroup, error) {
	gormDB := db.GetGormDB(tx...)
	var groups []model.FriendGroup
	err := gormDB.Where("user_id = ?", userId).Find(&groups).Error
	return groups, err
}

func (r *FriendGroupRepository) UpdateGroup(group *model.FriendGroup, tx ...*gorm.DB) error {
	gormDB := db.GetGormDB(tx...)
	// 只更新FriendIdList和Name字段，避免其他字段被误修改
	return gormDB.Model(group).
		Select("friend_id_list", "name").
		Updates(group).
		Error
}

func (r *FriendGroupRepository) DeleteGroupById(groupId uint, tx ...*gorm.DB) error {
	gormDB := db.GetGormDB(tx...)
	return gormDB.Where("id = ?", groupId).Delete(&model.FriendGroup{}).Error
}

func (r *FriendGroupRepository) GetGroupById(groupId uint, tx ...*gorm.DB) (*model.FriendGroup, error) {
	gormDB := db.GetGormDB(tx...)
	var group model.FriendGroup
	err := gormDB.Where("id = ?", groupId).First(&group).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &group, err
}
