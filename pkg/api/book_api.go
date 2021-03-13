package api

import (
	"net/http"
	"strconv"

	"github.com/suaas21/library-management-api/model"
	"github.com/suaas21/library-management-api/pkg/db"
	"gopkg.in/macaron.v1"
)

func AddNewBook(ctx *macaron.Context, book model.Book) {
	currentUserType := ctx.Req.Header.Get("current_user_type")
	if currentUserType != "admin" {
		ctx.JSON(http.StatusBadGateway, "user is not admin")
		return
	}

	resultBook, err := db.AddBook(book)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, resultBook)
}

func ShowAllBook(ctx *macaron.Context) {
	resultBooks, err := db.ShowAllBook()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resultBooks)
	return
}

func ShowBook(ctx *macaron.Context) {
	key := ctx.Params(":bookId")
	bookId, err := strconv.Atoi(key)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}

	resultBook, err := db.ShowBook(bookId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resultBook)
	return
}

func DeleteBook(ctx *macaron.Context) {
	currentUserType := ctx.Req.Header.Get("current_user_type")
	if currentUserType != "admin" {
		ctx.JSON(http.StatusBadRequest, "user is not admin")
		return
	}

	bookId, err := strconv.Atoi(ctx.Params(":bookId"))
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}

	_, err = db.DeleteBookMethod(bookId)
	if err == nil {
		ctx.JSON(http.StatusResetContent, "book has been deleted")
		return
	}

	ctx.JSON(http.StatusBadGateway, err.Error())
	return
}
