package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)


func NewDb() (*sql.DB, error){
    connStr := "user=postgres password=root1234 dbname=robinsonir sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
        return nil, err
    }

    createTablesSQL := `CREATE TABLE IF NOT EXISTS "library_sys".Books (
        ID TEXT PRIMARY KEY,
        Title TEXT NOT NULL,
        Author TEXT NOT NULL,
        Genre TEXT NOT NULL
    );
    
    CREATE TABLE IF NOT EXISTS "library_sys".User (
        ID TEXT PRIMARY KEY,
        Username TEXT NOT NULL,
        Email TEXT NOT NULL,
        Password TEXT NOT NULL
    );
    `



    _, err = db.Exec(createTablesSQL)
    if err != nil {
        panic(err)
    }

    if err := db.Ping(); err != nil {
        return nil, err
    }

    return db, nil
}
