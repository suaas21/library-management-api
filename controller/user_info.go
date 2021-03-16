package controller

import (
	"fmt"
	"github.com/suaas21/library-management-api/controller/authentication"
	"github.com/suaas21/library-management-api/database"
	"net/http"
	"strconv"

	"gopkg.in/macaron.v1"
)

func Register(ctx *macaron.Context, user database.User) {
	result, err := database.CreateUser(user)
	if result == nil || err != nil {
		ctx.JSON(http.StatusNotImplemented, fmt.Sprintf("the user already exist, err: %v", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, result)
}

func UserProfile(ctx *macaron.Context) {
	key := ctx.Params(":userId")
	userId, err := strconv.Atoi(key)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid user profile")
		return
	}

	userResult, err := database.GetUserInfo(userId)
	if userResult == nil || err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid user profile")
		return
	}

	ctx.JSON(http.StatusOK, userResult)
	return
}

func UpdateUserProfile(ctx *macaron.Context, user database.User) {
	currentUserType := ctx.Req.Header.Get("current_user_type")
	currentUserMail := ctx.Req.Header.Get("current_user_mail")
	if currentUserType != "user" {
		ctx.JSON(http.StatusNotAcceptable, "type of user didn't match")
		return
	}
	if currentUserMail != "" {
		user.Mail = currentUserMail
		resultUser, err := database.UpdateUserProfileToDB(user)
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

func Login(ctx *macaron.Context, user database.User) {
	userLoginInfo, err := database.GetUserLoginInfo(user)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err.Error())
		return
	}
	if userLoginInfo == nil {
		ctx.JSON(http.StatusNotFound, "User credential not fount in database")
		return
	}

	tokenString, err := authentication.GenerateJWT(userLoginInfo.Mail, userLoginInfo.UserType, userLoginInfo.Id)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, tokenString)
	return
}
