package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ianrr/library/internal/api"
	"github.com/ianrr/library/internal/db"
	"github.com/ianrr/library/internal/db/user"
)

func initServer() {
    database, err := db.NewDb()
    if err != nil {
        log.Fatal(err) 
    }
    
    bookRepo := &db.BookRepository{DB: database}
    userRepo := &user.UserRepository{DB: database}

    bookHandler := &api.BookHandler{Repo: bookRepo, IsTest: false}
    userHandler := &api.UserHandler{Repo: userRepo}

    
    r := chi.NewRouter()
    r.Use(middleware.Logger)

    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })

    // methods of bookhandler as HTTP handlers
    r.Get("/books", bookHandler.GetBooks)
    r.Get("/book/{id}", bookHandler.GetBook)
    r.Post("/book", bookHandler.SubmitBook)

    // methods of userhandlers as HTTP handlers
    r.Post("/user", userHandler.RegisterUser)
    r.Post("/user/login", userHandler.LoginUser)
    
    http.ListenAndServe(":8080", r)
}


func main() {
    log.Println("Hello, World!")
    initServer()
}
