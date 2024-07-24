package api

import (
	"encoding/json"
	"net/http"

	"github.com/ianrr/library/internal/auth"
	"github.com/ianrr/library/internal/db"
	"github.com/ianrr/library/internal/utils"
)


type APIBook struct {
    db.Book
    ID      string `json:"id"`
    Title   string `json:"title"` 
    Author  string `json:"author"`
    Genre   string `json:"genre"`
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

    secret, err := utils.ReadFile("../../privatekey.txt")
    if err != nil {
        http.Error(w, "Error reading token", http.StatusInternalServerError)
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
    id := r.FormValue("id")

    book, err := h.Repo.FindBookByID(id)
    if err != nil {
        http.Error(w, "Error retrieving book", http.StatusInternalServerError)
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(book); err != nil {
        http.Error(w, "Error encoding json", http.StatusInternalServerError)
        return
    }
}


// TODO: Implement Add book
func (h *BookHandler) AddBook(w http.ResponseWriter, r *http.Request) {
    if !h.IsTest {
        checkAuth(w, r)
    }
    
    var book APIBook
    err := json.NewDecoder(r.Body).Decode(&book)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

}
