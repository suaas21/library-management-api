package model

import (
	"time"
)

type User struct {
	Id       int    `xorm:"pk autoincr id" json:"id"`
	Name     string `xorm:"name" json:"name"`
	Mail     string `xorm:"mail" json:"mail"`
	Password string `xorm:"password" json:"password"`
	PhoneNo  string `xorm:"phone_no" json:"phone_no"`
	// type of user is handled by UserType : `admin` and `member`
	UserType string `xorm:"user_type" json:"user_type"`

	CreatedAt time.Time `xorm:"created" json:"created_at" `
	UpdatedAt time.Time `xorm:"updated" json:"updated_at" `
}

type Users struct {
	Users []User
}

func (User) TableName() string {
	return "users"
}
func (Users) TableName() string {
	return "users"
}
