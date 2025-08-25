package service

import (
	"errors"
	"go-chat/internal/db"
	interfaces "go-chat/internal/interfaces/handler"
	interfacerepository "go-chat/internal/interfaces/repository"
	"go-chat/internal/model"
	request "go-chat/internal/model/request"
	response "go-chat/internal/model/response"
	"go-chat/internal/utils"
	"gorm.io/gorm"
	"sort"
	"strings"
)

type FriendService struct {
	friendRepository        interfacerepository.FriendRepositoryInterface
	friendRequestRepository interfacerepository.FriendRequestRepositoryInterface
	friendGroupRepository   interfacerepository.FriendGroupRepositoryInterface
	userRepository          interfacerepository.UserRepositoryInterface
	wsHandler               interfaces.WsHandlerInterface
}

var (
	FriendServiceInstance *FriendService
)

func InitFriendService(friendRepository interfacerepository.FriendRepositoryInterface,
	friendRequestRepository interfacerepository.FriendRequestRepositoryInterface,
	friendGroupRepository interfacerepository.FriendGroupRepositoryInterface,
	userRepository interfacerepository.UserRepositoryInterface,
	wsHandler interfaces.WsHandlerInterface) {

	FriendServiceInstance = &FriendService{
		friendRepository:        friendRepository,
		friendRequestRepository: friendRequestRepository,
		friendGroupRepository:   friendGroupRepository,
		userRepository:          userRepository,
		wsHandler:               wsHandler,
	}
}

func (s *FriendService) Add(userId uint, friendIdList []uint) error {
	if len(friendIdList) == 0 {
		return errors.New("friendIdList不能为空")
	}
	for _, fid := range friendIdList {
		if fid == userId {
			return errors.New("不能添加自己为好友")
		}
	}

	return db.Mysql.Transaction(func(tx *gorm.DB) error {
		existingFriends, err := s.friendRepository.GetFriendsByUserId(userId, tx)
		if err != nil {
			return err
		}
		friendMap := make(map[uint]struct{}, len(existingFriends))
		for _, f := range existingFriends {
			friendMap[f.FriendId] = struct{}{}
		}

		existingRequests, err := s.friendRequestRepository.GetFriendRequestsByUser(userId, tx)
		if err != nil {
			return err
		}
		requestMap := make(map[uint]struct{}, len(existingRequests))
		for _, r := range existingRequests {
			if r.Status == model.Todo {
				requestMap[r.FriendId] = struct{}{}
			}
		}

		var requests []model.FriendRequest
		for _, fid := range friendIdList {
			if _, exists := friendMap[fid]; exists {
				continue
			}
			if _, requested := requestMap[fid]; requested {
				continue
			}
			requests = append(requests, model.FriendRequest{
				UserId:   userId,
				FriendId: fid,
				Status:   model.Todo,
			})
		}

		if len(requests) == 0 {
			return errors.New("好友已存在或申请已发出，无需重复申请")
		}

		return s.friendRepository.CreateFriendRequests(requests, tx)
	})
}

func (s *FriendService) ListReq(id uint) ([]response.FriendRequestVo, error) {
	sent, err := s.friendRequestRepository.GetSentFriendRequests(id)
	if err != nil {
		return nil, err
	}

	received, err := s.friendRequestRepository.GetReceivedFriendRequests(id)
	if err != nil {
		return nil, err
	}

	allRequests := append(sent, received...)

	sort.Slice(allRequests, func(i, j int) bool {
		return allRequests[i].Status < allRequests[j].Status
	})
	return allRequests, nil
}

func (s *FriendService) HandleReq(requestId int64, status model.Status) error {
	return db.Mysql.Transaction(func(tx *gorm.DB) error {
		req, err := s.friendRequestRepository.GetById(requestId, tx)
		if err != nil {
			return err
		}
		if req == nil {
			return errors.New("好友申请不存在")
		}
		if req.Status != model.Todo {
			return errors.New("该好友申请已处理")
		}

		err = s.friendRequestRepository.UpdateStatus(requestId, status, tx)
		if err != nil {
			return err
		}

		if status == model.Accept {
			friends := []model.Friend{
				{UserId: req.UserId, FriendId: req.FriendId},
				{UserId: req.FriendId, FriendId: req.UserId},
			}
			err = s.friendRepository.BatchCreate(friends, tx)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *FriendService) Remove(userId uint, friendIdList []int64) error {
	if len(friendIdList) == 0 {
		return errors.New("friendIdList不能为空")
	}

	for _, fid := range friendIdList {
		if fid == int64(userId) {
			return errors.New("不能删除自己")
		}
	}

	uintIDs := make([]uint, 0, len(friendIdList))
	for _, fid := range friendIdList {
		uintIDs = append(uintIDs, uint(fid))
	}

	return db.Mysql.Transaction(func(tx *gorm.DB) error {
		// 1. 删除 userId 对好友的关系
		if err := s.friendRepository.BatchDelete(userId, uintIDs, tx); err != nil {
			return err
		}
		// 2. 删除好友对 userId 的反向关系
		if err := s.friendRepository.BatchDeleteManyInverse(userId, uintIDs, tx); err != nil {
			return err
		}

		// 3. 更新 userId 的好友分组，删除这些好友
		groups, err := s.friendGroupRepository.GetGroupsByUserId(userId, tx)
		if err != nil {
			return err
		}
		for _, group := range groups {
			if group.FriendIdList == nil {
				continue
			}
			updatedList := model.FriendIdList{}
			removed := false
			for _, fid := range *group.FriendIdList {
				if utils.Contains(uintIDs, fid) {
					removed = true
					continue
				}
				updatedList = append(updatedList, fid)
			}
			if removed {
				group.FriendIdList = &updatedList
				if err := s.friendGroupRepository.UpdateGroup(&group, tx); err != nil {
					return err
				}
			}
		}

		// 4. 更新好友的分组，删除userId
		for _, fid := range uintIDs {
			friendGroups, err := s.friendGroupRepository.GetGroupsByUserId(fid, tx)
			if err != nil {
				return err
			}
			for _, group := range friendGroups {
				if group.FriendIdList == nil {
					continue
				}
				updatedList := model.FriendIdList{}
				removed := false
				for _, id := range *group.FriendIdList {
					if id == userId {
						removed = true
						continue
					}
					updatedList = append(updatedList, id)
				}
				if removed {
					group.FriendIdList = &updatedList
					if err := s.friendGroupRepository.UpdateGroup(&group, tx); err != nil {
						return err
					}
				}
			}
		}

		return nil
	})
}

func (s *FriendService) GroupCreate(userId uint, req request.FriendGroupCreateRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("分组名称不能为空")
	}

	var friendIds model.FriendIdList
	if req.FriendIdList != nil {
		friendIds = *req.FriendIdList
	} else {
		friendIds = model.FriendIdList{}
	}

	group := &model.FriendGroup{
		UserId:       userId,
		Name:         &req.Name,
		FriendIdList: &friendIds,
	}

	return s.friendGroupRepository.CreateGroup(group)
}

func (s *FriendService) GroupDelete(groupId int64) error {
	return s.friendGroupRepository.DeleteGroupById(uint(groupId))
}

func (s *FriendService) GroupUpdate(req request.FriendGroupUpdateRequest) error {
	// 1. 验证分组存在
	group, err := s.friendGroupRepository.GetGroupById(req.GroupId)
	if err != nil {
		return err
	}
	if group == nil {
		return errors.New("分组不存在")
	}

	// 2. 更新字段
	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		if trimmed == "" {
			return errors.New("分组名称不能为空")
		}
		group.Name = &trimmed
	}

	friendIds := req.FriendIdList
	group.FriendIdList = &friendIds

	// 3. 保存更新
	return s.friendGroupRepository.UpdateGroup(group)
}

func (s *FriendService) GroupList(userId uint) ([]response.GroupVo, error) {

	friends, err := s.friendRepository.GetFriendsWithUserInfo(userId)
	if err != nil {
		return nil, err
	}

	groups, err := s.friendGroupRepository.GetGroupsByUserId(userId)
	if err != nil {
		return nil, err
	}

	friendMap := make(map[uint]response.FriendVo)
	for _, f := range friends {
		friendMap[f.UserId] = f
	}

	groupedFriendIds := make(map[uint]struct{})
	var groupVos []response.GroupVo
	for _, g := range groups {
		var friendList response.FriendVoList
		if g.FriendIdList != nil {
			for _, fid := range *g.FriendIdList {
				if friend, ok := friendMap[fid]; ok {
					friendList = append(friendList, friend)
					groupedFriendIds[fid] = struct{}{}
				}
			}
		}

		groupVos = append(groupVos, response.GroupVo{
			GroupId:      g.ID,
			Name:         g.Name,
			FriendVoList: &friendList,
		})
	}

	var ungroupedFriends response.FriendVoList
	for fid, friend := range friendMap {
		if _, grouped := groupedFriendIds[fid]; !grouped {
			ungroupedFriends = append(ungroupedFriends, friend)
		}
	}

	if len(ungroupedFriends) > 0 {
		name := "未分组"
		groupVos = append(groupVos, response.GroupVo{
			GroupId:      0,
			Name:         &name,
			FriendVoList: &ungroupedFriends,
		})
	}

	return groupVos, nil
}
