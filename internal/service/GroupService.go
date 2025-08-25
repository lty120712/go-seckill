package service

import (
	"errors"
	"fmt"
	"github.com/lty120712/gorm-pagination/pagination"
	"go-chat/internal/db"
	interfacerepository "go-chat/internal/interfaces/repository"
	"go-chat/internal/model"
	request "go-chat/internal/model/request"
	response "go-chat/internal/model/response"
	"go-chat/internal/repository"
	"go-chat/internal/utils/idUtil"
	"gorm.io/gorm"
	"sort"
	"sync"
	"time"
)

type GroupService struct {
	groupRepository             interfacerepository.GroupRepositoryInterface
	messageRepository           interfacerepository.MessageRepositoryInterface
	userRepository              interfacerepository.UserRepositoryInterface
	groupMemberRepository       interfacerepository.GroupMemberRepositoryInterface
	groupAnnouncementRepository interfacerepository.GroupAnnouncementRepositoryInterface
}

var (
	GroupServiceInstance *GroupService
	groupOnce            sync.Once
)

func InitGroupService(groupRepository interfacerepository.GroupRepositoryInterface,
	messageRepository interfacerepository.MessageRepositoryInterface,
	userRepository interfacerepository.UserRepositoryInterface,
	groupMemberRepository interfacerepository.GroupMemberRepositoryInterface,
	groupAnnouncementRepository interfacerepository.GroupAnnouncementRepositoryInterface,
) {
	groupOnce.Do(func() {
		GroupServiceInstance = &GroupService{
			groupRepository:             groupRepository,
			messageRepository:           messageRepository,
			userRepository:              userRepository,
			groupMemberRepository:       groupMemberRepository,
			groupAnnouncementRepository: groupAnnouncementRepository,
		}
	})
}

// Create 创建群组
func (s GroupService) Create(req *request.GroupCreateRequest) error {
	err := db.Mysql.Transaction(func(tx *gorm.DB) error {
		maxAttempts := 3
		code := idUtil.GenerateId()
		for attempts := 0; attempts < maxAttempts; attempts++ {
			if exists := s.groupRepository.ExistsByCode(code, tx); !exists {
				break
			}
			code = idUtil.GenerateId()
		}

		if exists := s.groupRepository.ExistsByCode(code, tx); exists {
			return errors.New("error code, please try again")
		}

		group := &model.Group{
			OwnerId: req.UserId,
			Name:    req.Name,
			MaxNum:  req.MaxNum,
			Code:    code,
			Avatar:  "https://th.bing.com/th/id/ODF.HcSOJGKbi4khTMFQYDxyIA?w=32&h=32&qlt=90&pcl=fffffa&o=6&pid=1.2",
			Desc:    "群主很懒,什么也没留下~",
			Status:  model.Enable,
		}

		if err := s.groupRepository.Save(group, tx); err != nil {
			return err
		}

		if req.MemberList != nil && len(*req.MemberList) > 0 {
			groupMemberList := make([]*model.GroupMember, len(*req.MemberList))

			idNickNameMap, err := s.userRepository.GetNickNamesByIds(*req.MemberList, tx)
			if err != nil {
				return err
			}

			for i, memberId := range *req.MemberList {
				nickname := idNickNameMap[memberId]
				var role model.Role
				if req.UserId == memberId {
					role = model.Owner // 群主角色
				} else {
					role = model.Member // 普通成员角色
				}
				groupMemberList[i] = &model.GroupMember{
					GroupId:   group.ID,
					MemberId:  memberId,
					GNickName: nickname,
					Role:      role,
				}
			}

			if err := s.groupMemberRepository.SaveBatch(groupMemberList, tx); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
func (s GroupService) Update(req *request.GroupUpdateRequest) error {
	// 检查是否存在群组
	_, err := s.groupRepository.GetByID(req.GroupId)
	if err != nil {
		return fmt.Errorf("群组不存在: %v", err)
	}
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Avatar != nil {
		updates["avatar"] = *req.Avatar
	}
	if req.Desc != nil {
		updates["desc"] = *req.Desc
	}
	if req.MaxNum != nil {
		updates["max_num"] = *req.MaxNum
	}

	return s.groupRepository.Update(req.GroupId, updates)
}
func (s GroupService) Join(groupId uint, userId uint) error {
	return db.Mysql.Transaction(func(tx *gorm.DB) error {
		exists := s.groupMemberRepository.ExistsByGroupIdAndUserId(groupId, userId, tx)
		if exists {
			return errors.New("用户已加入该群组")
		}
		rejoin := s.groupMemberRepository.RejoinGroupIfDeleted(groupId, userId, tx)
		if rejoin {
			return nil
		}
		nickname, err := s.userRepository.GetNickNamesById(userId, tx)
		if err != nil {
			return err
		}
		member := &model.GroupMember{
			GroupId:   groupId,
			MemberId:  userId,
			GNickName: nickname,
			Role:      model.Member,
		}
		if err := s.groupMemberRepository.Save(member, tx); err != nil {
			return err
		}
		return nil
	})
}

func (s GroupService) Quit(groupId uint, memberId uint) error {
	//群主不能退
	if s.groupMemberRepository.IsOwner(groupId, memberId) {
		return errors.New("owner can not quit")
	}
	return s.groupMemberRepository.DeleteByGroupIdAndUserId(groupId, memberId)
}

func (s GroupService) Member(groupId uint) (memberList []response.MemberVo, err error) {
	memberList, err = s.groupMemberRepository.GetMemberListByGroupId(groupId)
	if err != nil {
		return
	}
	sort.Slice(memberList, func(i, j int) bool {
		return memberList[i].Role > memberList[j].Role
	})
	return
}

func (s GroupService) Mute(userId uint, req request.GroupMuteRequest) error {
	return db.Mysql.Transaction(func(tx *gorm.DB) error {

		if !s.groupMemberRepository.IsOwner(req.GroupId, userId) {
			return errors.New("only owner can mute")
		}

		return s.groupRepository.Update(req.GroupId, map[string]interface{}{
			"mute_end": time.Unix(req.MuteEnd, 0),
		})
	})
}

// CreateAnnouncement 创建群组公告
func (s GroupService) CreateAnnouncement(groupId uint, req *request.GroupAnnouncementCreateRequest) error {

	// 插入公告数据
	announcement := &model.GroupAnnouncement{
		GroupID:   groupId,
		Content:   req.Content,
		Publisher: req.Publisher,
	}

	// 保存公告
	err := repository.GroupAnnouncementRepositoryInstance.Create(announcement)
	if err != nil {
		return fmt.Errorf("创建群组公告失败: %w", err)
	}

	return nil
}

// UpdateAnnouncement 更新群组公告
func (s GroupService) UpdateAnnouncement(groupId uint, req *request.GroupAnnouncementUpdateRequest) error {
	// 获取现有公告
	announcement, err := s.groupAnnouncementRepository.GetByID(uint(req.AnnouncementId))
	if err != nil {
		return fmt.Errorf("查询群组公告失败: %w", err)
	}
	if announcement.GroupID != groupId {
		return fmt.Errorf("群组ID不匹配")
	}

	// 更新公告内容
	announcement.Content = req.Content
	announcement.UpdatedAt = time.Now()

	// 保存更新后的公告
	err = repository.GroupAnnouncementRepositoryInstance.Update(announcement)
	if err != nil {
		return fmt.Errorf("更新群组公告失败: %w", err)
	}

	return nil
}

// DeleteAnnouncement 删除群组公告
func (s GroupService) DeleteAnnouncement(groupId uint, announcementId uint) error {
	// 获取现有公告
	announcement, err := s.groupAnnouncementRepository.GetByID(announcementId)
	if err != nil {
		return fmt.Errorf("查询群组公告失败: %w", err)
	}
	if announcement.GroupID != groupId {
		return fmt.Errorf("群组ID不匹配")
	}

	// 设置为删除状态
	err = repository.GroupAnnouncementRepositoryInstance.Delete(announcementId)
	if err != nil {
		return fmt.Errorf("删除群组公告失败: %w", err)
	}

	return nil
}

// GetAnnouncement 获取群组单个公告
func (s GroupService) GetAnnouncement(groupId uint) (*model.GroupAnnouncement, error) {
	// 查询群组的公告（可以按最新时间排序等）
	announcement, err := s.groupAnnouncementRepository.GetLatestByGroupID(groupId)
	if err != nil {
		return nil, fmt.Errorf("查询群组公告失败: %w", err)
	}

	return announcement, nil
}

// GetAnnouncementList 获取群组公告列表
func (s GroupService) GetAnnouncementList(groupId uint) ([]model.GroupAnnouncement, error) {
	// 查询群组公告列表
	announcements, err := s.groupAnnouncementRepository.GetListByGroupID(groupId)
	if err != nil {
		return nil, fmt.Errorf("获取群组公告列表失败: %w", err)
	}

	return announcements, nil
}

func (s GroupService) KickMember(operatorId, groupId, targetMemberId uint) error {
	// 获取操作者在该群的身份
	operator, err := s.groupMemberRepository.GetGroupMember(groupId, operatorId)
	if err != nil {
		return fmt.Errorf("获取操作人信息失败: %v", err)
	}
	if operator == nil {
		return errors.New("你不是该群成员")
	}
	if operator.Role != 1 && operator.Role != 2 {
		return errors.New("无权限踢人，仅群主或管理员可踢人")
	}

	// 群主只能被自己踢（禁止踢群主）
	target, err := s.groupMemberRepository.GetGroupMember(groupId, targetMemberId)
	if err != nil {
		return fmt.Errorf("获取目标成员失败: %v", err)
	}
	if target == nil {
		return errors.New("目标成员不存在")
	}
	if target.Role == 1 {
		return errors.New("不能踢群主")
	}
	if operator.Role == 2 && target.Role == 2 {
		return errors.New("管理员不能踢管理员")
	}

	// 删除成员记录（逻辑删除或物理删除皆可）
	if err := s.groupMemberRepository.RemoveMember(groupId, targetMemberId); err != nil {
		return fmt.Errorf("踢人失败: %v", err)
	}

	return nil
}

func (s GroupService) SetAdmin(operatorId, groupId, memberId uint) error {
	// 获取操作者信息
	operator, err := s.groupMemberRepository.GetGroupMember(groupId, operatorId)
	if err != nil {
		return fmt.Errorf("获取操作者信息失败: %v", err)
	}
	if operator == nil {
		return errors.New("操作者不是该群成员")
	}
	if operator.Role != 1 { // 只有群主有权限设置管理员
		return errors.New("只有群主才能设置管理员")
	}

	target, err := s.groupMemberRepository.GetGroupMember(groupId, memberId)
	if err != nil {
		return fmt.Errorf("获取目标成员失败: %v", err)
	}
	if target == nil {
		return errors.New("目标成员不存在")
	}
	if target.Role == model.Admin {
		return errors.New("群主不能被设置为管理员")
	}

	err = s.groupMemberRepository.Update(groupId, memberId, map[string]interface{}{
		"role": model.Admin,
	})
	if err != nil {
		return fmt.Errorf("设置管理员失败: %v", err)
	}

	return nil
}

func (s GroupService) UnsetAdmin(operatorId, groupId, targetMemberId uint) error {
	operator, err := s.groupMemberRepository.GetGroupMember(groupId, operatorId)
	if err != nil {
		return fmt.Errorf("获取操作人信息失败: %v", err)
	}
	if operator == nil || operator.Role != 1 {
		return errors.New("仅群主可以取消管理员权限")
	}

	target, err := s.groupMemberRepository.GetGroupMember(groupId, targetMemberId)
	if err != nil {
		return fmt.Errorf("获取目标成员失败: %v", err)
	}
	if target == nil {
		return errors.New("目标成员不存在")
	}
	if target.Role != model.Admin {
		return errors.New("该成员不是管理员")
	}

	return s.groupMemberRepository.Update(groupId, targetMemberId, map[string]interface{}{
		"role": model.Member,
	})
}

func (s GroupService) MuteMember(operatorId, groupId, targetMemberId uint, duration int64) error {
	operator, err := s.groupMemberRepository.GetGroupMember(groupId, operatorId)
	if err != nil {
		return fmt.Errorf("获取操作人信息失败: %v", err)
	}
	if operator == nil || (operator.Role != model.Owner && operator.Role != model.Admin) {
		return errors.New("仅群主或管理员可进行禁言操作")
	}

	target, err := s.groupMemberRepository.GetGroupMember(groupId, targetMemberId)
	if err != nil {
		return fmt.Errorf("获取目标成员失败: %v", err)
	}
	if target == nil {
		return errors.New("目标成员不存在")
	}
	if target.Role == 1 {
		return errors.New("不能禁言群主")
	}
	if operator.Role == 2 && target.Role == 2 {
		return errors.New("管理员不能禁言管理员")
	}

	muteUntil := time.Now().Add(time.Duration(duration) * time.Second)
	return s.groupMemberRepository.Update(groupId, targetMemberId, map[string]interface{}{
		"mute_end": muteUntil,
	})
}

func (s GroupService) UnmuteMember(operatorId, groupId, targetId uint) error {
	operator, err := s.groupMemberRepository.GetGroupMember(groupId, operatorId)
	if err != nil {
		return fmt.Errorf("获取操作人信息失败: %v", err)
	}
	if operator == nil {
		return errors.New("你不是该群成员")
	}
	if operator.Role != 1 && operator.Role != 2 {
		return errors.New("无权限操作，仅群主或管理员可解除禁言")
	}

	target, err := s.groupMemberRepository.GetGroupMember(groupId, targetId)
	if err != nil {
		return fmt.Errorf("获取目标成员失败: %v", err)
	}
	if target == nil {
		return errors.New("目标成员不存在")
	}
	if target.IsOwner() {
		return errors.New("无法解除群主禁言（群主默认不禁言）")
	}
	if operator.IsAdmin() && target.IsAdmin() {
		return errors.New("管理员不能解除管理员禁言")
	}

	return s.groupMemberRepository.Update(groupId, targetId, map[string]interface{}{
		"mute_end": nil,
	})
}

func (s GroupService) Search(req request.GroupSearchRequest) (*pagination.PageResult[model.Group], error) {

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}
	res, err := s.groupRepository.Page(req)
	return res, err
}
func (s GroupService) Dissolve(userId uint, groupId uint) error {

	group, err := s.groupRepository.GetByID(groupId)
	if err != nil {
		return fmt.Errorf("获取群组信息失败: %w", err)
	}
	if group == nil {
		return errors.New("群组不存在")
	}
	if group.OwnerId != userId {
		return errors.New("无权限：只有群主可以解散群组")
	}
	err = db.Mysql.Transaction(func(tx *gorm.DB) error {
		if err := s.groupRepository.Delete(groupId, tx); err != nil {
			return err
		}
		if err := s.groupMemberRepository.DeleteByGroupID(groupId, tx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("解散群组失败: %w", err)
	}
	//todo 发送通知

	return nil
}

func (s GroupService) TransferOwnership(userId uint, req request.GroupTransferRequest) error {
	return db.Mysql.Transaction(func(tx *gorm.DB) error {
		// 1. 校验群是否存在且 userId 是群主
		group, err := s.groupRepository.GetByID(req.GroupID, tx)
		if err != nil {
			return err
		}
		if group.OwnerId != userId {
			return errors.New("只有群主可以转移群主权限")
		}

		isMember := s.groupMemberRepository.ExistsByGroupIdAndUserId(req.GroupID, req.NewOwnerID, tx)
		if !isMember {
			return errors.New("新群主不是该群成员")
		}

		group.OwnerId = req.NewOwnerID
		if err := s.groupRepository.Save(group, tx); err != nil {
			return err
		}

		if err := s.groupMemberRepository.Update(req.GroupID, userId, map[string]interface{}{
			"role": model.Member,
		}, tx); err != nil {
			return err
		}
		if err := s.groupMemberRepository.Update(req.GroupID, req.NewOwnerID, map[string]interface{}{
			"role": model.Owner,
		}, tx); err != nil {
			return err
		}

		return nil
	})
}

func (s GroupService) Limit(userID uint, req request.GroupLimitRequest) error {
	return db.Mysql.Transaction(func(tx *gorm.DB) error {
		group, err := s.groupRepository.GetByID(req.GroupId, tx)
		if err != nil {
			return err
		}
		isOwnerOrAdmin := s.groupMemberRepository.IsOwnerOrAdmin(req.GroupId, userID)
		if !isOwnerOrAdmin {
			return errors.New("非群主或管理员")
		}
		group.LimitInterval = req.LimitInterval
		group.LimitCount = req.LimitCount
		return s.groupRepository.Save(group, tx)
	})
}
