package controller

import (
	"github.com/go-xorm/xorm"
	"net/http"
	"strconv"

	"github.com/suaas21/library-management-api/model"
	"gopkg.in/macaron.v1"
)

type Controller struct {
	Eng *xorm.Engine
}

func (c Controller) AddBook(ctx *macaron.Context, book model.Book) {
	currentUserType := ctx.Req.Header.Get("current_user_type")
	if currentUserType != "admin" {
		ctx.JSON(http.StatusBadGateway, "user is not admin")
		return
	}

	resultBook, err := c.AddBookToDB(book)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, resultBook)
}

func (c Controller) ShowBooks(ctx *macaron.Context) {
	resultBooks, err := c.ShowBooksFromDB()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resultBooks)
	return
}

func (c Controller) ShowBookById(ctx *macaron.Context) {
	key := ctx.Params(":bookId")
	bookId, err := strconv.Atoi(key)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}

	resultBook, err := c.ShowBookByIdFromDB(bookId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resultBook)
	return
}

func (c Controller) UpdateBook(ctx *macaron.Context, book model.Book) {
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

	result, err := c.UpdateBookToDB(book)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
	}

	ctx.JSON(http.StatusCreated, result)
}

func (c Controller) DeleteBook(ctx *macaron.Context) {
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

	_, err = c.DeleteBookFromDB(bookId)
	if err == nil {
		ctx.JSON(http.StatusResetContent, "book has been deleted")
		return
	}

	ctx.JSON(http.StatusBadGateway, err.Error())
	return
}