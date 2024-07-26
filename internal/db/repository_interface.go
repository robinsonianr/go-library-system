package db

type BookRepositoryInterface interface {
    RegisterBook(title, author, genre string) error
    FindAllBooks() ([]Book, error)
    FindBookByID(id string) (*Book, error)
}


type UserRepositoryInterface interface {
    RegisterUser(username, email, password string) error
    LoginUser(username, password string) (string, error)
}
