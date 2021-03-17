package database

import (
	"errors"
	"fmt"
	"time"
)

type Book struct {
	Id           int       `xorm:"id pk autoincr" json:"id"`
	BookName     string    `xorm:"book_name unique" json:"book_name"`
	Author       string    `xorm:"author" json:"author"`
	NotAvailable bool      `xorm:"not_available DEFAULT FALSE" json:"not_available"`
	CreatedAt    time.Time `xorm:"created" json:"created_at"`
}

func (Book) TableName() string {
	return "books"
}

func AddBookToDB(book Book) (*Book, error) {
	// insert new book in the databases
	_, err := eng.Insert(book)
	if err != nil {
		return nil, err
	}
	_, err = eng.Where("book_name=?", book.BookName).Get(&book)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func ShowBooksFromDB(condition string) ([]Book, error) {
	// get all book from database in two away
	// 1. searching by author , if getting `condition` value
	// 2. no filtering, if getting `condition` empty
	var books []Book
	if condition == "" {
		err := eng.Find(&books)
		if err != nil {
			return nil, err
		}
	} else {
		err := eng.Where("author=?", condition).Find(&books)
		if err != nil {
			return nil, err
		}
	}
	return books, nil
}

func ShowBookByIdFromDB(bookId int) (*Book, error) {
	// get all book by book id
	var book Book
	okk, _ := eng.Where("id=?", bookId).Get(&book)
	if okk {
		return &book, nil
	}
	return nil, errors.New("the book didn't find")
}

func DeleteBookFromDB(bookId int) (bool, error) {
	// delete a specific book by using book id
	session := eng.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return false, err
	}

	ok, err := session.Id(bookId).Delete(&Book{})
	if err != nil {
		session.Rollback()
		return false, errors.New("roll backed")
	}
	_, err = eng.Where("book_id =?", bookId).Delete(&BookLoanHistory{})
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

	return false, errors.New("no book found last one")
}

func UpdateBookToDB(book Book) (*Book, error) {
	// update book author by using book name
	// book name is unique, so we are not permitted to update it
	var retrievedBook Book
	// 1st find the book from database using bookId and returned = FALSE
	ok, err := eng.Where("book_name=?", book.BookName).Get(&retrievedBook)
	if err != nil {
		return nil, err
	}
	if ok {
		// update book author
		if book.Author != "" {
			retrievedBook.Author = book.Author
		}
		_, err := eng.ID(retrievedBook.Id).Update(retrievedBook)
		if err != nil {
			return nil, err
		}
		return &retrievedBook, nil
	}
	return nil, fmt.Errorf("no book found data by using name")
}
