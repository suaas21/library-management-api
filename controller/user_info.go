package controller

import (
	"fmt"
	"github.com/suaas21/library-management-api/controller/authentication"
	"net/http"
	"strconv"

	"github.com/suaas21/library-management-api/model"
	"gopkg.in/macaron.v1"
)

var Users []model.User

func (c Controller) Register(ctx *macaron.Context, user model.User) {
	result, err := c.CreateUser(user)
	if result == nil || err != nil {
		ctx.JSON(http.StatusNotImplemented, fmt.Sprintf("the user already exist, err: %v", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, result)
}

func (c Controller) UserProfile(ctx *macaron.Context) {
	key := ctx.Params(":userId")
	userId, err := strconv.Atoi(key)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid user profile")
		return
	}

	userResult, err := c.GetUserInfo(userId)
	if userResult == nil || err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid user profile")
		return
	}

	ctx.JSON(http.StatusOK, userResult)
	return
}

func (c Controller) UpdateUserProfile(ctx *macaron.Context, user model.User) {
	currentUserType := ctx.Req.Header.Get("current_user_type")
	currentUserMail := ctx.Req.Header.Get("current_user_mail")
	if currentUserType != "user" {
		ctx.JSON(http.StatusNotAcceptable, "type of user didn't match")
		return
	}
	if currentUserMail != "" {
		user.Mail = currentUserMail
		resultUser, err := c.UpdateUserProfileToDB(user)
		if err != nil {
			ctx.JSON(http.StatusNotImplemented, "profile updating failed")
			return
		}
		ctx.JSON(http.StatusOK, resultUser)
		return
	}

	ctx.JSON(http.StatusNotAcceptable, "mail is not valid")
	return
}

func (c Controller) Login(ctx *macaron.Context, user model.User) {
	userLoginInfo, err := c.GetUserLoginInfo(user)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err.Error())
		return
	}
	if userLoginInfo == nil {
		ctx.JSON(http.StatusNotFound, "User credential not fount in database")
		return
	}

	tokenString, err := authentication.GenerateJWT(userLoginInfo.Mail, userLoginInfo.UserType, userLoginInfo.ID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, tokenString)
	return
}