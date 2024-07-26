package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ianrr/library/internal/auth"
	"github.com/ianrr/library/internal/db"
	"github.com/ianrr/library/internal/utils"
)


type APIBook struct {
    ID      string `json:"id"`
    Title   string `json:"title"` 
    Author  string `json:"author"`
    Genre   string `json:"genre"`
    Pages   int    `json:"pages"` 
    Stock   int    `json:"stock"` 
}


type BookHandler struct {
    Repo    db.BookRepositoryInterface
    IsTest  bool
}

func checkAuth(w http.ResponseWriter, r *http.Request) {
    token := r.Header.Get("Authorization")

    if token == "" {
        http.Error(w, "Authorization token is required", http.StatusUnauthorized)
        return
    }

    secret, err := utils.ReadFile("../../privateKey.txt")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    isValid := auth.IsValid(token, secret)

    if !isValid {
        http.Error(w, "Invalid token", http.StatusUnauthorized)
        return
    }

}

func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request){
    if !h.IsTest {
        checkAuth(w, r) 
    }

    books, err := h.Repo.FindAllBooks()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(books); err != nil {
        http.Error(w, "Error encoding json", http.StatusInternalServerError)
        return
    }
}

func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
    if !h.IsTest {
        checkAuth(w, r)
    }

    id := r.PathValue("id")

    book, err := h.Repo.FindBookByID(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }


    if book == nil {
        http.NotFound(w, r)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(book); err != nil {
        http.Error(w, "Error encoding json", http.StatusInternalServerError)
        return
    }
}

func (h *BookHandler) SearchBook(w http.ResponseWriter, r *http.Request) {
    if !h.IsTest {
        checkAuth(w, r)
    }
 
    err := r.ParseMultipartForm(32 << 20)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    identifier := r.FormValue("identifier")

    book, err := h.Repo.FindBookByIdentifier(identifier)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if book == nil {
        http.NotFound(w, r)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(book); err != nil {
        http.Error(w, "Error encoding json", http.StatusInternalServerError)
        return
    }
}


func (h *BookHandler) SubmitBook(w http.ResponseWriter, r *http.Request) {

    var book APIBook
    err := json.NewDecoder(r.Body).Decode(&book)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if book.Title == "" || book.Author == "" || book.Genre == "" {
        http.Error(w, "Missing required field(s)", http.StatusBadRequest)
        return
    }

    err = h.Repo.RegisterBook(book.Title, book.Author, book.Genre, book.Pages)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

     log.Println("Successfuly added book")
     w.WriteHeader(http.StatusCreated)
}
