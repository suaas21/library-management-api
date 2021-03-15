package controller

import (
	"github.com/suaas21/library-management-api/model"
)

func (c Controller) CreateUser(user model.User) (*model.User, error) {
	var retrievedUser model.User
	ok, err := c.Eng.Where("mail= ?", user.Mail).Get(&retrievedUser)
	if !ok || err != nil {
		// user not found in store, create a new one
		_, err := c.Eng.Insert(user)
		if err != nil {
			return nil, err
		}
		ok, err := c.Eng.Where("mail=?", user.Mail).Get(&retrievedUser)
		if !ok || err != nil {
			return nil, err
		}
		return &retrievedUser, nil
	}
   // user already registered in store, no need to create, just return
   return nil, nil
}

func (c Controller) GetUserInfo(userId int) (*model.User, error) {
	var user model.User
	ok, err := c.Eng.Where("id= ?", userId).Get(&user)
	if !ok || err != nil {
		return nil, err
	}
	return &user, nil
}

func (c Controller) GetUserLoginInfo(user model.User) (*model.User, error) {
	var retrievedUser model.User
	ok, err := c.Eng.Where("mail = ? AND password = ?", user.Mail, user.Password).Get(&retrievedUser)
	if !ok || err != nil {
		return nil, err
	}
	// user found in store, return the result
	return &retrievedUser, nil
}

func (c Controller) UpdateUserProfileToDB(user model.User) (*model.User, error) {
	var retrievedUser model.User
	okk, err := c.Eng.Where("mail = ? ", user.Mail).Get(&retrievedUser)
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
		_, err := c.Eng.ID(retrievedUser.ID).Update(retrievedUser)
		if err != nil {
			return nil, err
		}
	}

	return &retrievedUser, nil
}
