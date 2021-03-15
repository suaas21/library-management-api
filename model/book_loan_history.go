package model

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
