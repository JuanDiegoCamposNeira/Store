package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

/*
 Function to add a list of transactions
 @param (http.ResponseWriter, *http.Request)
 @return (void)
*/
func postTransactions(response http.ResponseWriter, request *http.Request) {

	// Check if date is given or not
	date := checkDate(request)

	//------------- Process http request -------------
	// Aux struct to decode request's body
	type TransactionAux struct {
		Type, Uid, BuyerId, Ip, Device string
		Products                       []string
	}
	// Create slice to store transactions
	transactionsEncoded := []TransactionAux{} // Transactions without the reference to the buyer or the products
	// Decode request's body into transactions slice
	err := json.NewDecoder(request.Body).Decode(&transactionsEncoded)
	if err != nil {
		message := fmt.Sprintf("PostTransactions:Error ... %v", err)
		response.Write([]byte(message))
		log.Fatal(message)
	}

	//------------- Process transactions -------------
	transactionsDecoded := []Transaction{} // Slice to store transactions with the reference to the UIDs in the DB
	// Traverse encoded transactions
	for _, transaction := range transactionsEncoded {

		// Slices to store products
		currentTransactionDecodedProducts := []Product{} // Slice to store the UIDs of the products given by Dgraph
		currentTransactionProducts := transaction.Products

		//------------- Create require DGraph format -------------
		// Create DGraph required format to link an existing product node to the transaction
		// Format example :
		//		{
		// 			...
		// 			Product : [ { Uid: <uid_given_by_DGraph> } [, { Uid :  <uid_given_by_DGraph> }, ...] ]
		// 			...
		// 		}
		//
		for _, productId := range currentTransactionProducts {
			currentTransactionDecodedProducts = append(currentTransactionDecodedProducts, Product{Uid: products[productId]})
		}

		//------------- Fill suggestions map -------------
		for i, i_productId := range currentTransactionProducts {
			for j, j_productId := range currentTransactionProducts {
				// Don't include the product being evaluated
				if i == j {
					continue
				}
				// Add one to the suggestions map of the product
				suggestions[products[i_productId]][products[j_productId]] += 1
			}
		}

		//------------- Decode transaction -------------
		// Create transaction with references (UIDs) to Products and Buyer in DB
		decodedTransaction := Transaction{
			Type:     "Transaction",
			Uid:      "_:" + transaction.Uid,
			Buyer:    Person{Uid: buyers[transaction.BuyerId]},
			Ip:       transaction.Ip,
			Device:   transaction.Device,
			Products: currentTransactionDecodedProducts,
		}
		// Add decoded transaction
		transactionsDecoded = append(transactionsDecoded, decodedTransaction)
	}

	//------------- Make request to DB -------------
	// Create schema for the fields
	schema := `
		type:  		string @index(exact) . 
		buyer:  	uid .
		ip:  		string @index(term) .
		device: 	string . 
		products: 	[uid]  .
	`
	// Convert slice into JSON format
	transactionsJson, err := json.Marshal(transactionsDecoded)
	if err != nil {
		log.Fatal(err)
	}
	// Send mutation to DB
	dbResponse, err := dbRequest(schema, transactionsJson)
	if err != nil {
		log.Fatalf("PostProducts : Error while making request to DB => %v", err)
	}

	//------------- Process DB response -------------
	// Save the UIDs assigned by DGraph to the transactions
	transactions = dbResponse.Uids

	//------------- Send succsessfull response -------------
	message := fmt.Sprintf("PostTransactions: Added [%v] transactions, date [%v]", len(transactionsDecoded), date)
	response.Write([]byte(message))
}
