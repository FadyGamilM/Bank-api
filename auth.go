package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

// a method to create a token for the given account
func CreateJwtToken(account *Account) (string, error) {
	// create our claims
	claims := &jwt.MapClaims{
		"expiresAt": 15000,
		// the issuer of the token is marked by the number of the account ..
		"accountNumber": account.Number,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwt_Secret := os.Getenv("JWT_SECRET")

	return token.SignedString(jwt_Secret) // this method returns string, error
}

// a method to validate the token
func ValidateJwtToken(token string) (*jwt.Token, error) {
	jwt_Secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method : %v", token.Header["alg"])
		}
		return []byte(jwt_Secret), nil
	})
}

// a method that authorize the given http.HandlerFunc
func Authorize(handler http.HandlerFunc) http.HandlerFunc {
	// return a http handler wrapper for the received http handler func
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Calling the jwt authorization middleware")
		// get the token from the request header
		userToken := r.Header.Get("jwt-token")

		_, err := ValidateJwtToken(userToken)
		if err != nil {
			SendJson(w, http.StatusForbidden, ApiError{Error: "Invalid Token"})
		}

		handler(w, r)
	}
}
