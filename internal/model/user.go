package model

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	*Model
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	State    uint8  `json:"state"`
}

func (u User) TableName() string {
	return "blog_user"
}

func (u User) Count(db *gorm.DB) (int, error) {
	var count int
	if u.Username != "" {
		db = db.Where("username = ?", u.Username)
	}
	if u.Email != "" {
		db = db.Where("email = ?", u.Email)
	}
	db = db.Where("state = ?", u.State)
	if err := db.Model(&u).Where("is_del = ?", 0).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (u User) List(db *gorm.DB, pageOffset, pageSize int) ([]*User, error) {
	var users []*User
	if pageOffset >= 0 && pageSize > 0 {
		db = db.Offset(pageOffset).Limit(pageSize)
	}
	if u.Username != "" {
		db = db.Where("username = ? ", u.Username)
	}
	if u.Email != "" {
		db = db.Where("email = ?", u.Email)
	}
	db = db.Where("state = ?", u.State)
	if err := db.Where("is_del = ? ", 0).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil

}

func (u User) Create(db *gorm.DB) error {
	return db.Create(&u).Error
}
