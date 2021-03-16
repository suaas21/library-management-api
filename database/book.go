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

type Books struct {
	Books []Book
}

func (Book) TableName() string {
	return "books"
}

func (Books) TableName() string {
	return "books"
}

func AddBookToDB(book Book) (*Book, error) {
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

func ShowBooksFromDB() (*Books, error) {
	booksdb := Books{}
	var books []Book
	err := eng.Find(&books)
	if err != nil {
		return nil, err
	}
	booksdb.Books = append(booksdb.Books, books...)
	return &booksdb, nil
}

func ShowBookByIdFromDB(bookId int) (*Book, error) {
	var book Book
	okk, _ := eng.Where("id=?", bookId).Get(&book)
	if okk {
		return &book, nil
	}
	return nil, errors.New("the book didn't find")
}

func DeleteBookFromDB(bookId int) (bool, error) {
	session := eng.NewSession()
	defer session.Close()
	err := session.Begin()

	ok, err := eng.Id(bookId).Delete(&Book{})
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
