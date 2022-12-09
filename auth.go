package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("my_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func CheckUserSingIn(w http.ResponseWriter, r *http.Request, email string, password string) bool {

	userData := getUserByEmail(email)

	if userData.Email == "" {
		return false
	}

	if bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(password)) != nil {
		return false
	}

	expirationTime := time.Now().Add(3 * time.Minute)

	claims := &Claims{
		Username: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return false
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	http.SetCookie(w, &http.Cookie{
		Name:    "userID",
		Value:   fmt.Sprint(userData.Id),
		Expires: expirationTime,
	})

	http.SetCookie(w, &http.Cookie{
		Name:    "isAdmin",
		Value:   fmt.Sprint(userData.IsAdmin),
		Expires: expirationTime,
	})
	return true
}

func CheckUserAuth(w http.ResponseWriter, r *http.Request) bool {

	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {

			return false
		}

		return false
	}

	tknStr := c.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return false
		}
		return false
	}
	if !tkn.Valid {
		return false
	}

	return true
}

func UserLogOut(w http.ResponseWriter, r *http.Request) {

	expirationTime := time.Now().Add(0 * time.Minute)

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: expirationTime,
	})

	http.SetCookie(w, &http.Cookie{
		Name:    "userID",
		Value:   "",
		Expires: expirationTime,
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "isAdmin",
		Value:   "",
		Expires: expirationTime,
	})
}
