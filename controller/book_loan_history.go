package controller

import (
	"net/http"

	"gopkg.in/macaron.v1"

	"github.com/suaas21/library-management-api/model"
)

func (c Controller) AddBookLoan(ctx *macaron.Context, bookHistory model.BookLoanHistory) {
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
	bookLoan, err := c.AddBookLoanToDB(userId, bookId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, bookLoan)
}

func (c Controller) ShowLoanHistory(ctx *macaron.Context) {
	loanHistories, err := c.ShowBookLoanHistories()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, loanHistories)
	return
}

func (c Controller) ReturnBook(ctx *macaron.Context, bookHistory model.BookLoanHistory) {
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

	result, err := c.UpdateBookLoanHistory(userId, bookId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
	}

	ctx.JSON(http.StatusCreated, result)
}
