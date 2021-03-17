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
	imageName, err := FileUpload(ctx)
	if err != nil {
		// don't return, because we ignore if the image is not upload
		fmt.Println("image not uploaded, because:", err.Error())
	}
	// create the user and store in database
	if imageName != "" {
		user.Image = imageName
	}
	result, err := database.CreateUser(user)
	if err != nil {
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

func ChangeUserImage(ctx *macaron.Context) {
	currentUserType := ctx.Req.Header.Get("current_user_type")
	if currentUserType != "user" {
		ctx.JSON(http.StatusNotAcceptable, "user is not authenticated, need bearer token to upload image")
		return
	}

	key := ctx.Params(":userId")
	userId, err := strconv.Atoi(key)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}

	imageName, err := FileUpload(ctx)
	//here we call the function we made to get the image and save it
	if err != nil {
		// don't return, because we ignore if the image is not upload
		fmt.Println("image not uploaded, because:", err.Error())
	}
    updatedUser, err := database.ChangeUserImageNameToDB(userId, imageName)
    if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, updatedUser)
	return
}
