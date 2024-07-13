package dao

import (
	"github.com/go-programming-tour-book/blog-service/internal/model"
	"github.com/go-programming-tour-book/blog-service/pkg/app"
)

func (d *Dao) CountUser(username string, email string, state uint8) (int, error) {
	user := model.User{Username: username, Email: email, State: state}
	return user.Count(d.engine)
}

func (d *Dao) ListUser(username string, email string, state uint8, page, pageSize int) ([]*model.User, error) {
	user := model.User{Username: username, Email: email, State: state}
	pageOffset := app.GetPageOffset(page, pageSize)
	return user.List(d.engine, pageOffset, pageSize)
}

func (d *Dao) RegisterUser(username string, email string, password string, state uint8) error {
	user := model.User{Username: username, Password: password, Email: email, State: state}
	return user.Create(d.engine)
}
