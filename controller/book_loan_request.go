package controller

import (
	"github.com/suaas21/library-management-api/database"
	"gopkg.in/macaron.v1"
	"net/http"
	"strconv"
)

func AddBookRequest(ctx *macaron.Context, bookRequest database.BookLoanRequest) {
	currentUserType := ctx.Req.Header.Get("current_user_type")
	if currentUserType != "user" {
		ctx.JSON(http.StatusUnauthorized, "user type didn't match")
		return
	}
	requestedBook, err := database.AddBookRequestToDB(bookRequest)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, requestedBook)
}

func ShowBookRequests(ctx *macaron.Context) {
	resultBooks, err := database.ShowBookRequestsFromDB()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resultBooks)
	return
}

func ShowBookRequestById(ctx *macaron.Context) {
	key := ctx.Params(":id")
	requestId, err := strconv.Atoi(key)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}

	resultBook, err := database.ShowBookRequestByIdFromDB(requestId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resultBook)
	return
}

func UpdateBookRequest(ctx *macaron.Context, bookRequest database.BookLoanRequest) {
	currentUserType := ctx.Req.Header.Get("current_user_type")
	if currentUserType != "admin" {
		ctx.JSON(http.StatusBadGateway, "user is not admin")
		return
	}
	bookId := bookRequest.BookId
	userId := bookRequest.UserId
	if userId <= 0 || bookId <= 0 {
		ctx.JSON(http.StatusBadGateway, "invalid user/book id, must provide the valid id")
		return
	}

	result, err := database.UpdateBookRequestToDB(userId, bookId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
	}

	ctx.JSON(http.StatusCreated, result)
}

func DeleteBookRequest(ctx *macaron.Context) {
	currentUserType := ctx.Req.Header.Get("current_user_type")
	if currentUserType != "admin" {
		ctx.JSON(http.StatusBadGateway, "user is not admin")
		return
	}
	bookRequestId, err := strconv.Atoi(ctx.Params(":id"))
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}
	_, err = database.DeleteBookRequestFromDB(bookRequestId)
	if err == nil {
		ctx.JSON(http.StatusResetContent, "book request has been deleted")
		return
	}

	ctx.JSON(http.StatusBadGateway, err.Error())
	return
}