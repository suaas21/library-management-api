package server

import (
	"fmt"
	"github.com/go-macaron/binding"
	"github.com/suaas21/library-management-api/controller"
	"github.com/suaas21/library-management-api/controller/authentication"
	"github.com/suaas21/library-management-api/database"
	"gopkg.in/macaron.v1"
	"log"
	"net/http"
)

type Server struct {
	ServerPort string
	DBPort     string
	DBPassword string
	DBName     string
	DBUser     string
}


func StartAPIServer(port string) {
	database.InitializeDB()

	m := macaron.Classic()

	m.Use(macaron.Renderer())
	m.Use(authentication.JwtMiddleWare)

	m.Post("/register", binding.Json(database.User{}),controller.Register)
	m.Patch("/change-user-image/:userId", controller.ChangeUserImage)
	m.Get("/login", binding.Json(database.User{}), controller.Login)

	m.Get("/user-profile/:userId([0-9]+)", controller.UserProfile)
	m.Patch("/edit-profile", binding.Json(database.User{}), controller.UpdateUserProfile)

	m.Post("/book", binding.Json(database.Book{}), controller.AddBook)
	m.Patch("/edit-book", binding.Json(database.Book{}), controller.UpdateBook)
	m.Get("/books", controller.ShowBooks)
	m.Get("/book/:bookId([0-9]+)", controller.ShowBookById)
	m.Delete("/delete-book/:bookId([0-9]+)", controller.DeleteBook)

	m.Post("/request", binding.Json(database.BookLoanRequest{}), controller.AddBookRequest)
	m.Get("/requests", controller.ShowBookRequests)
	m.Get("/request/:id([0-9]+)", controller.ShowBookRequestById)
	m.Patch("/edit-request", binding.Json(database.BookLoanRequest{}), controller.UpdateBookRequest)
	m.Delete("/delete-request/:id([0-9]+)", controller.DeleteBookRequest)

	m.Post("/loan-book", binding.Json(database.BookLoanHistory{}), controller.AddBookLoan)
	m.Get("/loan-history", binding.Json(database.BookLoanHistory{}), controller.ShowLoanHistory)
	m.Put("/return-book", binding.Json(database.BookLoanHistory{}), controller.ReturnBook)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", port), m))
}