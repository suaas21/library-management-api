# REST API Server for Library Management

The purpose of the API is to provide a management system for a library. There are two types of users in a Library.
- `Admin` 
- `Member/User`

**Admin :** Admin have the permission to Create, Update, Remove books from library. S/he has also the permission of accept and reject book loan requested by the users. If the book loan is issued and returned then S/he update the database.

**Users :** Users have only the permission to view books details, book-loans and also search by an author.

## Prerequisites

- This api server has implemented by using `go-macaron`. That's why at first you need to have a knowledge about [go-macaron](https://github.com/go-macaron/docs/blob/master/starter_guide.md). Actually this library is compatible to `go http`.
- Then you need to have a knowledge of about [go-xrom](https://github.com/go-xorm/xorm)
- And also have a little bit of knowledge about [cobra flag](https://github.com/spf13/cobra)

## To Start API Server
**Clone repository and enter the working dir :**  

```console
$ git clone git@github.com:suaas21/library-management-api.git
$ cd library-management-api
```

**Import database postgress :**
```console
$ sudo -u postgres psql < database.sql
```

**Print server version and start the server :** 

```console
$ go install
$ library-management-api version
$ library-management-api start
```


## Data Model

There are three types of data model we have introduced to design the api. the following data model are given below : 

- User Model
``````
type User struct {
	Id       int    `xorm:"pk autoincr id" json:"id"`
	Name     string `xorm:"name" json:"name"`
	Mail     string `xorm:"mail" json:"mail"`
	Password string `xorm:"password" json:"password"`
	PhoneNo  string `xorm:"phone_no" json:"phone_no"`
	// type of user is handled by UserType : `admin` and `member`
	UserType string `xorm:"user_type" json:"user_type"`

	CreatedAt time.Time `xorm:"created" json:"created_at" `
	UpdatedAt time.Time `xorm:"updated" json:"updated_at" `
}
``````
- Book Model
``````
type Book struct {
	Id           int       `xorm:"id pk autoincr" json:"id"`
	BookName     string    `xorm:"book_name unique" json:"book_name"`
	Author       string    `xorm:"author" json:"author"`
	NotAvailable bool      `xorm:"not_available DEFAULT FALSE" json:"not_available"`
	CreatedAt    time.Time `xorm:"created" json:"created_at"`
}
``````
- Book Loan Model
``````
type BookLoanHistory struct {
	Id     int `xorm:"pk autoincr id" json:"id"`
	BookId int `xorm:"book_id" json:"book_id"`
	UserId int `xorm:"user_id" json:"user_id"`

	Returned      bool   `xorm:"returned DEFAULT FALSE" json:"returned"`
	PurchasedDate string `xorm:"created" json:"purchased_date"`
	ReturnDate    string `xorm:"update updated " json:"return_date"`
}
``````
- Book Loan Request Model
``````
type BookLoanRequest struct {
	Id         int       `xorm:"id pk autoincr" json:"id"`
	UserId     int       `xorm:"user_id" json:"user_id"`
	BookId     int       `xorm:"book_id" json:"book_id"`
	// requet status: `Pending`, `Rejected` or `Accepted` 
	Status     string    `xorm:"status" json:"status"`
}
``````

## Available API Endpoints

|  Method | API Endpoint  | Authentication Type | Access Permission | Description |
|---|---|---|---|---|
|POST| /register | No auth | Any type of user | Registration for user/member |
|GET| /login | Basic or Bearer token | Any type of user | Return jwt token in response for successful authentication |
|GET| /user-profile/{userId} | No auth | Any type of user | Return a specific user in response | 
|PATCH| /edit-profile | Bearer token | User/Member | Return the updated user profile data in response | 
|POST| /loan-book | Bearer token | Admin | Admin can issue a book loan for the users | 
|PUT| /returned-book | Bearer token | Admin | Admin can put the returned book | 
|GET| /loan-history | No Auth | Any type of User | User can view loan book details | 
|POST| /book | Bearer token | Admin | Admin can create a new book |
|PATCH| /edit-book | Bearer token | Admin | Admin can update book author |
|GET| /book/{bookId} | No Auth | Any type of user  | Users can view specific book |
|GET| /books | No Auth | Any type of user  | Users can view all listed book |
|DELETE| /delete-book/{id} | Bearer token | Admin | Delete the book data and returned the updated data in response | 
|POST| /request | Bearer token | User | User can request for book loan |
|GET| /requests | No Auth | Any type of user | User/Admin can get requested book loan info |
|PATCH| /edit-request | Bearer token | Admin | Only amin can edit the request for book loan(Accepted/Rejected) |
|DELETE| /delete-request | Bearer token | Admin | Only admin delete the request for book loan |
## Available Flags

| Flag | Shorthand | Default value | Example | Description
|---|---|---|---|---|
|server-port|-|4000|  library-management-api start --port=8090 | Start API server in the given port otherwise in default port|
|db-port|-|5432| library-management-api --db-port=5432 | database will be started in this port|
|db-password|-|pass| library-management-api --db-password=pass | database password |
|db-name|-|library_management| library-management-api --db-name=library_management | database name |
|db-user|-|postgres| library-management-api --db-user=postgres | database user |


## Some Sample Curl commands to the server

Initialize database 

```console
$ sudo -u postgres psql < database.sql
```

Run API server

```console
$ library-management-api start
``` 

Registration for Admin

```console
$ curl -X POST -v -H "Content-Type:application/json" -d '{"name":"azad","mail":"azad@gmail.com","password":"password","phone_no":"017771","user_type":"admin"}' http://localhost:4000/register
```

Registration for user

```console
$ curl -X POST -v -H "Content-Type:application/json" -d '{"name":"sagor","mail":"sagor@gmail.com","password":"password","phone_no":"017771","user_type":"user"}' http://localhost:4000/register
```

Login for Admin/User(you will get bearer token)

```console
$ curl -X GET -H "Content-Type:application/json" -d '{"mail":"sagor@gmail.com","password":"password"}' http://localhost:4000/login

$ curl -X GET -H "Content-Type:application/json" -d '{"mail":"azad@gmail.com","password":"password"}' http://localhost:4000/login
```

Update user profile

```console
$ curl -X PATCH -H "Authorization: Bearer <user bearer token>" -d '{"name":"sagor","mail":"sagor@gmail.com","password":"password","phone_no":"017771885","user_type":"user"}' http://localhost:4000/edit-profile

```

Get user profile with id 1

```console
$ curl -X GET http://localhost:4000/user-profile/1
``` 

Add some New book 

```console
$ curl -X POST -H "Authorization: Bearer <admin bearer token>" -H "Content-Type:application/json" -d '{"book_name":"hello world", "author":"sagor"}' http://localhost:4000/book

$ curl -X POST -H "Authorization: Bearer <admin bearer token>" -H "Content-Type:application/json" -d '{"book_name":"Azad life story", "author":"azad"}' http://localhost:4000/book
```

Update book

```console
$ curl -X PATCH -H "Authorization: Bearer <admin bearer token>" -d '{"book_name":"Azad life story","author":"arya azad"}' http://localhost:4000/edit-book
```

Get loan book

```console
$ curl -X POST -H "Authorization: Bearer <admin bearer token>" -d '{"user_id":2,"book_id":2}' http://localhost:4000/loan-book
```

Show book loan history

```console
$ curl -X GET http://localhost:4000/loan-history
```

Return book from loan

```console
$ curl -X PUT -H "Authorization: Bearer <admin bearer token>" -d '{"user_id":2,"book_id":2}' http://localhost:4000/return-book
```

Get all books info 

```console
$ curl -X GET  http://localhost:4000/books
``` 

Get a specific book with id 1

```console
$ curl -X GET http://localhost:4000/book/1
``` 

Request book loan

```console
$ curl -X POST -v -H "Content-Type:application/json" -H "Authorization: Bearer <user bearer token>" -d '{"user_id":1,"book_id":1}' http://localhost:4000/request
```

Show Requested book loan
```console
$ curl -X GET http://localhost:4000/requests
$ curl -X GET http://localhost:4000/request/2
```

Update Requested Book Loan(Accepted/Rejected depends on availability)

```console
$ curl -X PATCH -H -H "Content-Type:application/json" -H "Authorization: Bearer <admin bearer token>" -d '{"user_id":1,"book_id":1}' http://localhost:4000/edit-request
```

Delete Request Book Loan

```console
$ curl -X DELETE -H "Authorization: Bearer <admin bearer token>"  http://localhost:4000/delete-request/1
```

Delete book with given id

```console
$ curl -X DELETE -H "Authorization: Bearer <admin berear token>" http://localhost:4000/delete-book/1
``` 