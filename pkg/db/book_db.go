package db

import (
	"errors"
	"github.com/suaas21/library-management-api/model"
)

func AddBook(book model.Book) (*model.Book, error) {
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

func ShowAllBook() (*model.Books, error) {
	booksdb := model.Books{}
	var books []model.Book
	err := eng.Find(&books)
	if err != nil {
		return nil, err
	}
	booksdb.Books = append(booksdb.Books, books...)
	return &booksdb, nil
}

func ShowBook(bookId int) (*model.Book, error) {
	var book model.Book
	okk, _ := eng.Where("id=?", bookId).Get(&book)
	if okk {
		return &book, nil
	}
	return nil, errors.New("the book didn't find")
}

func DeleteBookMethod(bookId int) (bool, error) {
	session := eng.NewSession()
	defer session.Close()
	err := session.Begin()

	ok, err := eng.Id(bookId).Delete(&model.Book{})
	if err != nil {
		session.Rollback()
		return false, errors.New("roll backed")
	}
	_, err = eng.Where("book_id =?", bookId).Delete(&model.BookHistory{})
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

func ShowLoanHistory() (*model.AllBookHistory, error) {
	AllbookHistory := model.AllBookHistory{}
	var bookHistories []model.BookHistory
	err := eng.Find(&bookHistories)
	if err != nil {
		return nil, err
	}
	AllbookHistory.BookHistory = append(AllbookHistory.BookHistory, bookHistories...)
	return &AllbookHistory, nil
}

