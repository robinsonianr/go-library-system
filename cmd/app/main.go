package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ianrr/library/internal/api"
	"github.com/ianrr/library/internal/db"
)

func initServer() {
    database, err := db.NewDb()
    if err != nil {
        log.Fatal(err) 
    }
    
    bookRepo := &db.BookRepository{DB: database}

    bookHandler := &api.BookHandler{Repo: bookRepo, IsTest: true}

    
    r := chi.NewRouter()

    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })
    r.Get("/books", bookHandler.GetBooks)
    
    http.ListenAndServe(":8080", r)

}


func main() {
    //server.MountHandler()
    //server := app.CreateNewServer()
    log.Println("Hello, World!")
    initServer()
}
