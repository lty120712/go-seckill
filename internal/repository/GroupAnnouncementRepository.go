package repository

import (
	"go-chat/internal/db"
	"go-chat/internal/model"
	"gorm.io/gorm"
	"sync"
)

type GroupAnnouncementRepository struct {
}

var (
	GroupAnnouncementRepositoryInstance *GroupAnnouncementRepository
	groupAnnouncementOnce               sync.Once
)

func InitGroupAnnouncementRepository() {
	groupAnnouncementOnce.Do(func() {
		GroupAnnouncementRepositoryInstance = &GroupAnnouncementRepository{}
	})
}

func (r *GroupAnnouncementRepository) Create(announcement *model.GroupAnnouncement, tx ...*gorm.DB) error {
	gormDB := db.GetGormDB(tx...)
	return gormDB.Create(announcement).Error
}

func (r *GroupAnnouncementRepository) Update(announcement *model.GroupAnnouncement, tx ...*gorm.DB) error {
	gormDB := db.GetGormDB(tx...)
	return gormDB.Save(announcement).Error
}

func (r *GroupAnnouncementRepository) Delete(announcementId uint, tx ...*gorm.DB) error {
	gormDB := db.GetGormDB(tx...)
	return gormDB.Delete(&model.GroupAnnouncement{}, announcementId).Error
}

func (r *GroupAnnouncementRepository) GetByID(announcementId uint, tx ...*gorm.DB) (*model.GroupAnnouncement, error) {
	gormDB := db.GetGormDB(tx...)
	var announcement model.GroupAnnouncement
	err := gormDB.Where("id = ?", announcementId).First(&announcement).Error
	if err != nil {
		return nil, err
	}
	return &announcement, nil
}

func (r *GroupAnnouncementRepository) GetLatestByGroupID(groupId uint, tx ...*gorm.DB) (*model.GroupAnnouncement, error) {
	gormDB := db.GetGormDB(tx...)
	var announcement model.GroupAnnouncement
	err := gormDB.Where("group_id = ?", groupId).
		Order("created_at DESC"). // 按创建时间倒序排列，获取最新公告
		First(&announcement).Error
	if err != nil {
		return nil, err
	}
	return &announcement, nil
}

func (r *GroupAnnouncementRepository) GetListByGroupID(groupId uint, tx ...*gorm.DB) ([]model.GroupAnnouncement, error) {
	gormDB := db.GetGormDB(tx...)
	var announcements []model.GroupAnnouncement
	err := gormDB.Where("group_id = ?", groupId).
		Order("created_at DESC"). // 按创建时间倒序排列
		Find(&announcements).Error
	if err != nil {
		return nil, err
	}
	return announcements, nil
}
