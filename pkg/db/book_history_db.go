package db

import (
	"errors"
	"fmt"

	"github.com/suaas21/library-management-api/model"
)

func AddNewLoan(userId int, bookId int) (*model.BookHistory, error) {
	var user model.UserDB
	var book model.Book
	_, err := eng.ID(userId).Get(&user)
	if err != nil {
		return nil, err
	}
	_, err = eng.ID(bookId).Get(&book)
	if err != nil {
		return nil, err
	}
	if book.NotAvailable {
		return nil, fmt.Errorf("book is not available in stack")
	}

	bookHistory := model.BookHistory{
		UserId:   user.ID,
		UserName: user.Name,
		UserMail: user.Mail,
		BookId:   book.Id,
		BookName: book.BookName,
	}
	book.NotAvailable = true

	session := eng.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return nil, err
	}

	// update the book db
	_, err = eng.Id(bookId).UseBool().Update(&book)
	if err != nil {
		session.Rollback()
		return nil, err
	}

	// insert bookHistory
	_, err = eng.Insert(bookHistory)
	if err != nil {
		session.Rollback()
		return nil, err
	}
	err = session.Commit()
	if err != nil {
		return nil, errors.New("server failed")
	}

	// get the updated book history
	_, err = eng.Where("user_id =? AND book_id = ? AND returned= FALSE", userId, bookId).Get(&bookHistory)
	if err != nil {
		return nil, errors.New("not found the bookHistory in db")
	}

	return &bookHistory, nil
}

func UpdateBookHistory(userId int, bookId int) (*model.BookHistory, error) {
	var bookHistory model.BookHistory
	// 1st find the bookHistory from database using userId and bookId and returned = FALSE
	okk, err := eng.Where("book_id=? AND user_id=? AND returned = FALSE", bookId, userId).Get(&bookHistory)
	if err != nil {
		return nil, err
	}
	if okk {
		var book model.Book
		_, err := eng.Id(bookId).Get(&book)
		if err != nil {
			// no book found in db using bookId
			return nil, err
		}
		// book found in db, so update the db
		book.NotAvailable = false
		_, err = eng.Id(bookId).UseBool().Update(book)
		if err != nil {
			return nil, err
		}

		// now the find the bookHistory from database using userId and bookId
		_, err = eng.Where("book_id=? AND user_id=?", bookId, userId).Get(&bookHistory)
		if err != nil {
			return nil, err
		}

		// bookHistory found in db, so update the db for specific bookHistory
		bookHistory.Returned = true
		_, err = eng.Id(bookHistory.HistoryId).UseBool().Update(&bookHistory)
		if err != nil {
			return nil, err
		}

		// update the full bookHistory
		_, err = eng.Where("book_id=? AND user_id=?", bookId, userId).UseBool().Update(&bookHistory)
		if err != nil {
			return nil, err
		}

		// get the updated bookHistory
		_, err = eng.Where("user_id =? AND book_id = ?", userId, bookId).Get(&bookHistory)
		if err != nil {
			return nil, err
		}

		return &bookHistory, nil

	}
	return nil, fmt.Errorf("no returned book data found")
}
