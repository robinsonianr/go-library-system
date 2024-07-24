package db

type BookRepositoryInterface interface {
    RegisterBook()
    FindAllBooks() ([]Book, error)
    FindBookByID(id string) (*Book, error)
}
