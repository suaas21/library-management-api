package database

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

func (User) TableName() string {
	return "users"
}

func CreateUser(user User) (*User, error) {
	// create a user if not present in the databases
	var retrievedUser User
	ok, err := eng.Where("mail= ?", user.Mail).Get(&retrievedUser)
	if !ok || err != nil {
		// user not found in store, create a new one
		_, err := eng.Insert(user)
		if err != nil {
			return nil, err
		}
		ok, err := eng.Where("mail=?", user.Mail).Get(&retrievedUser)
		if !ok || err != nil {
			return nil, err
		}
		return &retrievedUser, nil
	}
	// user already registered in store, no need to create, just return
	return nil, nil
}

func GetUserInfo(userId int) (*User, error) {
	// get user information from databases
	var user User
	ok, err := eng.Where("id= ?", userId).Get(&user)
	if !ok || err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserLoginInfo(user User) (*User, error) {
	// get user login info from databases
	var retrievedUser User
	ok, err := eng.Where("mail = ? AND password = ?", user.Mail, user.Password).Get(&retrievedUser)
	if !ok || err != nil {
		return nil, err
	}
	// user found in store, return the result
	return &retrievedUser, nil
}

func UpdateUserProfileToDB(user User) (*User, error) {
	// update user profile info into databases
	// user can update her/his profile like name, phone number or password.
	// mail is unique nad not updatable.
	var retrievedUser User
	okk, err := eng.Where("mail = ? ", user.Mail).Get(&retrievedUser)
	if err != nil {
		return nil, err
	}
	if okk {
		// update userInfo in store
		if user.Name != "" {
			retrievedUser.Name = user.Name
		}
		if user.PhoneNo != "" {
			retrievedUser.PhoneNo = user.PhoneNo
		}
		if user.Password != "" {
			retrievedUser.Password = user.Password
		}
		_, err := eng.ID(retrievedUser.Id).Update(retrievedUser)
		if err != nil {
			return nil, err
		}
	}

	return &retrievedUser, nil
}
