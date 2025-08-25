package repository

import (
	"github.com/lty120712/gorm-pagination/pagination"
	"go-chat/internal/db"
	"go-chat/internal/model"
	request "go-chat/internal/model/request"
	"gorm.io/gorm"
	"sync"
)

type GroupRepository struct {
}

var (
	GroupRepositoryInstance *GroupRepository
	groupOnce               sync.Once
)

func InitGroupRepository() {
	groupOnce.Do(func() {
		GroupRepositoryInstance = &GroupRepository{}
	})

}

func (g *GroupRepository) ExistsByCode(code string, tx ...*gorm.DB) bool {
	gormDB := db.GetGormDB(tx...)
	var exists bool
	rawSql := "SELECT EXISTS(SELECT 1 FROM `groups` WHERE code = ?)"
	err := gormDB.Raw(rawSql, code).Scan(&exists).Error

	if err != nil {
		return false
	}
	return exists
}

func (g *GroupRepository) Save(group *model.Group, tx ...*gorm.DB) error {
	gormDB := db.GetGormDB(tx...)
	return gormDB.Create(group).Error
}

func (g *GroupRepository) Page(req request.GroupSearchRequest, tx ...*gorm.DB) (*pagination.PageResult[model.Group], error) {
	gormDB := db.GetGormDB(tx...)
	query := gormDB.Model(&model.Group{})

	if req.Code != "" {
		query = query.Where("code = ?", req.Code)
	}
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.UserId != 0 {
		query = query.
			Joins("JOIN group_members ON group_members.group_id = groups.id").
			Where("group_members.member_id = ?", req.UserId)
	}
	result := &pagination.PageResult[model.Group]{Records: []model.Group{}}

	// 调用分页函数并获取结果
	_, err := pagination.Paginate(query, req.Page, req.PageSize, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (g *GroupRepository) GetByID(groupID uint, tx ...*gorm.DB) (*model.Group, error) {
	gormDB := db.GetGormDB(tx...)
	var group model.Group
	err := gormDB.First(&group, groupID).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (g *GroupRepository) Delete(groupID uint, tx ...*gorm.DB) error {
	gormDB := db.GetGormDB(tx...)
	return gormDB.Delete(&model.Group{}, groupID).Error
}

func (g *GroupRepository) Update(groupId uint, m map[string]interface{}, tx ...*gorm.DB) error {
	if len(m) == 0 {
		return nil
	}
	gormDB := db.GetGormDB(tx...)
	return gormDB.Model(&model.Group{}).Where("id = ?", groupId).Updates(m).Error
}
