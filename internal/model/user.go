package model

import (
	"errors"

	"github.com/jinzhu/gorm"
)

var ErrUsernameAlreadyExists = errors.New("username already exists")

type User struct {
	*Model
	Username string `json:"username" gorm:"uniqueIndex;not null"`
	Password string `json:"password"`
	Email    string `json:"email"`
	State    uint8  `json:"state"`
}

func (u User) TableName() string {
	return "blog_user"
}

// 添加检查用户名是否存在的方法
func (u User) IsUsernameExist(db *gorm.DB, username string) (bool, error) {
	var count int64
	err := db.Model(&User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
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
	exist, err := u.IsUsernameExist(db, u.Username)
	if err != nil {
		return err

	}
	if exist {
		return ErrUsernameAlreadyExists
	}

	return db.Create(&u).Error
}

func (u User) Update(db *gorm.DB, values interface{}) error {
	if err := db.Model(&u).Updates(values).Where("id = ? AND is_del = ?", u.ID, 0).Error; err != nil {
		return err
	}

	return nil
}

func (u User) Delete(db *gorm.DB) error {
	return db.Where("id = ? AND is_del = ?", u.Model.ID, 0).Delete(&u).Error
}
