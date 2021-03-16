package database

import (
	"errors"
	"fmt"
)

type BookLoanRequest struct {
	Id         int       `xorm:"id pk autoincr" json:"id"`
	UserId     int       `xorm:"user_id" json:"user_id"`
	BookId     int       `xorm:"book_id" json:"book_id"`
	Status     string    `xorm:"status" json:"status"`
}

type BookLoanRequests struct {
	BookRequests []BookLoanRequest
}

func (BookLoanRequest) TableName() string {
	return "book_loan_requests"
}

func (BookLoanRequests) TableName() string {
	return "book_loan_requests"
}

func AddBookRequestToDB(bookLoanRequest BookLoanRequest) (*BookLoanRequest, error) {
	fmt.Println(bookLoanRequest)
	_, err := eng.Insert(bookLoanRequest)
	if err != nil {
		return nil, err
	}
	_, err = eng.Where("id=?", bookLoanRequest.Id).Get(&bookLoanRequest)
	if err != nil {
		return nil, err
	}

	return &bookLoanRequest, nil
}

func ShowBookRequestsFromDB() (*BookLoanRequests, error) {
	bookLoanRequestsdb := BookLoanRequests{}
	var bookLoanRequests []BookLoanRequest
	err := eng.Find(&bookLoanRequests)
	if err != nil {
		return nil, err
	}
	bookLoanRequestsdb.BookRequests = append(bookLoanRequestsdb.BookRequests, bookLoanRequests...)
	return &bookLoanRequestsdb, nil
}

func ShowBookRequestByIdFromDB(requestId int) (*BookLoanRequest, error) {
	var bookLoanRequest BookLoanRequest
	okk, _ := eng.Where("id=?", requestId).Get(&bookLoanRequest)
	if okk {
		return &bookLoanRequest, nil
	}
	return nil, errors.New("the book request didn't find")
}

func UpdateBookRequestToDB(userId, bookId int) (*BookLoanRequest, error) {
	var retrievedBook Book
	var retrievedUser User
	bookLoanRequest := BookLoanRequest{
		UserId: userId,
		BookId: bookId,
	}
	// 1st find the book and user from database using bookId and and userId
	isBookId, err := eng.Where("id=?", bookId).Get(&retrievedBook)
	if err != nil {
		return nil, err
	}
	isUserId, err := eng.Where("id=?", userId).Get(&retrievedUser)
	if err != nil {
		return nil, err
	}
	if isBookId && isUserId {
		// update book book request status
		if retrievedBook.NotAvailable == false {
			bookLoanRequest.Status = "Accepted"
		} else {
			bookLoanRequest.Status = "Rejected"
		}
		_, err := eng.Update(&bookLoanRequest)
		if err != nil {
			return nil, err
		}
		return &bookLoanRequest, nil
	}
	return nil, fmt.Errorf("book/user data not found in db")
}

func DeleteBookRequestFromDB(requestId int) (bool, error) {
	session := eng.NewSession()
	defer session.Close()
	err := session.Begin()

	ok, err := eng.Id(requestId).Delete(&BookLoanRequest{})
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
