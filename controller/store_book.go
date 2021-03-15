package controller

import (
	"errors"
	"fmt"
	"github.com/suaas21/library-management-api/model"
)

func (c Controller) AddBookToDB(book model.Book) (*model.Book, error) {
	_, err := c.Eng.Insert(book)
	if err != nil {
		return nil, err
	}
	_, err = c.Eng.Where("book_name=?", book.BookName).Get(&book)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (c Controller) ShowBooksFromDB() (*model.Books, error) {
	booksdb := model.Books{}
	var books []model.Book
	err := c.Eng.Find(&books)
	if err != nil {
		return nil, err
	}
	booksdb.Books = append(booksdb.Books, books...)
	return &booksdb, nil
}

func (c Controller) ShowBookByIdFromDB(bookId int) (*model.Book, error) {
	var book model.Book
	okk, _ := c.Eng.Where("id=?", bookId).Get(&book)
	if okk {
		return &book, nil
	}
	return nil, errors.New("the book didn't find")
}

func (c Controller) DeleteBookFromDB(bookId int) (bool, error) {
	session := c.Eng.NewSession()
	defer session.Close()
	err := session.Begin()

	ok, err := c.Eng.Id(bookId).Delete(&model.Book{})
	if err != nil {
		session.Rollback()
		return false, errors.New("roll backed")
	}
	_, err = c.Eng.Where("book_id =?", bookId).Delete(&model.BookLoanHistory{})
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

func (c Controller) UpdateBookToDB(book model.Book) (*model.Book, error) {
	var retrievedBook model.Book
	// 1st find the book from database using bookId and returned = FALSE
	ok, err := c.Eng.Where("book_name=?", book.BookName).Get(&retrievedBook)
	if err != nil {
		return nil, err
	}
	if ok {
		// update book author
		if book.Author != "" {
			retrievedBook.Author = book.Author
		}
		_, err := c.Eng.ID(retrievedBook.Id).Update(retrievedBook)
		if err != nil {
			return nil, err
		}
		return &retrievedBook, nil
	}
	return nil, fmt.Errorf("no book found data by using name")
}

