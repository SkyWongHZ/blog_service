package service

import (
	"errors"

	"github.com/go-programming-tour-book/blog-service/internal/model"
	"github.com/go-programming-tour-book/blog-service/pkg/app"
	"github.com/jinzhu/gorm"
)

type CountUserRequest struct {
	Username string `form:"username" binding:"max=100"`
	Email    string `form:"email" binding:"max=255"`
	State    uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type ListUserRequest struct {
	Username string `form:"username" binding:"max=100"`
	Email    string `form:"email" binding:"max=255"`
	State    uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type RegisterUserRequest struct {
	Username string `form:"username" binding:"max=100"`
	Email    string `form:"email" binding:"max=255"`
	State    uint8  `form:"state,default=1" binding:"oneof=0 1"`
	Password string `form:"password" binding:"max=255,min=6"`
}

type UpdateUserRequest struct {
	ID         uint32 `form:"id" binding:"required,gte=1"`
	Name       string `form:"name" binding:"max=100"`
	State      uint8  `form:"state" binding:"oneof=0 1"`
	ModifiedBy string `form:"modified_by" min=2,max=100"`
}

type DeleteUserRequest struct {
	ID uint32 `form:"id" binding:"required,gte=1"`
}

func (svc *Service) CountUser(param *CountUserRequest) (int, error) {
	return svc.dao.CountUser(param.Username, param.Email, param.State)
}

func (svc *Service) ListUser(param *ListUserRequest, pager *app.Pager) ([]*model.User, error) {
	return svc.dao.ListUser(param.Username, param.Email, param.State, pager.Page, pager.PageSize)
}

func (svc *Service) RegisterUser(param *RegisterUserRequest) error {
	// return svc.dao.RegisterUser(param.Username, param.Email, param.Password, param.State)
	err := svc.dao.RegisterUser(param.Username, param.Email, param.Password, param.State)
	if err != nil {
		// 检查是否是唯一性约束违反错误
		if gorm.IsRecordNotFoundError(err) {
			return errors.New("username already exists")
		}
		return err
	}
	return nil
}

func (svc *Service) UpdateUser(param *UpdateUserRequest) error {
	return svc.dao.UpdateUser(param.ID, param.Name, param.State, param.ModifiedBy)
}

func (svc *Service) DeleteUser(param *DeleteUserRequest) error {
	return svc.dao.DeleteUser(param.ID)
}
