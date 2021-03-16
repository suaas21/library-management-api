package database

import (
	"errors"
	"fmt"
	"time"
)

type BookLoanHistory struct {
	Id     int `xorm:"pk autoincr id" json:"id"`
	BookId int `xorm:"book_id" json:"book_id"`
	UserId int `xorm:"user_id" json:"user_id"`

	Returned      bool   `xorm:"returned DEFAULT FALSE" json:"returned"`
	PurchasedDate string `xorm:"created" json:"purchased_date"`
	ReturnDate    string `xorm:"update updated " json:"return_date"`
}

type BookLoanHistories struct {
	BookLoanHistories []BookLoanHistory
}

func (BookLoanHistories) TableName() string {
	return "book_loan_histories"
}

func (BookLoanHistory) TableName() string {
	return "book_loan_histories"
}

func AddBookLoanToDB(userId int, bookId int) (*BookLoanHistory, error) {
	var user User
	var book Book
	isUser, err := eng.ID(userId).Get(&user)
	if err != nil {
		return nil, err
	}

	isBook, err := eng.ID(bookId).Get(&book)
	if err != nil {
		return nil, err
	}

	if !(isUser && isBook) {
		return nil, fmt.Errorf("book/UserInfo is not available in store")
	}

	if book.NotAvailable {
		return nil, fmt.Errorf("book is not available in stack")
	}

	bookHistory := BookLoanHistory{
		UserId:   user.Id,
		BookId:   book.Id,
	}
	book.NotAvailable = true

	session := eng.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return nil, err
	}

	// update the book store
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
		return nil, errors.New("not found the bookHistory in store")
	}

	return &bookHistory, nil
}

func UpdateBookLoanHistory(userId int, bookId int) (*BookLoanHistory, error) {
	var bookHistory BookLoanHistory
	// 1st find the bookHistory from database using userId and bookId and returned = FALSE
	ok, err := eng.Where("book_id=? AND user_id=? AND returned = FALSE", bookId, userId).Get(&bookHistory)
	if err != nil {
		return nil, err
	}
	if ok {
		var book Book
		_, err := eng.Id(bookId).Get(&book)
		if err != nil {
			// no book found in store using bookId
			return nil, err
		}
		// book found in store, so update the store
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

		// bookHistory found in store, so update the store for specific bookHistory
		bookHistory.Returned = true
		bookHistory.ReturnDate = time.Now().String()
		_, err = eng.Id(bookHistory.Id).UseBool().Update(&bookHistory)
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

func ShowBookLoanHistories() (*BookLoanHistories, error) {
	bookLoanHistories := BookLoanHistories{}
	var bookHistories []BookLoanHistory
	err := eng.Find(&bookHistories)
	if err != nil {
		return nil, err
	}
	bookLoanHistories.BookLoanHistories = append(bookLoanHistories.BookLoanHistories, bookHistories...)
	return &bookLoanHistories, nil
}
