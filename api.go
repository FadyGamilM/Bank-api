// => this file constructs my api server

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// our server type
type apiServer struct {
	port string
}

// constructor
func NewApiServer(listenAddress string) *apiServer {
	return &apiServer{
		port: listenAddress,
	}
}

// an api error type to customize my errors into one type
type ApiError struct {
	Error string
}

// my custom snyntax of the http.handlerFunc
type apiHttpHandlerFunc func(http.ResponseWriter, *http.Request) error

// some handlers to handle http requests
// actually these handlers are not perfectly tied to the http.handler syntax .. because i don't need to handle the error inside the handler, i need to return it and handle it somewhere else
func (server *apiServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	// select the appropriate handler according to the http method
	switch r.Method {
	case http.MethodGet:
		return server.handleGetAccount(w, r)
	case http.MethodPost:
		return server.handleCreateAccount(w, r)
	case http.MethodDelete:
		return server.handleDeleteAccount(w, r)
	default:
		return fmt.Errorf("Method Is Not Allowed %s", r.Method)
	}
}

func (server *apiServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	return SendJson(w, http.StatusOK, AccountFactory("Fady", "Gamil"))
}

func (server *apiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (server *apiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (server *apiServer) handleTransferAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// wrapper above my handlers to convert them to the http.handlerFunc syntax
func makeHttpHandlerFunc(apiFunc apiHttpHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO => handle the error from my custom handlers here
		if err := apiFunc(w, r); err != nil {
			// the normal standard http.HandlerFunc doesn't return anything .. so we will send the json response
			// TODO => send the appropriate status code
			SendJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

// a method to return a json response into the response writer variable
func SendJson(w http.ResponseWriter, statusCode int, value any) error {
	// set the header format
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// set the status code
	w.WriteHeader(statusCode)

	// send the json response
	return json.NewEncoder(w).Encode(value)
}

// method to construct the server router, register the handlers
func (server *apiServer) Run() {
	// instantiate the server router
	router := mux.NewRouter()

	// register the handler
	router.HandleFunc("/account", makeHttpHandlerFunc(server.handleAccount))

	//  logging
	log.Println("server is up and running on port ", server.port)

	// listen to a port
	http.ListenAndServe(server.port, router)
}
