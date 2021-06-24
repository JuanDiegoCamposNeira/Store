package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

/*
 Function to add a list of buyers to the DB
 @param (http.ResponseWriter, *http.Request)
 @return (void)
*/
func postBuyers(response http.ResponseWriter, request *http.Request) {

	// Check for date in request
	date := checkDate(request)

	// Create slice to store buyers
	buyersArr := []Person{}

	// Decode body request into buyers slice
	err := json.NewDecoder(request.Body).Decode(&buyersArr)
	// Check if error occurred
	if err != nil {
		message := fmt.Sprintf("PostBuyers: Error ... %v", err)
		response.Write([]byte(message))
		log.Fatal(message)
	}

	//------------- Make request to DB -------------
	// Create schema for the fields
	schema := `
		type:  string @index(exact) . 
		name:  string @index(term) .
		age:   int .
	`
	// Convert slice into JSON format
	buyersJson, err := json.Marshal(buyersArr)
	if err != nil {
		log.Fatal(err)
	}
	// Send mutation to DB
	dbResponse, err := dbRequest(schema, buyersJson)
	if err != nil {
		log.Fatalf("PostBuyers : Error while making request to DB => %v", err)
	}

	//------------- Process DB response -------------
	buyers = dbResponse.Uids // Save the UIDs assigned by DGraph to the buyers

	//------------- Send succsessfull response -------------
	response.Header().Set("Access-Control-Allow-Origin", "*")
	message := fmt.Sprintf("PostBuyers: Added [%v] buyers, date [%v]", len(buyersArr), date)
	response.Write([]byte(message))
}

/*
 Function to get all the buyers that have been buy in the plattform
 @param (http.ResponseWriter, *http.Request)
 @return (void)
*/
func getBuyers(response http.ResponseWriter, request *http.Request) {

	ctx := context.Background()
	conn, err := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("While trying to dial gRPC")
	}
	defer conn.Close()

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)

	// Query to get all the buyers
	const buyersQuery = `
							{
								buyers(func: eq(type, "Person")) {
									name
									id
									age
									uid
								}
							}`

	// Ask for the type of name and age.
	resp, err := dg.NewTxn().Query(ctx, buyersQuery)
	if err != nil {
		log.Fatal(err)
	}

	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Write(resp.GetJson())
}

/*
 Function to get a buyer given an Id
 @param (http.ResponseWriter, *http.Request)
 @return (void)
*/
func getBuyerById(response http.ResponseWriter, request *http.Request) {

	// Get requested buyer Id passed as parameter
	buyerId := chi.URLParam(request, "buyerId")

	// Struct to save reponse data
	type Response struct {
		TransactionsHistory []Transaction `json:"history"`
		SameIp              []Transaction `json:"sameIp"`
		Suggestions         map[string][]Product
	}
	// Instance of the response struct
	var res Response

	//------------- Make request to DB -------------
	ctx := context.Background()
	conn, err := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("While trying to dial gRPC")
	}
	defer conn.Close()

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)

	//------------- Query transactions history to the DB -------------
	// Construct query
	const buyersQuery = `query history($id: string) { 
							history(func: has(buyer)) @filter(uid_in(buyer, $id)) {
								uid
								device
								ip
								Date
								buyer {
									age
									uid
									name
								}
								products {
									uid
									price
									name
								}
							}
						}`
	// Send query to the DB
	variables := map[string]string{"$id": buyerId}
	resp, err := dg.NewTxn().QueryWithVars(ctx, buyersQuery, variables)
	if err != nil {
		log.Fatal(err)
	}
	//----- Process query response -----
	// Debug
	fmt.Println()
	// Decode transactions history
	err = json.Unmarshal(resp.Json, &res)
	if err != nil {
		log.Fatalf("GetBuyerById:Error while unmarshal JSON => %v", err)
	}

	//------------- Query same IP addesses to the DB -------------
	// Get IP addresses from every transaction made
	var ipAddresses string
	for _, transaction := range res.TransactionsHistory {
		ipAddresses += transaction.Ip + " "
	}
	// Create query to get the users using the same address of the given buyer
	const sameAddresses = `	query sameIp($addresses: string) { 
								sameIp(func: anyofterms(ip, $addresses))  {
									device
									ip
									buyer {
										uid
										age
										name
									}
								}
							}`
	// Send query to the DB
	variables = map[string]string{"$addresses": ipAddresses}
	resp, err = dg.NewTxn().QueryWithVars(ctx, sameAddresses, variables)
	if err != nil {
		log.Fatal(err)
	}
	//----- Process query response -----
	// Decode users with the sameIp
	err = json.Unmarshal(resp.Json, &res)
	if err != nil {
		log.Fatalf("GetBuyerById:Error while unmarshal SameAddress to JSON => %v", err)
	}

	//------------- Get suggestions based on the most buyed products -------------
	// Add suggestions to the reponse
	responseSuggestions := map[string][]Product{} // Map to store suggestions
	for _, transaction := range res.TransactionsHistory {
		// Traverse the products in the transactions
		for _, product := range transaction.Products {
			mostOrderedProduct := productsObj[product.Uid]
			secondMostOrderedProduct := productsObj[product.Uid]
			suggestions[product.Uid][product.Uid] = 0
			//
			for key := range suggestions[product.Uid] {
				if suggestions[product.Uid][key] > suggestions[product.Uid][mostOrderedProduct.Uid] {
					// The most Ordered product will be in second place now
					secondMostOrderedProduct = mostOrderedProduct
					mostOrderedProduct = productsObj[key]
				} else if suggestions[product.Uid][key] > suggestions[product.Uid][secondMostOrderedProduct.Uid] {
					secondMostOrderedProduct = productsObj[key]
				}
			}
			productSuggestions := []Product{mostOrderedProduct, secondMostOrderedProduct}
			// Save suggestions for the current product
			responseSuggestions[product.Name] = productSuggestions
		}
	}
	// Add slice with suggestions to response
	res.Suggestions = responseSuggestions

	responseJson, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("GetBuyerById: Error while marshal to JSON => %v", err)
	}

	//------------- Send succsessfull response -------------
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Write(responseJson)
}
