package controller

import (
	"github.com/suaas21/library-management-api/model"
	"gopkg.in/macaron.v1"
	"net/http"
	"strconv"
)

func (c Controller) AddBookRequest(ctx *macaron.Context, bookRequest model.BookRequest) {
	requestedBook, err := c.AddBookRequestToDB(bookRequest)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, requestedBook)
}

func (c Controller) ShowBookRequests(ctx *macaron.Context) {
	resultBooks, err := c.ShowBookRequestsFromDB()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resultBooks)
	return
}

func (c Controller) ShowBookRequestById(ctx *macaron.Context) {
	key := ctx.Params(":id")
	requestId, err := strconv.Atoi(key)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}

	resultBook, err := c.ShowBookRequestByIdFromDB(requestId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resultBook)
	return
}

func (c Controller) UpdateBookRequest(ctx *macaron.Context, bookRequest model.BookRequest) {
	bookId := bookRequest.BookId
	userId := bookRequest.UserId
	if userId <= 0 || bookId <= 0 || bookRequest.Id <= 0 {
		ctx.JSON(http.StatusBadGateway, "invalid user/book id, must provide the valid id")
		return
	}

	result, err := c.UpdateBookRequestToDB(bookRequest.Id, userId, bookId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
	}

	ctx.JSON(http.StatusCreated, result)
}

func (c Controller) DeleteBookRequest(ctx *macaron.Context) {
	bookRequestId, err := strconv.Atoi(ctx.Params(":id"))
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}
	_, err = c.DeleteBookRequestFromDB(bookRequestId)
	if err == nil {
		ctx.JSON(http.StatusResetContent, "book request has been deleted")
		return
	}

	ctx.JSON(http.StatusBadGateway, err.Error())
	return
}