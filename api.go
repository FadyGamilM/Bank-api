// => this file constructs my api server

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// our server type
type apiServer struct {
	port    string
	storage Storage
}

// constructor
func NewApiServer(listenAddress string, storage Storage) *apiServer {
	return &apiServer{
		port:    listenAddress,
		storage: storage,
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
// TODO => to handle [GET, POST] requests to `/api/accounts`
func (server *apiServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	// select the appropriate handler according to the http method
	switch r.Method {
	case http.MethodGet:
		return server.HandleGetAccounts(w, r)
	case http.MethodPost:
		return server.handleCreateAccount(w, r)
	default:
		return fmt.Errorf("Method Is Not Allowed %s", r.Method)
	}
}

// TODO => to handle [GET, PUT, DELETE] requests to `/api/accounts/{id}`
func (server *apiServer) handleAccountById(w http.ResponseWriter, r *http.Request) error {
	// select the appropriate method
	switch r.Method {
	case http.MethodGet:
		return server.handleAccountById(w, r)
	case http.MethodDelete:
		return server.handleDeleteAccount(w, r)
	default:
		return fmt.Errorf("Method is not allowed")
	}
}

func (server *apiServer) HandleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := server.storage.GetAll()
	if err != nil {
		log.Println("Error while fetching all accounts from db => ", err)
		return err
	}
	return SendJson(w, http.StatusOK, accounts)
}

func (server *apiServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Println("Error while formatting the passed account id from the url => ", err)
		return err
	}
	account, err := server.storage.GetById(id)
	if err != nil {
		log.Println("Error while fetching account by id => ", err)
		return err
	}

	return SendJson(w, http.StatusOK, account)
}

func (server *apiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	// utilize the dto request
	createAccDto := CreateAccountDTO{}

	// decode the given object to the account type we have defined
	if err := json.NewDecoder(r.Body).Decode(&createAccDto); err != nil {
		log.Println("Error while decoding the new account in POST request => ", err)
		return err
	}

	// use the storage apis we have defined in the storage interface to create the entity and add it to the database
	newAccount := AccountFactory(createAccDto.FirstName, createAccDto.LastName)
	// handle errors if there are any
	if err := server.storage.Create(newAccount); err != nil {
		log.Println("Error while creating the new account in POST request => ", err)
		return err
	}

	// return the response
	return SendJson(w, http.StatusCreated, newAccount)
}

func (server *apiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	accountID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return fmt.Errorf("Error while converting the passed id parameter")
	}
	err = server.storage.DeleteById(accountID)
	if err != nil {
		return fmt.Errorf("Error while deleting an account with id = %d", accountID)
	}
	return SendJson(w, http.StatusAccepted, map[string]int{"deleted": accountID})
}

func (server *apiServer) handleTransferAccount(w http.ResponseWriter, r *http.Request) error {
	// create instance of the transfer request dto we have defined in types file
	accountTransfer := new(TransferToAccountRequest)

	// decode the received object
	err := json.NewDecoder(r.Body).Decode(&accountTransfer)
	if err != nil {
		fmt.Errorf("Error while decoding the request body")
		return err
	}

	// send the response back
	return SendJson(w, http.StatusAccepted, accountTransfer)
}

// ! wrapper above my handlers to convert them to the http.handlerFunc syntax
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

// ! a method to return a json response into the response writer variable
func SendJson(w http.ResponseWriter, statusCode int, value any) error {
	// set the header format
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// set the status code
	w.WriteHeader(statusCode)

	// send the json response
	return json.NewEncoder(w).Encode(value)
}

// !method to construct the server router, register the handlers
func (server *apiServer) Run() {
	// instantiate the server router
	router := mux.NewRouter()

	// register the handler/s
	router.HandleFunc("/api/accounts", makeHttpHandlerFunc(server.handleAccount))
	router.HandleFunc("/api/accounts/{id}", makeHttpHandlerFunc(server.handleAccountById))
	router.HandleFunc("/api/accounts/transfer", makeHttpHandlerFunc(server.handleTransferAccount))

	//  logging
	log.Println("server is up and running on port ", server.port)

	// listen to a port
	http.ListenAndServe(server.port, router)
}
