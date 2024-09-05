package routers

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-programming-tour-book/blog-service/internal/middleware"
	"github.com/go-programming-tour-book/blog-service/internal/routers/api"
	v1 "github.com/go-programming-tour-book/blog-service/internal/routers/api/v1"

	_ "github.com/go-programming-tour-book/blog-service/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 添加 CORS 中间件
	r.Use(cors.New(corsConfig()))

	// 设置 Swagger文档地址
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	article := v1.NewArticle()
	tag := v1.NewTag()
	user := v1.NewUser()
	r.POST("auth", api.GetAuth)
	apiv1 := r.Group("/api/v1")
	// JWT中间件使用
	apiv1.Use(middleware.JWT())
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
		apiv1.PUT("/user/:id", user.Update)
		apiv1.DELETE("/user/:id", user.Delete)

	}

	return r
}

// CORS 中间件
func corsConfig() cors.Config {
	return cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
}
