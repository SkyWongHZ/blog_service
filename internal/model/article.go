package model

import (
	"github.com/go-programming-tour-book/blog-service/pkg/redis"
	"github.com/jinzhu/gorm"
)

type Article struct {
	*Model
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	State         uint8  `json:"state"`
}

func (a Article) TableName() string {
	return "blog_article"
}

func (a Article) Create(db *gorm.DB) (*Article, error) {
	if err := db.Create(&a).Error; err != nil {
		return nil, err
	}
	// 缓存最新文章id
	redis.CacheLatestArticleID(a.ID)

	return &a, nil
}

func (a Article) Update(db *gorm.DB, values interface{}) error {
	if err := db.Model(&a).Where("id = ? AND is_del = ?", a.ID, 0).Updates(values).Error; err != nil {
		return err
	}

	return nil
}

func (a Article) Get(db *gorm.DB) (Article, error) {
	var article Article
	db = db.Where("id = ? AND state = ? AND is_del = ?", a.ID, a.State, 0)
	err := db.First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return article, err
	}

	return article, nil
}

func (a Article) Delete(db *gorm.DB) error {
	if err := db.Where("id = ? AND is_del = ?", a.Model.ID, 0).Delete(&a).Error; err != nil {
		return err
	}

	return nil
}

func (a Article) GetHotArticles(db *gorm.DB) ([]*Article, error) {
	latestIDs, err := redis.GetLatestArticleIDs()
	if err != nil {
		return a.GetLatestArticlesFromDB(db)
	}

	// 从数据库获取完整文章信息
	var articles []*Article
	err = db.Where("id IN (?)", latestIDs).Find(&articles).Error
	return articles, err
}

func (a Article) GetLatestArticlesFromDB(db *gorm.DB) ([]*Article, error) {
	var articles []*Article
	// 按创建时间降序
	err := db.Order("created_at DESC").Limit(2).Find(&articles).Error
	return articles, err
}

type ArticleRow struct {
	ArticleID     uint32
	TagID         uint32
	TagName       string
	ArticleTitle  string
	ArticleDesc   string
	CoverImageUrl string
	Content       string
}

// 根据标签id获取文章列表
func (a Article) ListByTagID(db *gorm.DB, tagID uint32, pageOffset, pageSize int) ([]*ArticleRow, error) {
	fields := []string{"ar.id AS article_id", "ar.title AS article_title", "ar.desc AS article_desc", "ar.cover_image_url", "ar.content"}
	fields = append(fields, []string{"t.id AS tag_id", "t.name AS tag_name"}...)
	query := db.Select(fields).Table(Article{}.TableName()+" AS ar").
		Joins("LEFT JOIN `"+ArticleTag{}.TableName()+"` AS at ON ar.id = at.article_id").
		Joins("LEFT JOIN `"+Tag{}.TableName()+"` AS t ON at.tag_id = t.id").
		Where("ar.state = ? AND ar.is_del = ?", a.State, 0)
	if tagID != 0 {
		query = query.Where("at.`tag_id` = ?", tagID)
	}
	if pageOffset >= 0 && pageSize > 0 {
		query = query.Offset(pageOffset).Limit(pageSize)
	}
	rows, err := query.Rows()
	// 结果处理
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*ArticleRow
	for rows.Next() {
		r := &ArticleRow{}
		if err := rows.Scan(&r.ArticleID, &r.ArticleTitle, &r.ArticleDesc, &r.CoverImageUrl, &r.Content, &r.TagID, &r.TagName); err != nil {
			return nil, err
		}

		articles = append(articles, r)
	}

	return articles, nil
}

// 根据标签id获取文章数量
func (a Article) CountByTagID(db *gorm.DB, tagID uint32) (int, error) {
	var count int64
	query := db.Table(Article{}.TableName()+" AS ar").
		Where("ar.state = ? AND ar.is_del = ?", a.State, 0)

	if tagID != 0 {
		query = query.Joins("INNER JOIN `"+ArticleTag{}.TableName()+"` AS at ON ar.id = at.article_id").
			Where("at.tag_id = ?", tagID)
	}

	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}

	return int(count), nil
}
