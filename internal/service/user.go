package service

import (
	"github.com/go-programming-tour-book/blog-service/internal/model"
	"github.com/go-programming-tour-book/blog-service/pkg/app"
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
	Password string `form:"password" binding:"max=255"`
}

func (svc *Service) CountUser(param *CountUserRequest) (int, error) {
	return svc.dao.CountUser(param.Username, param.Email, param.State)
}

func (svc *Service) ListUser(param *ListUserRequest, pager *app.Pager) ([]*model.User, error) {
	return svc.dao.ListUser(param.Username, param.Email, param.State, pager.Page, pager.PageSize)
}

func (svc *Service) RegisterUser(param *RegisterUserRequest) error {
	return svc.dao.RegisterUser(param.Username, param.Email, param.Password, param.State)
}
