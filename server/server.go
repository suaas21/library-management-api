package server

import (
	"fmt"
	"github.com/go-macaron/binding"
	"github.com/suaas21/library-management-api/controller"
	"github.com/suaas21/library-management-api/controller/authentication"
	"github.com/suaas21/library-management-api/model"
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


func (svr Server) StartAPIServer() {
	svr.InitializeDB()

	m := macaron.Classic()

	m.Use(macaron.Renderer())
	m.Use(authentication.JwtMiddleWare)

	eng := controller.Controller{
		Eng: eng,
	}
	m.Get("/login", binding.Json(model.User{}), eng.Login)
	m.Post("/register", binding.Json(model.User{}),eng.Register)
	m.Get("/user-profile/:userId([0-9]+)", eng.UserProfile)
	m.Patch("/edit-profile", binding.Json(model.User{}), eng.UpdateUserProfile)

	m.Post("/loan-book", binding.Json(model.BookLoanHistory{}), eng.AddBookLoan)
	m.Get("/loan-history", binding.Json(model.BookLoanHistory{}), eng.ShowLoanHistory)
	m.Put("/return-book", binding.Json(model.BookLoanHistory{}), eng.ReturnBook)

	m.Post("/book", binding.Json(model.Book{}), eng.AddBook)
	m.Patch("/edit-book", binding.Json(model.Book{}), eng.UpdateBook)
	m.Get("/books", eng.ShowBooks)
	m.Get("/book/:bookId([0-9]+)", eng.ShowBookById)
	m.Delete("/delete-book/:bookId([0-9]+)", eng.DeleteBook)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", svr.ServerPort), m))
}