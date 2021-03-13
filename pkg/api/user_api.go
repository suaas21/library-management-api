package api

import (
	"fmt"
	"net/http"
	"strconv"

	"gopkg.in/macaron.v1"

	"github.com/suaas21/library-management-api/pkg/middleware"

	"github.com/suaas21/library-management-api/model"

	"github.com/suaas21/library-management-api/pkg/db"
)

var Users []model.User

func Register(ctx *macaron.Context, user model.User) {
	result, err := db.CreateUser(model.UserDBFormat(user))
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

	userResult, err := db.GetUserInfo(userId)
	if userResult == nil || err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid user profile")
		return
	}

	ctx.JSON(http.StatusOK, userResult)
	return
}

func UpdateUserProfile(ctx *macaron.Context, user model.User) {
	currentUserType := ctx.Req.Header.Get("current_user_type")
	currentUserMail := ctx.Req.Header.Get("current_user_mail")
	if currentUserType != "user" {
		ctx.JSON(http.StatusNotAcceptable, "type of user didn't match")
		return
	}
	if currentUserMail != "" {
		user.Mail = currentUserMail
		resultUser, err := db.UpdateUserProfile(model.UserDBFormat(user))
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

func Login(ctx *macaron.Context, user model.User) {
	userLoginInfo, err := db.GetUserLoginInfo(model.UserDBFormat(user))
	if err != nil {
		ctx.JSON(http.StatusNotFound, err.Error())
		return
	}
	if userLoginInfo == nil {
		ctx.JSON(http.StatusNotFound, "User credential not fount in database")
		return
	}

	tokenString, err := middleware.GenerateJWT(userLoginInfo.Mail, userLoginInfo.UserType, userLoginInfo.ID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, tokenString)
	return
}
