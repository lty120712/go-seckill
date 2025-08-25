package model

import (
	"github.com/go-playground/validator/v10"
	"go-chat/internal/model"
)

type GroupCreateRequest struct {
	UserId     uint    `json:"user_id" binding:"required"`                // 创建者ID
	Name       string  `json:"name" binding:"required" validate:"max=50"` // 群名称
	MaxNum     int     `json:"max_num" validate:"min=1,max=2000"`         // 群最大人数
	MemberList *[]uint `json:"member_list"`                               // 初始群成员
}

func (r *GroupCreateRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Name" {
				if err.Tag() == "max" {
					return &model.ValidationError{
						Field:   err.Field(),
						Message: "群名称应小于50字",
					}
				}
			}
			if err.Field() == "MaxNum" {
				if err.Tag() == "min" {
					return &model.ValidationError{
						Field:   err.Field(),
						Message: "群人数应大于1",
					}
				}
				if err.Tag() == "max" {
					return &model.ValidationError{
						Field:   err.Field(),
						Message: "群人数应小于2000",
					}
				}
			}
		}
	}
	return nil
}
