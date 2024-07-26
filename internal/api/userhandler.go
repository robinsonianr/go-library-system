package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ianrr/library/internal/db"
	"github.com/ianrr/library/internal/db/user"
)



type APIUser struct {
    ID          string `json:"id"`
    Username    string `json:"username"`
    Email       string `json:"email"`
    Password    string `json:"password"`
}


type UserHandler struct {
    Repo    db.UserRepositoryInterface
}


func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
    var user APIUser

    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if user.Username == "" || user.Email == "" || user.Password == "" {
        http.Error(w, "Missing required field(s)", http.StatusBadRequest)
        return
    }

    err = h.Repo.RegisterUser(user.Username, user.Email, user.Password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

     log.Println("Successfuly regisered user")
     w.WriteHeader(http.StatusCreated)
}


func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
    err := r.ParseMultipartForm(32 << 20)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    identifier := r.FormValue("identifier")
    password := r.FormValue("password")
    
    if identifier == "" || password == "" {
        http.Error(w, "Missng required field(s)", http.StatusBadRequest)
        return
    }

    token, err := h.Repo.LoginUser(identifier, password)
    if err != nil {
        if err == user.UserNotFound {
            http.Error(w, "User not found", http.StatusNotFound)
            return
        }

        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    expiration := time.Now().Add(24 * time.Hour)
    cookie := http.Cookie {
        Name: "token",
        Value: token,
        Expires: expiration,
        HttpOnly: true,
        Path: "/",
    }
    http.SetCookie(w, &cookie)

    w.WriteHeader(http.StatusCreated)
}
