package controller

import (
	"github.com/suaas21/library-management-api/database"
	"net/http"

	"gopkg.in/macaron.v1"
)

func AddBookLoan(ctx *macaron.Context, bookHistory database.BookLoanHistory) {
	currentUserType := ctx.Req.Header.Get("current_user_type")
	if currentUserType != "admin" {
		ctx.JSON(http.StatusBadGateway, "only admin can update purchase book info")
		return
	}

	userId := bookHistory.UserId
	bookId := bookHistory.BookId
	if userId <= 0 || bookId <= 0 {
		ctx.JSON(http.StatusBadGateway, "invalid user/book id")
		return
	}
	bookLoan, err := database.AddBookLoanToDB(userId, bookId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, bookLoan)
}

func ShowLoanHistory(ctx *macaron.Context) {
	bookLoanHistories, err := database.ShowBookLoanHistories()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, bookLoanHistories)
	return
}

func ReturnBook(ctx *macaron.Context, bookHistory database.BookLoanHistory) {
	currentUserType := ctx.Req.Header.Get("current_user_type")
	if currentUserType != "admin" {
		ctx.JSON(http.StatusUnauthorized, "user type didn't match")
		return
	}

	userId := bookHistory.UserId
	bookId := bookHistory.BookId
	if userId <= 0 || bookId <= 0 {
		ctx.JSON(http.StatusBadGateway, "invalid user/book id")
		return
	}

	result, err := database.UpdateBookLoanHistory(userId, bookId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
	}

	ctx.JSON(http.StatusCreated, result)
}
