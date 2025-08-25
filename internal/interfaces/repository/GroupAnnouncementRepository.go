package interfaces

import (
	"go-chat/internal/model"
	"gorm.io/gorm"
)

type GroupAnnouncementRepositoryInterface interface {
	// Create 创建群组公告
	Create(announcement *model.GroupAnnouncement, tx ...*gorm.DB) error

	// Update 更新群组公告
	Update(announcement *model.GroupAnnouncement, tx ...*gorm.DB) error

	// Delete 删除群组公告
	Delete(announcementId uint, tx ...*gorm.DB) error

	// GetByID 根据公告ID查询群组公告
	GetByID(announcementId uint, tx ...*gorm.DB) (*model.GroupAnnouncement, error)

	// GetLatestByGroupID 获取群组的最新公告
	GetLatestByGroupID(groupId uint, tx ...*gorm.DB) (*model.GroupAnnouncement, error)

	// GetListByGroupID 获取群组的所有公告
	GetListByGroupID(groupId uint, tx ...*gorm.DB) ([]model.GroupAnnouncement, error)
}
