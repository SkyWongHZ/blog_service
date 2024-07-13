package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/go-programming-tour-book/blog-service/global"
	"github.com/go-programming-tour-book/blog-service/internal/service"
	"github.com/go-programming-tour-book/blog-service/pkg/app"
	"github.com/go-programming-tour-book/blog-service/pkg/errcode"
)

type User struct{}

func NewUser() User {
	return User{}
}

func (t User) List(c *gin.Context) {
	param := service.ListUserRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf(c, "app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	pager := app.Pager{Page: app.GetPage(c), PageSize: app.GetPageSize(c)}
	totalRows, err := svc.CountUser(&service.CountUserRequest{Username: param.Username, Email: param.Email, State: param.State})
	if err != nil {
		global.Logger.Errorf(c, "svc.CountUser err: %v", err)
		response.ToErrorResponse(errcode.ErrorCountUserFail)
		return
	}

	users, err := svc.ListUser(&param, &pager)
	if err != nil {
		global.Logger.Errorf(c, "svc.ListUser err: %v", err)
		response.ToErrorResponse(errcode.ErrorListUserFail)
		return
	}

	response.ToResponseList(users, totalRows)
	return
}

func (t User) Create(c *gin.Context) {
	param := service.RegisterUserRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf(c, "app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())

	err := svc.RegisterUser(c.Request.Context(), &param)
	if err != nil {
		global.Logger.Errorf(c, "svc.RegisterUser err: %v", err)
		response.ToErrorResponse(errcode.ErrorRegisterUserFail)
		return
	}

	response.ToResponse(gin.H{"新增接口成功": "200"})
	return
}
