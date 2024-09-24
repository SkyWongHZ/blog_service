package api

import (
	"github.com/go-programming-tour-book/blog-service/pkg/redis"

	"github.com/gin-gonic/gin"
	"github.com/go-programming-tour-book/blog-service/global"
	"github.com/go-programming-tour-book/blog-service/internal/service"
	"github.com/go-programming-tour-book/blog-service/pkg/app"
	"github.com/go-programming-tour-book/blog-service/pkg/errcode"
)

// @Summary 获取授权Token
// @Produce json
// @Param app_key query string true "应用关键字"
// @Param app_secret query string true "应用密钥"
// @Success 200 {object} map[string]interface{} "成功返回Token"
// @Failure 400 {object} errcode.Error "请求参数错误"
// @Failure 401 {object} errcode.Error "授权失败"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /auth [post]
func GetAuth(c *gin.Context) {
	param := service.AuthRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf(c, "app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.CheckAuth(&param)
	if err != nil {
		global.Logger.Errorf(c, "svc.CheckAuth err: %v", err)
		response.ToErrorResponse(errcode.UnauthorizedAuthNotExist)
		return
	}

	// Redis黑名单
	banned, err := redis.IsUserBanned(param.AppKey)
	if err != nil {
		global.Logger.Errorf(c, "redis.IsUserBanned err: %v", err)
		response.ToErrorResponse(errcode.ServerError)
		return
	}
	if banned {
		response.ToErrorResponse(errcode.UnauthorizedAppKeyBanned)
		return
	}

	token, err := app.GenerateToken(param.AppKey, param.AppSecret)
	if err != nil {
		global.Logger.Errorf(c, "app.GenerateToken err: %v", err)
		response.ToErrorResponse(errcode.UnauthorizedTokenGenerate)
		return
	}

	response.ToResponse(gin.H{
		"token": token,
	})
}
