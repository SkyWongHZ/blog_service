package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/url"
	"path"
	"path/filepath"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-programming-tour-book/blog-service/global"
	"github.com/go-programming-tour-book/blog-service/internal/dao"
	"github.com/go-programming-tour-book/blog-service/internal/model"
	"github.com/go-programming-tour-book/blog-service/pkg/app"
	"github.com/go-programming-tour-book/blog-service/pkg/minio"
)

type ArticleRequest struct {
	ID    uint32 `form:"id" binding:"required,gte=1"`
	State uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type ArticleListRequest struct {
	TagID *uint32 `form:"tag_id" binding:"omitempty,gte=1"`
	State uint8   `form:"state,default=1" binding:"oneof=0 1"`
}

type CreateArticleRequest struct {
	TagID         uint32                `form:"tag_id" binding:"required,gte=1"`
	Title         string                `form:"title" binding:"required,min=2,max=100"`
	Desc          string                `form:"desc" binding:"required,min=2,max=255"`
	Content       string                `form:"content" binding:"required,min=2,max=4294967295"`
	CoverImageUrl *multipart.FileHeader `form:"cover_image_url" binding:"omitempty"`
	CreatedBy     string                `form:"created_by" binding:"required,min=2,max=100"`
	State         uint8                 `form:"state,default=1" binding:"oneof=0 1"`
}

type UpdateArticleRequest struct {
	ID            uint32                `form:"id" binding:"required,gte=1"`
	TagID         uint32                `form:"tag_id" binding:"required,gte=1"`
	Title         string                `form:"title" binding:"min=2,max=100"`
	Desc          string                `form:"desc" binding:"min=2,max=255"`
	Content       string                `form:"content" binding:"min=2,max=4294967295"`
	CoverImageUrl *multipart.FileHeader `form:"cover_image_url" binding:"omitempty"`
	ModifiedBy    string                `form:"modified_by" binding:"required,min=2,max=100"`
	State         uint8                 `form:"state,default=1" binding:"oneof=0 1"`
}

type DeleteArticleRequest struct {
	ID uint32 `form:"id" binding:"required,gte=1"`
}

type Article struct {
	ID            uint32     `json:"id"`
	Title         string     `json:"title"`
	Desc          string     `json:"desc"`
	Content       string     `json:"content"`
	CoverImageUrl string     `json:"cover_image_url"`
	State         uint8      `json:"state"`
	Tag           *model.Tag `json:"tag"`
	CoverImageKey string     `json:"cover_image_key"`
}

func (svc *Service) GetArticle(param *ArticleRequest) (*Article, error) {
	article, err := svc.dao.GetArticle(param.ID, param.State)
	if err != nil {
		return nil, err
	}

	articleTag, err := svc.dao.GetArticleTagByAID(article.ID)
	if err != nil {
		return nil, err
	}

	tag, err := svc.dao.GetTag(articleTag.TagID, model.STATE_OPEN)
	if err != nil {
		return nil, err
	}
	fmt.Println("article.CoverImageKey", article.CoverImageKey)
	// 使用minio生成封面图片URL
	coverImageURL, err := svc.GenerateImageURL(article.CoverImageKey)

	// oss
	// coverImageURL, err := svc.GenerateOSSImageURL(article.CoverImageKey)
	if err != nil {
		return nil, err
	}

	return &Article{
		ID:            article.ID,
		Title:         article.Title,
		Desc:          article.Desc,
		Content:       article.Content,
		CoverImageUrl: coverImageURL,
		State:         article.State,
		Tag:           &tag,
	}, nil
}

// 新增方法用于生成图片URL
func (svc *Service) GenerateImageURL(imageKey string) (string, error) {

	// 确保 imageKey 不为空
	if imageKey == "" {
		return "", errors.New("image key is empty")
	}
	baseURL := "http://118.178.184.13:9000"
	bucket := "articles" // 使用代码中定义的 bucket 名称

	// URL 编码 imageKey
	encodedKey := url.PathEscape(imageKey)

	// 构建完整的 URL
	fullURL := fmt.Sprintf("%s/%s/%s", baseURL, bucket, encodedKey)

	// 规范化 URL 路径
	cleanURL := path.Clean(fullURL)

	return cleanURL, nil
}

// oss生成url
func (svc *Service) GenerateOSSImageURL(imageKey string) (string, error) {
	if imageKey == "" {
		return "", errors.New("image key is empty")
	}

	// 假设我们已经在全局变量中设置了OSSClient
	client := global.OSSClient
	bucket, err := client.Bucket(global.OSSSetting.BucketName)
	if err != nil {
		return "", fmt.Errorf("failed to get bucket: %w", err)
	}

	// 生成签名URL，有效期为1小时
	signedURL, err := bucket.SignURL(imageKey, oss.HTTPGet, 3600)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return signedURL, nil
}

func (svc *Service) GetArticleList(param *ArticleListRequest, pager *app.Pager) ([]*Article, int, error) {
	var tagID uint32
	if param.TagID != nil {
		tagID = *param.TagID
	} else {
		tagID = 0
	}

	articleCount, err := svc.dao.CountArticleListByTagID(tagID, param.State)
	fmt.Println("articleCount", articleCount)
	if err != nil {
		return nil, 0, err
	}

	articles, err := svc.dao.GetArticleListByTagID(tagID, param.State, pager.Page, pager.PageSize)
	if err != nil {
		return nil, 0, err
	}

	var articleList []*Article
	for _, article := range articles {
		// 生成图片 URL
		coverImageURL := ""
		if article.CoverImageKey != "" {
			var err error
			coverImageURL, err = svc.GenerateImageURL(article.CoverImageKey)
			if err != nil {
				// 如果生成 URL 失败，记录错误并使用空字符串
				global.Logger.Errorf(context.Background(), "Failed to generate image URL for article %d: %v", article.ArticleID, err)
			}
		}

		articleList = append(articleList, &Article{
			ID:            article.ArticleID,
			Title:         article.ArticleTitle,
			Desc:          article.ArticleDesc,
			Content:       article.Content,
			CoverImageUrl: coverImageURL,
			// CoverImageKey: article.CoverImageKey,
			Tag:           &model.Tag{Model: &model.Model{ID: article.TagID}, Name: article.TagName},
		})
	}

	return articleList, articleCount, nil
}

func (svc *Service) CreateArticle(param *CreateArticleRequest) error {
	// minio处理文件上传
	coverImageKey := ""
	if param.CoverImageUrl != nil {
		fmt.Printf("Processing cover image: %s", param.CoverImageUrl.Filename)
		objectName := fmt.Sprintf("cover_images/%d%s", time.Now().UnixNano(), filepath.Ext(param.CoverImageUrl.Filename))
		bucketName := "articles"

		file, err := param.CoverImageUrl.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = minio.MinioClient.PutObject(
			context.Background(),
			bucketName,
			objectName,
			file,
			param.CoverImageUrl.Size,
			minio.PutObjectOptions{ContentType: param.CoverImageUrl.Header.Get("Content-Type")},
		)
		if err != nil {
			return err
		}

		coverImageKey = objectName
	}

	// 阿里云oss处理文件上传
	// coverImageKey := ""
	// if param.CoverImageUrl != nil {
	// 	file, err := param.CoverImageUrl.Open()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	defer file.Close()

	// 	objectKey := fmt.Sprintf("cover_images/%d%s", time.Now().UnixNano(), filepath.Ext(param.CoverImageUrl.Filename))

	// 	bucket, err := global.OSSClient.Bucket(global.OSSSetting.BucketName)
	// 	if err != nil {
	// 		fmt.Printf("Failed to get bucket: %v", err)
	// 		return err
	// 	}

	// 	err = bucket.PutObject(objectKey, file)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	coverImageKey = objectKey
	// }

	article, err := svc.dao.CreateArticle(&dao.Article{
		Title:         param.Title,
		Desc:          param.Desc,
		Content:       param.Content,
		CoverImageKey: coverImageKey,
		State:         param.State,
		CreatedBy:     param.CreatedBy,
	})
	if err != nil {
		return err
	}

	err = svc.dao.CreateArticleTag(article.ID, param.TagID, param.CreatedBy)
	if err != nil {
		return err
	}

	return nil
}

func (svc *Service) UpdateArticle(param *UpdateArticleRequest) error {

	// 获取原有的文章信息
	existingArticle, err := svc.dao.GetArticle(param.ID, param.State)
	if err != nil {
		return err
	}

	// MinIO 处理文件上传
	coverImageKey := existingArticle.CoverImageKey // 默认使用原有的封面图片键
	if param.CoverImageUrl != nil {
		fmt.Printf("Processing new cover image: %s", param.CoverImageUrl.Filename)
		objectName := fmt.Sprintf("cover_images/%d%s", time.Now().UnixNano(), filepath.Ext(param.CoverImageUrl.Filename))
		bucketName := "articles"

		file, err := param.CoverImageUrl.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = minio.MinioClient.PutObject(
			context.Background(),
			bucketName,
			objectName,
			file,
			param.CoverImageUrl.Size,
			minio.PutObjectOptions{ContentType: param.CoverImageUrl.Header.Get("Content-Type")},
		)
		if err != nil {
			return err
		}

		// 如果上传成功，更新封面图片键
		coverImageKey = objectName

		// 删除旧的封面图片（如果存在）
		if existingArticle.CoverImageKey != "" {
			err = minio.MinioClient.RemoveObject(context.Background(), bucketName, existingArticle.CoverImageKey, minio.RemoveObjectOptions{})
			if err != nil {
				fmt.Printf("Failed to remove old cover image: %v", err)
				// 不返回错误，继续更新文章
			}
		}
	}

	err = svc.dao.UpdateArticle(&dao.Article{
		ID:            param.ID,
		Title:         param.Title,
		Desc:          param.Desc,
		Content:       param.Content,
		CoverImageKey: coverImageKey,
		State:         param.State,
		ModifiedBy:    param.ModifiedBy,
	})
	if err != nil {
		return err
	}

	err = svc.dao.UpdateArticleTag(param.ID, param.TagID, param.ModifiedBy)
	if err != nil {
		return err
	}

	return nil
}

func (svc *Service) DeleteArticle(param *DeleteArticleRequest) error {
	err := svc.dao.DeleteArticle(param.ID)
	if err != nil {
		return err
	}

	err = svc.dao.DeleteArticleTag(param.ID)
	if err != nil {
		return err
	}

	return nil
}

func (svc *Service) GetHotArticles() ([]*Article, error) {
	articles, err := svc.dao.GetHotArticles()
	if err != nil {
		return nil, err
	}

	var hotArticles []*Article
	for _, article := range articles {

		coverImageURL, err := svc.GenerateImageURL(article.CoverImageKey)
		if err != nil {
			// 如果生成URL失败，可以选择记录错误并继续，或者返回错误
			// 这里我们选择记录错误并继续，使用空字符串作为URL
			c := context.Background()
			global.Logger.Errorf(c, "Failed to generate image URL for article %d: %v", article.ID, err)
			coverImageURL = ""
		}

		hotArticles = append(hotArticles, &Article{
			ID:            article.ID,
			Title:         article.Title,
			Desc:          article.Desc,
			Content:       article.Content,
			CoverImageKey: article.CoverImageKey,
			CoverImageUrl: coverImageURL,
			State:         article.State,
		})
	}

	return hotArticles, nil
}
