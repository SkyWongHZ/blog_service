package v1

import "github.com/gin-gonic/gin"

type Article struct{}

func NewArticle() Tag {
	return Tag{}
}

func (t Article) Get(c *gin.Context)    {}
func (t Article) List(c *gin.Context)   {}
func (t Article) Create(c *gin.Context) {}
func (t Article) Update(c *gin.Context) {}
func (t Article) Delete(c *gin.Context) {}
