package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/go-programming-tour-book/blog-service/global"
	"github.com/go-programming-tour-book/blog-service/internal/service"
	"github.com/go-programming-tour-book/blog-service/pkg/app"
	"github.com/go-programming-tour-book/blog-service/pkg/convert"
	"github.com/go-programming-tour-book/blog-service/pkg/errcode"
)

type User struct{}

func NewUser() User {
	return User{}
}

// @Summary 用户列表
// @Produce json
// @Param username query string false "用户名"
// @Param email query string false "邮箱"
// @Param state query int false "状态"
// @Param page query int false "页码"
// @Param page_size query int false "每页显示数量"
// @Param Authorization header string true "Bearer Token"
// @Success 200 {array} User "用户列表"
// @Failure 400 {object} errcode.Error "请求参数错误"
// @Router /api/v1/user [get]
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

// @Summary 创建用户
// @Produce json
// @Param username body string true "用户名"
// @Param email body string true "邮箱"
// @Param password body string true "密码"
// @Param state body int false "状态"
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} map[string]interface{} "用户创建成功"
// @Failure 400 {object} errcode.Error "请求参数错误"
// @Router /api/v1/user [post]
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

	err := svc.RegisterUser(&param)

	if err != nil {
		global.Logger.Errorf(c, "svc.RegisterUser err: %v", err)

		if err.Error() == "username already exists" {
			response.ToErrorResponse(errcode.NewError(20010001, "用户名已存在"))

			response.ToErrorResponse(errcode.ErrorRegisterUserFail)
			return
		} else {
			response.ToResponse(gin.H{"新增接口成功": "200"})
		}
	}

	return
}

// @Summary 更新用户
// @Produce json
// @Param id path int true "用户ID"
// @Param username body string false "用户名"
// @Param email body string false "邮箱"
// @Param state body int false "状态"
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} map[string]interface{} "用户更新成功"
// @Failure 400 {object} errcode.Error "请求参数错误"
// @Router /api/v1/user/{id} [put]
func (t User) Update(c *gin.Context) {
	param := service.UpdateUserRequest{ID: convert.StrTo(c.Param("id")).MustUInt32()}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf(c, "app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.UpdateUser(&param)
	if err != nil {
		global.Logger.Errorf(c, "svc.UpdateUser err: %v", err)
		response.ToErrorResponse(errcode.ErrorUpdateUserFail)
		return
	}

	response.ToResponse(gin.H{"更新成功": "200"})
	return
}

// @Summary 删除用户
// @Produce json
// @Param id path int true "用户ID"
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} map[string]interface{} "用户删除成功"
// @Failure 400 {object} errcode.Error "请求参数错误"
// @Router /api/v1/user/{id} [delete]
func (t User) Delete(c *gin.Context) {
	param := service.DeleteUserRequest{ID: convert.StrTo(c.Param("id")).MustUInt32()}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf(c, "app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.DeleteUser(&param)
	if err != nil {
		global.Logger.Errorf(c, "svc.DeleteUser err: %v", err)
		response.ToErrorResponse(errcode.ErrorDeleteUserFail)
		return
	}
	response.ToResponse(gin.H{"删除成功": "200"})
	return
}
