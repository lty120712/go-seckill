package interfaces

import (
	"go-chat/internal/model"
	response "go-chat/internal/model/response"
	"gorm.io/gorm"
)

type GroupMemberRepositoryInterface interface {
	SaveBatch(list []*model.GroupMember, tx *gorm.DB) error

	Save(member *model.GroupMember, tx ...*gorm.DB) error

	ExistsByGroupIdAndUserId(groupId uint, memberId uint, tx ...*gorm.DB) bool

	RejoinGroupIfDeleted(groupId uint, memberId uint, tx ...*gorm.DB) bool

	DeleteByGroupIdAndUserId(groupId uint, memberId uint, tx ...*gorm.DB) error

	GetMemberListByGroupId(groupId uint, tx ...*gorm.DB) ([]response.MemberVo, error)

	IsOwner(groupId uint, memberId uint, tx ...*gorm.DB) bool
	IsOwnerOrAdmin(groupId uint, memberId uint, tx ...*gorm.DB) bool
	GetRelatedMemberByUserId(id uint, tx ...*gorm.DB) (memberList []response.MemberVo, err error)

	GetGroupMember(groupId, userId uint, tx ...*gorm.DB) (*model.GroupMember, error)
	RemoveMember(groupId, userId uint, tx ...*gorm.DB) error
	DeleteByGroupID(groupID uint, tx ...*gorm.DB) error
	Update(groupID, memberID uint, updates map[string]interface{}, tx ...*gorm.DB) error
}
