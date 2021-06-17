package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

/*
 Function to add a list of products to the DB
 @param (http.ResponseWriter, *http.Request)
 @return (void)
*/
func postProducts(response http.ResponseWriter, request *http.Request) {

	// Check if date is given or not
	date := checkDate(request)

	// Create slice of products
	currentTransactionProducts := []Product{}

	// Decode request body into products slice
	jsonErr := json.NewDecoder(request.Body).Decode(&currentTransactionProducts)
	// Check if error occurred
	if jsonErr != nil {
		message := fmt.Sprintf("PostProducts:Error ... %v", jsonErr)
		response.Write([]byte(message))
		log.Fatal(message)
	}

	//------------- Make request to DB -------------
	// Create schema for the fields
	schema := `
		type:  string @index(exact) . 
		name:  string @index(term) .
		price: string @index(term) .
	`
	// Convert slice into JSON format
	productsJson, err := json.Marshal(currentTransactionProducts)
	if err != nil {
		log.Fatal(err)
	}
	// Send mutation to DB
	dbResponse, err := dbRequest(schema, productsJson)
	if err != nil {
		log.Fatalf("PostProducts : Error while making request to DB => %v", err)
	}

	//------------- Process DB response -------------
	// Save the UIDs assigned by DGraph to the products
	products = dbResponse.Uids
	// Fill suggestions map with empty map
	for _, productUid := range products {
		suggestions[productUid] = make(map[string]int)
	}
	// Save products objects in 'productsObj' map
	for _, productObject := range currentTransactionProducts {
		// Save product in map with the Uid as the key
		productId := productObject.Uid[2:]      // Id given in request.body => _:abcdef
		productUid := products[productId]       // Uid given by DGraph
		productsObj[productUid] = productObject // Assign product object
	}

	//------------- Send succsessfull response -------------
	message := fmt.Sprintf("PostProducts: Added [%v] products, date [%v]", len(products), date)
	response.Write([]byte(message))
}
