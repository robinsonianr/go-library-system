package db

import (
	"database/sql"
	"errors"
	"log"
)

var ErrNotFound = errors.New("book not found")

type BookRepository struct {
    DB *sql.DB
}

type Book struct {
    ID      string
    Title   string
    Author  string
    Genre   string
}

func (repo *BookRepository) FindAllBooks() ([]Book, error){
    rows, err := repo.DB.Query("SELECT id, title, author, genre FROM library_sys.Books")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var books []Book
    for rows.Next() {
        var b Book
        if err :=  rows.Scan(&b.ID, &b.Title, &b.Author, &b.Genre); err != nil {
            return nil, err
        }

        books = append(books, b)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return books, nil
}

func (repo *BookRepository) FindBookByID(id string) (*Book, error) {
    var b Book
    err := repo.DB.QueryRow("SELECT * FROM library_sys.Books WHERE id = ?", id).Scan(&b.ID, &b.Title, &b.Author, &b.Genre)
    if err != nil {
        return nil, err
    }

    return &b, nil

}


func (repo *BookRepository) RegisterBook(book Book) {
    sqlStatement := "INSERT INTO library_sys.Books (id, title, author, genre) VALUES (?, ?, ?, ?)"
    _, err := repo.DB.Exec(sqlStatement, book.ID, book.Title, book.Author, book.Genre)
    if err != nil {
        log.Fatal(err)
        return
    }
}