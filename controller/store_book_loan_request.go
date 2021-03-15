package controller

import (
	"errors"
	"fmt"
	"github.com/suaas21/library-management-api/model"
)

func (c Controller) AddBookRequestToDB(bookRequest model.BookRequest) (*model.BookRequest, error) {
	_, err := c.Eng.Insert(bookRequest)
	if err != nil {
		return nil, err
	}
	_, err = c.Eng.Where("id=?", bookRequest.Id).Get(&bookRequest)
	if err != nil {
		return nil, err
	}

	return &bookRequest, nil
}

func (c Controller) ShowBookRequestsFromDB() (*model.BookRequests, error) {
	bookRequestsdb := model.BookRequests{}
	var bookRequests []model.BookRequest
	err := c.Eng.Find(&bookRequests)
	if err != nil {
		return nil, err
	}
	bookRequestsdb.BookRequests = append(bookRequestsdb.BookRequests, bookRequests...)
	return &bookRequestsdb, nil
}

func (c Controller) ShowBookRequestByIdFromDB(requestId int) (*model.BookRequest, error) {
	var bookRequest model.BookRequest
	okk, _ := c.Eng.Where("id=?", requestId).Get(&bookRequest)
	if okk {
		return &bookRequest, nil
	}
	return nil, errors.New("the book request didn't find")
}

func (c Controller) UpdateBookRequestToDB(bookRequestId, userId, bookId int) (*model.BookRequest, error) {
	var retrievedBook model.Book
	var retrievedUser model.User
	bookRequest := model.BookRequest{
		Id:     bookRequestId,
		UserId: userId,
		BookId: bookId,
	}
	// 1st find the book and user from database using bookId and and userId
	isBookId, err := c.Eng.Where("id=?", bookId).Get(&retrievedBook)
	if err != nil {
		return nil, err
	}
	isUserId, err := c.Eng.Where("id=?", userId).Get(&retrievedUser)
	if err != nil {
		return nil, err
	}
	if isBookId && isUserId {
		// update book book request status
		if retrievedBook.NotAvailable == false {
			bookRequest.Status = "Accepted"
		} else {
			bookRequest.Status = "Rejected"
		}
		_, err := c.Eng.ID(bookRequestId).Update(&bookRequest)
		if err != nil {
			return nil, err
		}
		return &bookRequest, nil
	}
	return nil, fmt.Errorf("book/user data not found in db")
}

func (c Controller) DeleteBookRequestFromDB(bookRequestId int) (bool, error) {
	session := c.Eng.NewSession()
	defer session.Close()
	err := session.Begin()

	ok, err := c.Eng.Id(bookRequestId).Delete(&model.BookRequest{})
	if err != nil {
		session.Rollback()
		return false, errors.New("roll backed")
	}
	err = session.Commit()
	if err != nil {
		return false, errors.New("server failed")
	}
	if ok > 0 {
		return true, nil
	}

	return false, errors.New("no book request found")
}