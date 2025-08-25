package service

import (
	"errors"
	"fmt"
	"go-chat/internal/db"
	interfacehandler "go-chat/internal/interfaces/handler"
	interfacerepository "go-chat/internal/interfaces/repository"
	"go-chat/internal/model"
	request "go-chat/internal/model/request"
	response "go-chat/internal/model/response"
	"go-chat/internal/utils/jwtUtil"
	"go-chat/internal/utils/logUtil"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"sync"
	"time"
)

type UserService struct {
	userRepository interfacerepository.UserRepositoryInterface
	wsHandler      interfacehandler.WsHandlerInterface
}

var (
	UserServiceInstance *UserService
	once                sync.Once
)

func InitUserService(wsHandler interfacehandler.WsHandlerInterface, userRepository interfacerepository.UserRepositoryInterface) *UserService {
	once.Do(func() {
		UserServiceInstance = &UserService{
			wsHandler:      wsHandler,
			userRepository: userRepository,
		}
	})
	return UserServiceInstance
}

// Register 注册
func (u *UserService) Register(username, password, rePassword *string) (err error) {
	if *password != *rePassword {
		return errors.New("密码不一致")
	}
	user, err := u.userRepository.GetByName(username)
	if err != nil {
		logUtil.Errorf("GetUserByName error: %v", err)
		return err
	}
	if user != nil {
		logUtil.Warnf("用户(%v)已存在", username)
		return errors.New(fmt.Sprintf("用户(%v)已存在", username))
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user = &model.User{Username: *username,
		Password:     string(hashedPassword),
		Nickname:     username,
		Status:       model.Enable,
		OnlineStatus: model.Offline,
	}
	err = u.userRepository.Save(user)
	if err != nil {
		logUtil.Errorf("保存用户失败: %v", err)
		return err
	}
	return nil
}

func (u *UserService) Login(username, password *string) (token string, err error) {
	//根据username查询
	user, err := u.userRepository.GetByName(username)
	// todo 1.非封禁状态 之后维护redis 黑名单
	if user.Status == model.Disable {
		return "", errors.New("用户被封禁")
	}
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("用户名或密码错误")
	}
	//验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(*password)); err != nil {
		return "", errors.New("用户名或密码错误")
	}
	//jwt 返回
	token, err = jwtUtil.GenerateJWT(user.ID)
	if err != nil {
		return "", err
	}
	if user.OnlineStatus == model.Offline {
		user.OnlineStatus = model.Online
		fields := map[string]interface{}{
			"online_status":  user.OnlineStatus,
			"login_time":     time.Now().Unix(),
			"heartbeat_time": time.Now().Unix(),
		}
		err = u.userRepository.UpdateFields(user.ID, fields, db.Mysql)
		if err != nil {
			return "", err
		}
		onlineStatusNotice := model.OnlineStatusNotice{
			UserId:       user.ID,
			OnlineStatus: model.Online,
			ActionType:   model.LoginAction,
		}
		go u.wsHandler.OnlineStatusNotice(int64(user.ID), onlineStatusNotice)
	}
	return token, nil
}

func (u *UserService) Logout(id uint) {
	//更新等出时间
	updates := make(map[string]interface{})
	updates["logout_time"] = time.Now().Unix()
	_ = u.userRepository.UpdateFields(id, updates)
	u.wsHandler.OnlineStatusNotice(int64(id), model.OnlineStatusNotice{
		UserId:       id,
		OnlineStatus: model.Offline,
		ActionType:   model.LogoutAction,
	})
}

func (u *UserService) OnlineStatusChange(id uint, onlineStatus model.OnlineStatus) error {
	updates := make(map[string]interface{})
	updates["online_status"] = onlineStatus
	err := u.userRepository.UpdateFields(id, updates)
	if err != nil {
		return err
	}
	go u.wsHandler.OnlineStatusNotice(int64(id), model.OnlineStatusNotice{
		UserId:       id,
		OnlineStatus: onlineStatus,
		ActionType:   model.StatusChangeAction,
	})
	return nil
}

func (u *UserService) UpdateUser(updateRequest *request.UserUpdateRequest) error {
	err := db.Mysql.Transaction(func(tx *gorm.DB) error {

		user, err := u.userRepository.GetById(updateRequest.ID, tx)
		if err != nil {
			return fmt.Errorf("查找用户失败: %v", err)
		}
		if user == nil {
			return errors.New("用户不存在")
		}
		updates := make(map[string]interface{})
		if updateRequest.Nickname != nil {
			updates["nickname"] = updateRequest.Nickname
		}
		if updateRequest.Desc != nil {
			updates["desc"] = updateRequest.Desc
		}
		if updateRequest.Avatar != nil {
			updates["avatar"] = updateRequest.Avatar
		}
		if updateRequest.Phone != nil {
			updates["phone"] = updateRequest.Phone
		}
		if updateRequest.Email != nil {
			updates["email"] = updateRequest.Email
		}
		if len(updates) == 0 {
			return errors.New("没有可更新的字段")
		}
		err = u.userRepository.UpdateFields(user.ID, updates, tx)
		if err != nil {
			return fmt.Errorf("更新用户失败: %v", err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) GetUserInfo(id uint) (response.UserVO, error) {
	userVo, err := u.userRepository.GetVoById(id)
	return userVo, err
}

func (u *UserService) UpdateHeartbeatTime(userId int64, time int64) error {
	return u.userRepository.UpdateHeartbeatTime(userId, time)
}
func (u *UserService) CheckOfflineUsers() error {
	//todo 可以抽离为配置
	timeout := 30 * time.Second
	cutoffTime := time.Now().Add(-timeout).Unix()

	users, err := u.userRepository.GetUsersWithHeartbeatBefore(cutoffTime)
	if err != nil {
		return fmt.Errorf("查询心跳超时用户失败: %w", err)
	}

	for _, user := range users {
		if user.OnlineStatus == model.Online {

			updates := map[string]interface{}{
				"online_status": model.Offline,
			}
			if err := u.userRepository.UpdateFields(user.ID, updates); err != nil {
				continue
			}

			u.wsHandler.OnlineStatusNotice(int64(user.ID), model.OnlineStatusNotice{
				UserId:       user.ID,
				OnlineStatus: model.Offline,
				ActionType:   model.HeartbeatAction,
			})
			logUtil.Infof("用户(%d)心跳检测不通过,已主动下线", user.ID)
		}
	}

	return nil
}
