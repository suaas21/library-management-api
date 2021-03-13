package server

import (
	"fmt"
	"github.com/go-macaron/binding"
	"github.com/suaas21/library-management-api/model"
	"github.com/suaas21/library-management-api/pkg/api"
	"github.com/suaas21/library-management-api/pkg/middleware"
	"gopkg.in/macaron.v1"
	"log"
	"net/http"
)

func StartAPIServer(port string) {
	m := macaron.Classic()

	m.Use(macaron.Renderer())
	m.Use(middleware.JwtMiddleWare)

	m.Get("/login", binding.Json(model.User{}), api.Login)
	m.Post("/register", binding.Json(model.User{}), api.Register)
	m.Get("/user-profile/:userId([0-9]+)", api.UserProfile)
	m.Patch("/edit-profile", binding.Json(model.User{}), api.UpdateUserProfile)

	m.Post("/loan-book", binding.Json(model.BookHistory{}), api.AddNewLoan)
	m.Put("/return-book", binding.Json(model.BookHistory{}), api.ReturnBook)
	m.Post("/book", binding.Json(model.Book{}), api.AddNewBook)
	m.Get("/books", api.ShowAllBook)
	m.Get("/book/:bookId([0-9]+)", api.ShowBook)
	m.Delete("/delete-book/:bookId([0-9]+)", api.DeleteBook)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", port), m))
}