package interfacesservice

import (
	"github.com/lty120712/gorm-pagination/pagination"
	"go-chat/internal/model"
	request "go-chat/internal/model/request"
	response "go-chat/internal/model/response"
)

type GroupServiceInterface interface {
	// Create 创建群组
	Create(req *request.GroupCreateRequest) error
	Join(groupId uint, userId uint) error

	Quit(groupId uint, memberId uint) error

	Member(groupId uint) (memberList []response.MemberVo, err error)
	// CreateAnnouncement 创建群组公告
	CreateAnnouncement(groupId uint, req *request.GroupAnnouncementCreateRequest) error

	// UpdateAnnouncement 更新群组公告
	UpdateAnnouncement(groupId uint, req *request.GroupAnnouncementUpdateRequest) error

	// DeleteAnnouncement 删除群组公告
	DeleteAnnouncement(groupId uint, announcementId uint) error

	// GetAnnouncement 获取群组单个公告
	GetAnnouncement(groupId uint) (*model.GroupAnnouncement, error)

	// GetAnnouncementList 获取群组公告列表
	GetAnnouncementList(groupId uint) ([]model.GroupAnnouncement, error)
	// KickMember 踢人
	KickMember(operatorId, groupId, targetMemberId uint) error

	SetAdmin(operatorId, groupId, memberId uint) error

	UnsetAdmin(operatorId, groupId, targetMemberId uint) error

	MuteMember(operatorId, groupId, targetMemberId uint, duration int64) error
	UnmuteMember(operatorId, groupId, memberId uint) error
	Search(req request.GroupSearchRequest) (*pagination.PageResult[model.Group], error)
	Dissolve(userId uint, groupId uint) error
	TransferOwnership(userId uint, req request.GroupTransferRequest) error
	Mute(userId uint, req request.GroupMuteRequest) error
	Update(req *request.GroupUpdateRequest) error
	Limit(userID uint, req request.GroupLimitRequest) error
}
