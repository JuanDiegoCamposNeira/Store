package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//--------------------------------------------------------
//						Structures
//--------------------------------------------------------
type Product struct {
	Id, Name, Price string
}

type Person struct {
	Id, Name string
	Age      int
}

type Transaction struct {
	Id, BuyerId, Ip, Device string
	Products                []string
}

//--------------------------------------------------------
//						Methods
//--------------------------------------------------------
//----------------------- Utils --------------------------
/*
 Function to check date in a request
 param (httpRequest)
 return (string) date given in the url or the current day by default
*/
func checkDate(request *http.Request) string {
	date := request.URL.Query().Get("date")
	if date == "" {
		currentDay := time.Now().Format("2006-01-02")
		fmt.Println("Empty date")
		date = currentDay
	}
	return date
}

//----------------------- Products -----------------------
// Function to add a list of products to the DB
func postProducts(response http.ResponseWriter, request *http.Request) {

	// Check if date is given or not
	date := checkDate(request)

	// Create slice of products
	products := []Product{}

	// Decode request body into products slice
	err := json.NewDecoder(request.Body).Decode(&products)
	// Check whether an error occurred while parsing data or not
	if err != nil {
		message := fmt.Sprintf("PostProducts:Error ... %v", err)
		response.Write([]byte(message))
		return
	}

	// Send succsessfull message
	message := fmt.Sprintf("PostProducts: Added [%v] products, date [%v]", len(products), date)
	response.Write([]byte(message))
}

//------------------------ Buyers ------------------------
// Function to add a list of buyers to the DB
func postBuyers(response http.ResponseWriter, request *http.Request) {

	// Check for date in request
	date := checkDate(request)

	// Create slice to store buyers
	buyers := []Person{}

	// Decode body request into buyers slice
	err := json.NewDecoder(request.Body).Decode(&buyers)
	// Check whether an error occurred while parsing data or not
	if err != nil {
		message := fmt.Sprintf("PostBuyers: Error ... %v", err)
		response.Write([]byte(message))
		return
	}

	// Send succsessfull message
	message := fmt.Sprintf("PostBuyers: Added [%v] buyers, date [%v]", len(buyers), date)
	response.Write([]byte(message))
}

// Function to get all the buyers that have been buy in the plattform
func getBuyers(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("Hello from get buyers"))
}

// Function to get a buyer given an Id
func getBuyerById(response http.ResponseWriter, request *http.Request) {
	buyerId := chi.URLParam(request, "buyerId")
	fmt.Printf("Buyer Id : %v", buyerId)
	fmt.Println()
	response.Write([]byte("Hello from get buyer by id"))
}

//--------------------- Transactions ---------------------
// Function to add a list of transactions
func postTransactions(response http.ResponseWriter, request *http.Request) {

	// Check if date is given or not
	date := checkDate(request)

	// Create slice to store transactions
	transactions := []Transaction{}

	// Decode request body into products slice
	err := json.NewDecoder(request.Body).Decode(&transactions)
	// Check whether an error occurred while parsing data or not
	if err != nil {
		message := fmt.Sprintf("PostTransactions:Error ... %v", err)
		response.Write([]byte(message))
		return
	}

	fmt.Printf("Transaction %v \n", transactions[0])

	// Send succsessfull message
	message := fmt.Sprintf("PostTransactions: Added [%v] transactions, date [%v]", len(transactions), date)
	response.Write([]byte(message))

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
	router.Post("/products", postProducts) // Add products with or without a date

	//--------- BUYER endpoints ---------
	router.Post("/buyers", postBuyers)           // Add buyers with or without a date
	router.Get("/buyers", getBuyers)             // Get buyers
	router.Get("/buyer/{buyerId}", getBuyerById) // Get buyer given an id

	//--------- TRANSACTION endpoints ---------
	router.Post("/transactions", postTransactions) // Add transactions

	//--------- Listenning server ---------
	log.Fatal(http.ListenAndServe(":3000", router))
}
