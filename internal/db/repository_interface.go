package db

type BookRepositoryInterface interface {
    RegisterBook(id, title, author, genre string) error
    FindAllBooks() ([]Book, error)
    FindBookByID(id string) (*Book, error)
}
