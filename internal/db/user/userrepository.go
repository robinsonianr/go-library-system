package user

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/ianrr/library/internal/auth"
	"github.com/ianrr/library/internal/db"
	"github.com/ianrr/library/internal/utils"
	"golang.org/x/crypto/bcrypt"
)


type UserRepository struct {
   DB  *sql.DB 
}


type User struct {
    ID          string
    Username    string
    Email       string
    Password    string
    Books       []db.Book
}

var UserNotFound = errors.New("User not found")

func hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    if err != nil {
        return "", err
    }

    return string(bytes), nil
}

func (repo *UserRepository) RegisterUser(username, email, password string) error {
    userId := uuid.New().String()
    hashedPassword, err := hashPassword(password)
    if err != nil {
        return err
    }

    sql := "INSERT INTO library_sys.User (id, username, email, password) VALUES ($1, $2, $3, $4)"
    _, err = repo.DB.Exec(sql, userId, username, email, hashedPassword)
    if err != nil {
        return  err 
    }

    return nil
}


func (repo *UserRepository) LoginUser(identifier, password string) (string, error) {
    var user User

    sql := "SELECT * FROM library_sys.User WHERE username = $1 OR email = $2"

    err := repo.DB.QueryRow(sql, identifier, identifier).Scan(&user.ID, &user.Email, &user.Username, &user.Password)
    if err != nil {
        return "",  UserNotFound 
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return "", UserNotFound
    }

    secret, err := utils.ReadFile("../../privateKey.txt")
    if err != nil {
        return "", err
    }

    jwt, err := auth.GenerateJWT(user.ID, secret)
    if err != nil {
        return "", err
    }

    return jwt, nil
}
