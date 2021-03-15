package model

import "time"

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
