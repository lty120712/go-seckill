package repository

import (
	"errors"
	"go-chat/internal/db"
	"go-chat/internal/model"
	response "go-chat/internal/model/response"
	"gorm.io/gorm"
	"sync"
	"time"
)

type GroupMemberRepository struct {
}

var (
	GroupMemberRepositoryInstance *GroupMemberRepository
	groupMemberOnce               sync.Once
)

func InitGroupMemberRepository() {
	groupMemberOnce.Do(func() {
		GroupMemberRepositoryInstance = &GroupMemberRepository{}
	})
}

func (r *GroupMemberRepository) SaveBatch(list []*model.GroupMember, tx *gorm.DB) error {
	gormDB := db.GetGormDB(tx)
	return gormDB.CreateInBatches(list, len(list)/2).Error
}

func (r *GroupMemberRepository) Save(member *model.GroupMember, tx ...*gorm.DB) error {
	gormDB := db.GetGormDB(tx...)
	return gormDB.Create(member).Error
}

func (r *GroupMemberRepository) ExistsByGroupIdAndUserId(groupId uint, memberId uint, tx ...*gorm.DB) bool {
	gormDB := db.GetGormDB(tx...)
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = ? AND member_id = ? AND deleted_at is NULL)"
	err := gormDB.Raw(query, groupId, memberId).Scan(&exists).Error
	if err != nil {
		return false
	}
	return exists
}

func (r *GroupMemberRepository) RejoinGroupIfDeleted(groupId uint, memberId uint, tx ...*gorm.DB) bool {
	gormDB := db.GetGormDB(tx...)
	var groupMember model.GroupMember
	if err := gormDB.Unscoped().Where("group_id = ? AND member_id = ? AND deleted_at IS NOT NULL", groupId, memberId).First(&groupMember).Error; err == nil {
		groupMember.DeletedAt = gorm.DeletedAt{Time: time.Time{}, Valid: false}
		groupMember.CreatedAt = time.Now()
		groupMember.UpdatedAt = time.Now()
		if err := gormDB.Save(&groupMember).Error; err != nil {
			return false
		}
		return true
	}
	return false
}

func (r *GroupMemberRepository) DeleteByGroupIdAndUserId(groupId uint, memberId uint, tx ...*gorm.DB) error {
	gormDB := db.GetGormDB(tx...)
	return gormDB.Where("group_id = ? and member_id = ?", groupId, memberId).Delete(&model.GroupMember{}).Error
}

func (r *GroupMemberRepository) GetMemberListByGroupId(groupId uint, tx ...*gorm.DB) ([]response.MemberVo, error) {
	gormDB := db.GetGormDB(tx...)
	var memberList []response.MemberVo

	err := gormDB.Table("group_members as gm").
		Select("gm.group_id, gm.member_id AS user_id, IFNULL(gm.g_nick_name, u.nickname) AS nickname, gm.mute_end,gm.role,u.avatar,u.online_status").
		Joins("JOIN users u ON gm.member_id = u.id").
		Where("gm.group_id = ?", groupId).
		Order("u.online_status DESC").
		Scan(&memberList).Error

	if err != nil {
		return nil, err
	}
	return memberList, nil
}

func (r *GroupMemberRepository) IsOwner(groupId uint, memberId uint, tx ...*gorm.DB) bool {
	gormDB := db.GetGormDB(tx...)
	var groupMember model.GroupMember

	err := gormDB.Model(&model.GroupMember{}).
		Where("group_id = ? AND member_id = ? AND role = ?", groupId, memberId, model.Owner). // 检查是否为群主
		First(&groupMember).Error
	if err == nil {
		return true
	}
	return false
}

func (r *GroupMemberRepository) IsOwnerOrAdmin(groupId uint, memberId uint, tx ...*gorm.DB) bool {
	gormDB := db.GetGormDB(tx...)
	var groupMember model.GroupMember

	err := gormDB.Model(&model.GroupMember{}).
		Where("group_id = ? AND member_id = ? AND (role = ? OR role =?)", groupId, memberId, model.Owner, model.Admin). // 检查是否为群主
		First(&groupMember).Error
	if err == nil {
		return true
	}
	return false
}
func (r *GroupMemberRepository) GetRelatedMemberByUserId(id uint, tx ...*gorm.DB) (memberList []response.MemberVo, err error) {
	gormDB := db.GetGormDB(tx...)

	var groupIds []uint
	err = gormDB.Model(&model.GroupMember{}).
		Where("member_id = ?", id).
		Pluck("group_id", &groupIds).
		Error

	if err != nil {
		return nil, err
	}

	if len(groupIds) == 0 {
		return nil, nil
	}

	// 联表查询并去重
	err = gormDB.Model(&model.GroupMember{}).
		Select("group_members.member_id as user_id,group_members.group_id, group_members.g_nick_name as nickname, group_members.role, users.avatar, users.online_status").
		Joins("LEFT JOIN users ON users.id = group_members.member_id"). // 联接 User 表
		Where("group_members.group_id IN ? AND group_members.member_id != ? AND users.online_status = ?", groupIds, id, model.Online).
		Find(&memberList).
		Error

	if err != nil {
		return nil, err
	}

	return memberList, nil
}
func (r *GroupMemberRepository) GetGroupMember(groupId, userId uint, tx ...*gorm.DB) (*model.GroupMember, error) {
	gormDB := db.GetGormDB(tx...)
	var gm model.GroupMember
	err := gormDB.Where("group_id = ? AND member_id = ?", groupId, userId).First(&gm).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &gm, err
}

func (r *GroupMemberRepository) RemoveMember(groupId, userId uint, tx ...*gorm.DB) error {
	gormDB := db.GetGormDB(tx...)
	return gormDB.Where("group_id = ? AND member_id = ?", groupId, userId).Delete(&model.GroupMember{}).Error
}

func (r *GroupMemberRepository) DeleteByGroupID(groupID uint, tx ...*gorm.DB) error {
	gormDB := db.GetGormDB(tx...)
	return gormDB.Where("group_id = ?", groupID).Delete(&model.GroupMember{}).Error
}

func (r *GroupMemberRepository) Update(groupID, memberID uint, updates map[string]interface{}, tx ...*gorm.DB) error {
	if len(updates) == 0 {
		return nil
	}
	gormDB := db.GetGormDB(tx...)
	return gormDB.Model(&model.GroupMember{}).
		Where("group_id = ? AND member_id = ?", groupID, memberID).
		Updates(updates).Error
}
