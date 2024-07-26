package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
    Payload string `json:"payload"`
    jwt.StandardClaims
}


func GenerateJWT(payload string, secretKey string) (string, error) {
    // Set expiration date
    expirationTime := time.Now().Add(24 * time.Hour)

    // create custom claims using payload and expiration
    claims := CustomClaims {
        Payload: payload,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },

    }

    initToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    token, err := initToken.SignedString([]byte(secretKey))
    if err != nil {
        return "", err
    }

    return token, nil
}


func IsValid(token string, secretKey string) bool {
    _, err := jwt.Parse(token, func (token *jwt.Token) (interface{}, error) {
        return []byte(secretKey), nil
    })
   if err != nil {
		return false
	}
	return true
}
