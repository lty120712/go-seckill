package model

import (
	"github.com/go-playground/validator/v10"
	"go-chat/internal/model"
	"regexp"
)

// UserUpdateRequest 用户基础信息更新请求结构体
type UserUpdateRequest struct {
	ID       uint    `json:"id"`
	Desc     *string `json:"desc" validate:"omitempty,max=500"`
	Avatar   *string `json:"avatar"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Nickname *string `json:"nickname" validate:"omitempty,max=60"`
	Phone    *string `json:"phone" validate:"omitempty,phone"`
}

func phoneValidator(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	re := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`) // 校验手机号格式
	if !re.MatchString(phone) {
		return false
	}
	return true
}

func (r *UserUpdateRequest) Validate() error {
	validate := validator.New()
	// 注册自定义校验规则，确保它被注册
	err := validate.RegisterValidation("phone", phoneValidator)
	if err != nil {
		return err
	}
	err = validate.Struct(r)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.Tag() == "email" {
				return &model.ValidationError{
					Field:   err.Field(),
					Message: "邮箱格式错误",
				}
			} else if err.Tag() == "phone" {
				return &model.ValidationError{
					Field:   err.Field(),
					Message: "手机号格式错误",
				}
			}
		}
	}
	return nil
}
