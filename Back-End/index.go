package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//--------------------------------------------------------
//						Global variables
//--------------------------------------------------------
var _ map[string]map[string]int                   // [day] : { [Product]: quantity, [Product]: quantity, ... }
var suggestions = make(map[string]map[string]int) // [Product]: { [Product]: quantity, [Product]: quantity, ... }
var productsObj = make(map[string]Product)        // Store products in a map to access them in O(1)

// Store assigned Uids by Dgraph
var buyers = make(map[string]string)
var products = make(map[string]string)
var transactions = make(map[string]string)

//--------------------------------------------------------
//						Main
//--------------------------------------------------------
func main() {

	//-------------- Chi router --------------
	// Create chi router
	router := chi.NewRouter()
	// Chi middlewares
	router.Use(middleware.Logger)

	//--------------   PRODUCT endpoints   --------------
	router.Post("/products", postProducts) // Add products with or without a date

	//--------------   BUYER endpoints   --------------
	router.Post("/buyers", postBuyers)           // Add buyers with or without a date
	router.Get("/buyers", getBuyers)             // Get all buyers
	router.Get("/buyer/{buyerId}", getBuyerById) // Get buyer given an id

	//--------------   TRANSACTION endpoints   --------------
	router.Post("/transactions", postTransactions) // Add transactions

	//--------------   Listenning server   --------------
	fmt.Println("Server listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
