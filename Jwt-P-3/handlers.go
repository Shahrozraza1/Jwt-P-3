package main

import (
	"encoding/json"
	"fmt"
	"go/token"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)
var jwtKey = []byte("secret_Key")
var users = maps [string] string{
"user1" :"password1",
"user2" : "password2"
}
type Credentials struct {
	Username string `json: "username"`
	Password string `json: "password"`
}
type Claims struct {
	Username string `json: "username"`
	jwt.StandardClaims
}
func Login(w http.ResponseWriter, r *http.Request){
	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expectedPassword, ok := users [credentials.Username]
	if !ok || expectedPassword != credentials.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return

	}
	expirationTime := time.Now().Add(time.Minute*5)
	Claims := &Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString , err := token.SignedString(jwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w,
	&http.Cookie{
		Name: "token",
		Value: tokenString,
		Expires : expirationTime,
	} )

}
func Home (w http.ResponseWriter, r *http.Request) {
	cookies, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie{
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tokenStr := cookie.Value
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
	func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusBadRequest)
			request
		}
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Write([]byte(fmt.Sprintf("Hello, %s", claims.Username)))
	}
	func Refresh(w http.ResponseWriter, r *http.Request){
		cookies, err := r.Cookies("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		tokenStr := cookie.Value
		claims := &Claims{}
		tkn , err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface {}, error ) {
			return jwtKey, nil
		})
		if err != nil{
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			request
		}
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		expirationTime := time.Now().Add(time.Minute*5)
		claims.ExpiresAt = expirationTime.Unix()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		} 
		http.SetCookie(w,
		&http.Cookies{
			Name : "refresh_token",
			Value : tokenString,
			Expires : expirationTime,
		})
	}
}

