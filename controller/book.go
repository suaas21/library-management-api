package controller

import (
	"github.com/suaas21/library-management-api/database"
	"net/http"
	"strconv"

	"gopkg.in/macaron.v1"
)

func AddBook(ctx *macaron.Context, book database.Book) {
	currentUserType := ctx.Req.Header.Get("current_user_type")
	if currentUserType != "admin" {
		ctx.JSON(http.StatusBadGateway, "user is not admin")
		return
	}

	resultBook, err := database.AddBookToDB(book)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, resultBook)
}

func ShowBooks(ctx *macaron.Context) {
	condition := ctx.Query("author")
	resultBooks, err := database.ShowBooksFromDB(condition)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resultBooks)
	return
}

func ShowBookById(ctx *macaron.Context) {
	key := ctx.Params(":bookId")
	bookId, err := strconv.Atoi(key)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}

	resultBook, err := database.ShowBookByIdFromDB(bookId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resultBook)
	return
}

func UpdateBook(ctx *macaron.Context, book database.Book) {
	currentUserType := ctx.Req.Header.Get("current_user_type")
	if currentUserType != "admin" {
		ctx.JSON(http.StatusNotAcceptable, "type of user didn't match")
		return
	}

	bookName := book.BookName
	if bookName == "" {
		ctx.JSON(http.StatusBadRequest, "no book name is properly specified")
		return
	}

	result, err := database.UpdateBookToDB(book)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
	}

	ctx.JSON(http.StatusCreated, result)
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

	_, err = database.DeleteBookFromDB(bookId)
	if err == nil {
		ctx.JSON(http.StatusResetContent, "book has been deleted")
		return
	}

	ctx.JSON(http.StatusBadGateway, err.Error())
	return
}