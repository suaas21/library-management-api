package db

import (
	"github.com/suaas21/library-management-api/model"
)

func CreateUser(user model.UserDB) (*model.User, error) {
	var retrievedUser model.UserDB
	ok, err := eng.Where("mail= ?", user.Mail).Get(&retrievedUser)
	if !ok || err != nil {
		// user not found in db, create a new one
		_, err := eng.Insert(user)
		if err != nil {
			return nil, err
		}
		ok, err := eng.Where("mail=?", user.Mail).Get(&retrievedUser)
		if !ok || err != nil {
			return nil, err
		}
		return model.APIFormat(retrievedUser), nil
	}
   // user already registered in db, no need to create, just return
   return nil, nil
}

func GetUserInfo(userId int) (*model.User, error) {
	var user model.UserDB
	ok, err := eng.Where("id= ?", userId).Get(&user)
	if !ok || err != nil {
		return nil, err
	}
	return model.APIFormat(user), nil
}

func GetUserLoginInfo(user model.UserDB) (*model.User, error) {
	var retrievedUser model.UserDB
	ok, err := eng.Where("mail = ? AND password = ?", user.Mail, user.Password).Get(&retrievedUser)
	if !ok || err != nil {
		return nil, err
	}
	// user found in db, return the result
	return model.APIFormat(retrievedUser), nil
}

func UpdateUserProfile(user model.UserDB) (*model.User, error) {
	var retrievedUser model.UserDB
	okk, err := eng.Where("mail = ? ", user.Mail).Get(&retrievedUser)
	if err != nil {
		return nil, err
	}
	if okk {
		// update userInfo in db
		if user.Name != "" {
			retrievedUser.Name = user.Name
		}
		if user.PhoneNo != "" {
			retrievedUser.PhoneNo = user.PhoneNo
		}
		if user.Password != "" {
			retrievedUser.Password = user.Password
		}
		_, err := eng.ID(retrievedUser.ID).Update(retrievedUser)
		if err != nil {
			return nil, err
		}
	}

	return model.APIFormat(retrievedUser), nil
}
