package db

import (
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("Book not found")

type BookRepository struct {
    DB *sql.DB
}

type Book struct {
    ID      string
    Title   string
    Author  string
    Genre   string
    Pages   int
    Stock   int
}

func (repo *BookRepository) FindAllBooks() ([]Book, error){
    rows, err := repo.DB.Query("SELECT * FROM library_sys.Books")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var books []Book
    for rows.Next() {
        var b Book
        if err :=  rows.Scan(&b.ID, &b.Title, &b.Author, &b.Genre, &b.Pages, &b.Stock); err != nil {
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
    
    sql := "SELECT id, title, author, genre FROM library_sys.Books WHERE id = $1"
    err := repo.DB.QueryRow(sql, id).Scan(&b.ID, &b.Title, &b.Author, &b.Genre, &b.Pages, &b.Stock)
    if err != nil {
        return nil, err
    }

    return &b, nil
}

func (repo *BookRepository) FindBookByIdentifier(indentifier string) (*Book, error) {
    var book Book

    sql := `SELECT * FROM library_sys.Books WHERE LOWER(title) iLIKE '%' || $1 || '%' OR LOWER(author) iLIKE '%' || $2 || '%'`
    err := repo.DB.QueryRow(sql, indentifier, indentifier).Scan(&book.ID, &book.Title, &book.Author, &book.Genre, &book.Pages, &book.Stock)
    if err != nil {
        return nil, err
    }

return &book, nil
}

func (repo *BookRepository) RegisterBook(title, author, genre string, pages int) error {
    id := uuid.New().String()
    var book Book
    var stock int

    checkStock := `SELECT stock FROM library_sys.Books WHERE LOWER(title) = LOWER($1) OR LOWER(author) = LOWER($2)`
    err := repo.DB.QueryRow(checkStock, title, author).Scan(&book.Stock)
    if err != nil {
        log.Println("Book not found creating new entry.")
        stock = 0
    } else {
        stock = book.Stock + 1
    }
    
    if stock < 1 {
        stock += 1
        sqlStatement := "INSERT INTO library_sys.Books (id, title, author, genre, pages, stock) VALUES ($1, $2, $3, $4, $5, $6)"
        _, err = repo.DB.Exec(sqlStatement, id, title, author, genre, pages, stock)
        if err != nil {
            log.Fatal(err)
            return err
        }
    } else {
        updateSql := "UPDATE library_sys.Books SET stock = $1 WHERE LOWER(title) = LOWER($2) OR LOWER(author) = LOWER($3)"
        _, err := repo.DB.Exec(updateSql, stock, title, author)
        if err != nil {
            log.Fatal(err)
            return err
        }
    }


    return nil
}
