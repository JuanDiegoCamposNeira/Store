package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//--------------------------------------------------------
//						Structures
//--------------------------------------------------------
type Product struct {
	Id, Name string
	Price    int
}

type Person struct {
	Id, Name, Age string
}

type Transaction struct {
	Id, BuyerId, Ip, Device string
	ProductIds              []string
}

//--------------------------------------------------------
//						Methods
//--------------------------------------------------------
//----------------------- Products -----------------------
// Function to add a list of products to the DB
func postProducts(response http.ResponseWriter, request *http.Request) {
	// Create slice of products
	products := []Product{}

	// Decode request body into products slice
	err := json.NewDecoder(request.Body).Decode(&products)
	if err != nil {
		message := fmt.Sprintf("Error ... %v", err)
		response.Write([]byte(message))
		return
	}
	fmt.Println(products)
}

//------------------------ Buyers ------------------------
// Function to add a list of buyers to the DB
func postBuyers(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("Hello from post buyers"))
}

// Function to get all the buyers that have been buy in the plattform
func getBuyers(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("Hello from get buyers"))
}

// Function to get a buyer given an Id
func getBuyerById(response http.ResponseWriter, request *http.Request) {
	// buyerId := chi.URLParam(router)
	buyerId := "Some ID"
	fmt.Printf("Buyer Id : %v", buyerId)
	response.Write([]byte("Hello from get buyer by id"))
}

//--------------------- Transactions ---------------------
// Function to add a list of transactions
func postTransactions(response http.ResponseWriter, request *http.Request) {
	fmt.Println("addTransactions : Called")
	response.Write([]byte("Hello from transactions"))
}

//--------------------------------------------------------
//						Main
//--------------------------------------------------------
func main() {

	//--------- Chi router ---------
	// Create chi router
	router := chi.NewRouter()
	// Chi middlewares
	router.Use(middleware.Logger)

	//--------- PRODUCT endpoints ---------
	router.Post("/products", postProducts) // Add products

	//--------- BUYER endpoints ---------
	router.Post("/buyers", postBuyers)           // Add buyers
	router.Get("/buyers", getBuyers)             // Get buyers
	router.Get("/buyer/{buyerId}", getBuyerById) // Get buyer given an id

	//--------- TRANSACTION endpoints ---------
	router.Post("/transactions", postTransactions) // Add transactions

	//--------- Listenning server ---------
	log.Fatal(http.ListenAndServe(":3000", router))
}
