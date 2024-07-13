package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-programming-tour-book/blog-service/internal/routers/api"
	v1 "github.com/go-programming-tour-book/blog-service/internal/routers/api/v1"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	article := v1.NewArticle()
	tag := v1.NewTag()
	user := v1.NewUser()
	r.POST("auth", api.GetAuth)
	apiv1 := r.Group("/api/v1")
	// JWT中间件使用
	// apiv1.Use(middleware.JWT())
	{

		apiv1.POST("/tags", tag.Create)
		apiv1.DELETE("/tags/:id", tag.Delete)
		apiv1.PUT("/tags/:id", tag.Update)
		apiv1.PATCH("/tags/:id/state", tag.Update)
		apiv1.GET("/tags", tag.List)

		apiv1.POST("/articles", article.Create)
		apiv1.DELETE("/articles/:id", article.Delete)
		apiv1.PUT("/articles/:id", article.Update)
		apiv1.PATCH("/articles/:id/state", article.Update)
		apiv1.GET("/articles/:id", article.Get)
		apiv1.GET("/articles", article.List)

		apiv1.GET("/user", user.List)
		apiv1.POST("/user", user.Create)

	}

	return r
}
