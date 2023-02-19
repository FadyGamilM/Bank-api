package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

// a method to create a token for the given account
func CreateJwtToken(account *Account) (string, error) {
	// create our claims
	claims := &jwt.MapClaims{
		"expiresAt": 15000,
		// the issuer of the token is marked by the number of the account ..
		"accountNumber": int64(account.Number),
	}

	fmt.Println(claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwt_Secret := os.Getenv("JWT_SECRET")

	return token.SignedString([]byte(jwt_Secret)) // this method returns string, error
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
func Authorize(handler http.HandlerFunc, storage Storage) http.HandlerFunc {
	// return a http handler wrapper for the received http handler func
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Calling the jwt authorization middleware")
		// get the token from the request header
		userToken := r.Header.Get("jwt-token")

		validatedToken, err := ValidateJwtToken(userToken)
		if err != nil {
			SendJson(w, http.StatusForbidden, ApiError{Error: "Invalid Token"})
			return
		}

		if !validatedToken.Valid {
			SendJson(w, http.StatusForbidden, ApiError{Error: "Invalid Token"})
			return
		}

		claims := validatedToken.Claims.(jwt.MapClaims)

		log.Println("AFTER VALIDATIONG THE TOKEN => ", claims["accountNumber"])
		log.Println("AFTER VALIDATIONG THE TOKEN => ", claims)
		log.Println(reflect.TypeOf(claims["accountNumber"]))

		// fetch the account based on the given id to compare the embedded data on it with the embedded data within the jwt token to know if this is the right user
		userID, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			SendJson(w, http.StatusInternalServerError, ApiError{Error: "something wrong happened while processing your request"})
			return
		}
		accountFromDB, err := storage.GetById(userID)
		if err != nil {
			SendJson(w, http.StatusForbidden, ApiError{Error: "perfmission denied"})
			return
		}
		fmt.Println("the number of the fetched account => ", accountFromDB.Number)

		if accountFromDB.Number != int64(claims["accountNumber"].(float64)) {
			SendJson(w, http.StatusForbidden, ApiError{Error: "perfmission denied"})
			return
		}
		fmt.Println("user id => ", accountFromDB.Number)
		fmt.Println("token claims => ", int64(claims["accountNumber"].(float64)))
		handler(w, r)
	}
}
