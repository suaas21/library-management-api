package model

import (
	"time"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Mail     string `json:"mail"`
	Password string `json:"password"`
	PhoneNo  string `json:"phone_no"`
	// type of user is handled by UserType : `Admin` and `normal member`
	UserType string `json:"user_type"`
}

type UserDB struct {
	ID       int    `xorm:"pk autoincr id"`
	Name     string `xorm:"name"`
	Mail     string `xorm:"mail"`
	Password string `xorm:"password"`
	PhoneNo  string `xorm:"phone_no"`
	UserType string `xorm:"user_type"`

	CreatedAt time.Time `xorm:"created" `
	UpdatedAt time.Time `xorm:"updated"`
}

func APIFormat(u UserDB) *User {
	return &User{
		ID:       u.ID,
		Name:     u.Name,
		Mail:     u.Mail,
		PhoneNo:  u.PhoneNo,
		UserType: u.UserType,
		Password: u.Password,
	}
}
func UserDBFormat(u User) UserDB {
	return UserDB{
		ID:       u.ID,
		Name:     u.Name,
		Mail:     u.Mail,
		Password: u.Password,
		PhoneNo:  u.PhoneNo,
		UserType: u.UserType,
	}

}

type UsersDB struct {
	Users []UserDB
}

func (UserDB) TableName() string {
	return "users"
}
func (UsersDB) TableName() string {
	return "users"
}
