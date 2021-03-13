package api

import (
	"net/http"

	"gopkg.in/macaron.v1"

	"github.com/suaas21/library-management-api/model"
	"github.com/suaas21/library-management-api/pkg/db"
)

func AddNewLoan(ctx *macaron.Context, bookHistory model.BookHistory) {
	currentUserType := ctx.Req.Header.Get("current_user_type")
	if currentUserType != "admin" {
		ctx.JSON(http.StatusBadGateway, "only admin can update purchase book info")
		return
	}

	userId := bookHistory.UserId
	bookId := bookHistory.BookId
	resultPurchase, err := db.AddNewLoan(userId, bookId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resultPurchase)
}

func ReturnBook(ctx *macaron.Context, bookHistory model.BookHistory) {
	currentUserType := ctx.Req.Header.Get("current_user_type")
	if currentUserType != "admin" {
		ctx.JSON(http.StatusUnauthorized, "user type didn't match")
		return
	}

	userId := bookHistory.UserId
	bookId := bookHistory.BookId
	result, err := db.UpdateBookHistory(userId, bookId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, "can't execute return request")
	}

	ctx.JSON(http.StatusCreated, result)
}
