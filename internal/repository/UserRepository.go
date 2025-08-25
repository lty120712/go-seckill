package repository

import (
	"errors"
	"go-chat/internal/db"
	"go-chat/internal/model"
	response "go-chat/internal/model/response"
	"gorm.io/gorm"
	"sync"
)

type UserRepository struct {
}

var (
	UserRepositoryInstance *UserRepository
	userOnce               sync.Once
)

func InitUserRepository() {
	userOnce.Do(func() {
		UserRepositoryInstance = &UserRepository{}
	})

}

// GetById
func (r *UserRepository) GetById(id uint, tx ...*gorm.DB) (user *model.User, err error) {
	user = &model.User{}
	gormDB := db.GetGormDB(tx...)
	err = gormDB.Where("id = ?", id).First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}

// GetByName 根据用户名查询用户
func (r *UserRepository) GetByName(username *string, tx ...*gorm.DB) (user *model.User, err error) {
	user = &model.User{}
	gormDB := db.GetGormDB(tx...)
	err = gormDB.Where("username = ?", username).First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return user, err
}

// Save 保存用户
func (r *UserRepository) Save(user *model.User, tx ...*gorm.DB) (err error) {
	gormDB := db.GetGormDB(tx...)
	err = gormDB.Create(user).Error
	return
}

// UpdateFields 更新用户字段
func (r *UserRepository) UpdateFields(id uint, updates map[string]interface{}, tx ...*gorm.DB) error {
	if len(updates) == 0 {
		return nil
	}
	gormDB := db.GetGormDB(tx...)
	// 执行更新操作
	return gormDB.Model(&model.User{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *UserRepository) GetNickNamesByIds(ids []uint, tx ...*gorm.DB) (map[uint]string, error) {
	gormDB := db.GetGormDB(tx...)
	var users []*model.User
	// 执行查询，只选择 id 和 nickname 字段
	err := gormDB.Select("id, nickname").Where("id IN ?", ids).Find(&users).Error
	if err != nil {
		return nil, err
	}
	idNicknameMap := make(map[uint]string)
	for _, user := range users {
		idNicknameMap[user.ID] = *user.Nickname
	}
	return idNicknameMap, nil
}

func (r *UserRepository) GetNickNamesById(id uint, tx ...*gorm.DB) (nickname string, err error) {
	gormDB := db.GetGormDB(tx...)
	err = gormDB.Model(&model.User{}).Select("nickname").Where("id = ?", id).Scan(&nickname).Error
	return nickname, err
}

func (r *UserRepository) GetByIdList(userIdList []uint, tx ...*gorm.DB) (userList []model.User, err error) {
	if len(userIdList) == 0 {
		return []model.User{}, nil
	}
	gormDB := db.GetGormDB(tx...)
	err = gormDB.Where("id IN ?", userIdList).Find(&userList).Error
	return
}

func (r *UserRepository) GetVoById(id uint, tx ...*gorm.DB) (userVo response.UserVO, err error) {
	gormDB := db.GetGormDB(tx...)
	err = gormDB.Table("users").
		Select("id, username, nickname, `desc`, phone, email, avatar, client_ip, client_port, login_time, heartbeat_time, logout_time, `status`, online_status, device_info").
		Where("id = ?", id).
		Scan(&userVo).Error

	return userVo, err
}

func (r *UserRepository) UpdateHeartbeatTime(userId int64, heartbeatTime int64, tx ...*gorm.DB) error {
	gormDB := db.GetGormDB(tx...)

	// heartbeat_time 总是更新
	// online_status 仅当当前状态为 Offline 时更新为 Online
	return gormDB.Model(&model.User{}).
		Where("id = ?", userId).
		Updates(map[string]interface{}{
			"heartbeat_time": heartbeatTime,
			"online_status": gorm.Expr(
				"CASE WHEN online_status = ? THEN ? ELSE online_status END",
				model.Offline, model.Online,
			),
		}).Error
}
func (r *UserRepository) GetUsersWithHeartbeatBefore(cutoffTime int64, tx ...*gorm.DB) ([]model.User, error) {
	gormDB := db.GetGormDB(tx...)
	var users []model.User
	err := gormDB.Where("heartbeat_time < ? AND online_status = ?", cutoffTime, model.Online).
		Find(&users).Error
	return users, err
}
