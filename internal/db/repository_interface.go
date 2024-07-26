package db

type BookRepositoryInterface interface {
    RegisterBook(title, author, genre string, pages int) error
    FindAllBooks() ([]Book, error)
    FindBookByID(id string) (*Book, error)
    FindBookByIdentifier(identifier string) (*Book, error)
}


type UserRepositoryInterface interface {
    RegisterUser(username, email, password string) error
    LoginUser(username, password string) (string, error)
}
