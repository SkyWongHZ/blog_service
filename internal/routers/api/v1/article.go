package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/go-programming-tour-book/blog-service/global"
	"github.com/go-programming-tour-book/blog-service/internal/service"
	"github.com/go-programming-tour-book/blog-service/pkg/app"
	"github.com/go-programming-tour-book/blog-service/pkg/convert"
	"github.com/go-programming-tour-book/blog-service/pkg/errcode"
)

type Article struct {
}

func NewArticle() Article {
	return Article{}
}

// @Summary 获取单个文章
// @Produce json
// @Param id path int true "文章ID"
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} Article "文章详细信息"
// @Failure 400 {object} errcode.Error "请求错误"
// @Router /api/v1/articles/{id} [get]
func (a Article) Get(c *gin.Context) {
	param := service.ArticleRequest{ID: convert.StrTo(c.Param("id")).MustUInt32()}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf(c, "app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	article, err := svc.GetArticle(&param)
	if err != nil {
		global.Logger.Errorf(c, "svc.GetArticle err: %v", err)
		response.ToErrorResponse(errcode.ErrorGetArticleFail)
		return
	}

	response.ToResponse(article)
	return
}

// @Summary 获取文章列表
// @Produce json
// @Param name query string false "文章名称"
// @Param tag_id query int true "标签ID"
// @Param state query int false "状态" Enums(0, 1) default(1)
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param Authorization header string true "Bearer Token"
// @Success 200 {array} Article "文章列表"
// @Failure 400 {object} errcode.Error "请求错误"
// @Router /api/v1/articles [get]
func (a Article) List(c *gin.Context) {
	param := service.ArticleListRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf(c, "app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	pager := app.Pager{Page: app.GetPage(c), PageSize: app.GetPageSize(c)}
	articles, totalRows, err := svc.GetArticleList(&param, &pager)
	if err != nil {
		global.Logger.Errorf(c, "svc.GetArticleList err: %v", err)
		response.ToErrorResponse(errcode.ErrorGetArticlesFail)
		return
	}

	response.ToResponseList(articles, totalRows)
	return
}

// @Summary 创建文章
// @Produce json
// @Param title body string true "文章标题"
// @Param desc body string false "文章描述"
// @Param content body string true "文章内容"
// @Param cover_image_url body string false "封面图片地址"
// @Param state body int false "状态" Enums(0, 1) default(1)
// @Param tag_ids body []int true "标签ID集合"
// @Param Authorization header string true "Bearer Token"
// @Success 201 {object} Article "文章创建成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Router /api/v1/articles [post]
func (a Article) Create(c *gin.Context) {
	param := service.CreateArticleRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf(c, "app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.CreateArticle(&param)
	if err != nil {
		global.Logger.Errorf(c, "svc.CreateArticle err: %v", err)
		response.ToErrorResponse(errcode.ErrorCreateArticleFail)
		return
	}

	response.ToResponse(gin.H{})
	return
}

// @Summary 更新文章
// @Produce json
// @Param id path int true "文章ID"
// @Param title body string false "文章标题"
// @Param desc body string false "文章描述"
// @Param content body string false "文章内容"
// @Param cover_image_url body string false "封面图片地址"
// @Param state body int false "状态" Enums(0, 1)
// @Param tag_ids body []int false "标签ID集合"
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} Article "文章更新成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Router /api/v1/articles/{id} [put]
func (a Article) Update(c *gin.Context) {
	param := service.UpdateArticleRequest{ID: convert.StrTo(c.Param("id")).MustUInt32()}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf(c, "app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.UpdateArticle(&param)
	if err != nil {
		global.Logger.Errorf(c, "svc.UpdateArticle err: %v", err)
		response.ToErrorResponse(errcode.ErrorUpdateArticleFail)
		return
	}

	response.ToResponse(gin.H{})
	return
}

// @Summary 删除文章
// @Produce json
// @Param id path int true "文章ID"
// @Param Authorization header string true "Bearer Token"
// @Success 204 "文章删除成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Router /api/v1/articles/{id} [delete]
func (a Article) Delete(c *gin.Context) {
	param := service.DeleteArticleRequest{ID: convert.StrTo(c.Param("id")).MustUInt32()}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf(c, "app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.DeleteArticle(&param)
	if err != nil {
		global.Logger.Errorf(c, "svc.DeleteArticle err: %v", err)
		response.ToErrorResponse(errcode.ErrorDeleteArticleFail)
		return
	}

	response.ToResponse(gin.H{})
	return
}

// @Summary 获取热门文章
// @Produce json
// @Success 200 {array} Article "热门文章列表"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/articles/hot [get]
func (a Article) GetHotArticles(c *gin.Context) {
	response := app.NewResponse(c)
	svc := service.New(c.Request.Context())
	articles, err := svc.GetHotArticles()
	if err != nil {
		global.Logger.Errorf(c, "svc.GetHotArticles err: %v", err)
		response.ToErrorResponse(errcode.ErrorGetHotArticlesFail)
		return
	}

	response.ToResponse(articles)
}
