package entity

import (
	"oprec/go-blog/utils"

	"gorm.io/gorm"
)

type User struct {
	ID       uint64 `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	ListBlog []Blog `json:"list_blog,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	var err error
	u.Password, err = utils.HashAndSalt(u.Password)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	var err error
	if u.Password != "" {
		u.Password, err = utils.HashAndSalt(u.Password)
	}
	if err != nil {
		return err
	}
	return nil
}
