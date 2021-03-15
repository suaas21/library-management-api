package model

type BookRequest struct {
	Id     int    `xorm:"id pk autoincr" json:"id"`
	UserId int    `xorm:"user_id" json:"user_id"`
	BookId int    `xorm:"book_id" json:"book_id"`
	Status string `xorm:"status DEFAULT Pending" json:"status"`
}

type BookRequests struct {
	BookRequests []BookRequest
}

func (BookRequest) TableName() string {
	return "book_requests"
}

func (BookRequests) TableName() string {
	return "book_requests"
}
