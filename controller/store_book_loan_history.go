package controller

import (
	"errors"
	"fmt"
	"github.com/suaas21/library-management-api/model"
	"time"
)

func (c Controller) AddBookLoanToDB(userId int, bookId int) (*model.BookLoanHistory, error) {
	var user model.User
	var book model.Book
	isUser, err := c.Eng.ID(userId).Get(&user)
	if err != nil {
		return nil, err
	}

	isBook, err := c.Eng.ID(bookId).Get(&book)
	if err != nil {
		return nil, err
	}

	if !(isUser && isBook) {
		return nil, fmt.Errorf("book/UserInfo is not available in store")
	}

	if book.NotAvailable {
		return nil, fmt.Errorf("book is not available in stack")
	}

	bookHistory := model.BookLoanHistory{
		UserId:   user.Id,
		BookId:   book.Id,
	}
	book.NotAvailable = true

	session := c.Eng.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return nil, err
	}

	// update the book store
	_, err = c.Eng.Id(bookId).UseBool().Update(&book)
	if err != nil {
		session.Rollback()
		return nil, err
	}

	// insert bookHistory
	_, err = c.Eng.Insert(bookHistory)
	if err != nil {
		session.Rollback()
		return nil, err
	}
	err = session.Commit()
	if err != nil {
		return nil, errors.New("server failed")
	}

	// get the updated book history
	_, err = c.Eng.Where("user_id =? AND book_id = ? AND returned= FALSE", userId, bookId).Get(&bookHistory)
	if err != nil {
		return nil, errors.New("not found the bookHistory in store")
	}

	return &bookHistory, nil
}

func (c Controller) UpdateBookLoanHistory(userId int, bookId int) (*model.BookLoanHistory, error) {
	var bookHistory model.BookLoanHistory
	// 1st find the bookHistory from database using userId and bookId and returned = FALSE
	ok, err := c.Eng.Where("book_id=? AND user_id=? AND returned = FALSE", bookId, userId).Get(&bookHistory)
	if err != nil {
		return nil, err
	}
	if ok {
		var book model.Book
		_, err := c.Eng.Id(bookId).Get(&book)
		if err != nil {
			// no book found in store using bookId
			return nil, err
		}
		// book found in store, so update the store
		book.NotAvailable = false
		_, err = c.Eng.Id(bookId).UseBool().Update(book)
		if err != nil {
			return nil, err
		}

		// now the find the bookHistory from database using userId and bookId
		_, err = c.Eng.Where("book_id=? AND user_id=?", bookId, userId).Get(&bookHistory)
		if err != nil {
			return nil, err
		}

		// bookHistory found in store, so update the store for specific bookHistory
		bookHistory.Returned = true
		bookHistory.ReturnDate = time.Now().String()
		_, err = c.Eng.Id(bookHistory.Id).UseBool().Update(&bookHistory)
		if err != nil {
			return nil, err
		}

		// update the full bookHistory
		_, err = c.Eng.Where("book_id=? AND user_id=?", bookId, userId).UseBool().Update(&bookHistory)
		if err != nil {
			return nil, err
		}

		// get the updated bookHistory
		_, err = c.Eng.Where("user_id =? AND book_id = ?", userId, bookId).Get(&bookHistory)
		if err != nil {
			return nil, err
		}

		return &bookHistory, nil

	}
	return nil, fmt.Errorf("no returned book data found")
}

func (c Controller) ShowBookLoanHistories() (*model.BookLoanHistories, error) {
	bookLoanHistories := model.BookLoanHistories{}
	var bookHistories []model.BookLoanHistory
	err := c.Eng.Find(&bookHistories)
	if err != nil {
		return nil, err
	}
	bookLoanHistories.BookLoanHistories = append(bookLoanHistories.BookLoanHistories, bookHistories...)
	return &bookLoanHistories, nil
}
