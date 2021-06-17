package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
)

//--------------------------------------------------------
//						Structures
//--------------------------------------------------------
type Day struct {
	Type   string   `json:"type,omitempty"`
	Buyers []Person `json:"buyers,omitempty"`
}

type Person struct {
	Type string `json:"type,omitempty"`
	Uid  string `json:"uid,omitempty"`
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`
}

type Product struct {
	Type  string `json:"type,omitempty"`
	Uid   string `json:"uid,omitempty"`
	Name  string `json:"name,omitempty"`
	Price string `json:"price,omitempty"`
}

type Transaction struct {
	Type     string    `json:"type,omitempty"`
	Uid      string    `json:"uid,omitempty"`
	Buyer    Person    `json:"buyer,omitempty"`
	Ip       string    `json:"ip,omitempty"`
	Device   string    `json:"device,omitempty"`
	Products []Product `json:"products,omitempty"`
}

type Suggestion struct {
	Product  Product
	Quantity int
}

//--------------------------------------------------------
//						DGraph
//--------------------------------------------------------
// // var err error                                                    // To catch errors
// var ctx = context.Background()                                 // Context
// var conn, _ = grpc.Dial("127.0.0.1:9080", grpc.WithInsecure()) // GRPC client connection
// var dc = api.NewDgraphClient(conn)                             // DGrap client
// var dg = dgo.NewDgraphClient(dc)                               // Dgraph go
// var mu = &api.Mutation{CommitNow: true}                        // To make a mutation to the DB

// ctx = context.Background()
// conn, err = grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
// if err != nil {
// 	log.Fatal("While trying to dial gRPC")
// }
// defer conn.Close()

// dc = api.NewDgraphClient(conn)
// dg = dgo.NewDgraphClient(dc)

// mu = &api.Mutation{CommitNow: true}

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
//						Methods
//--------------------------------------------------------
//----------------------- Utils --------------------------
/*
 Function to check date in a request
 @param (httpRequest)
 @return (string) date given in the url or the current day by default
*/
func checkDate(request *http.Request) string {
	date := request.URL.Query().Get("date")
	if date == "" {
		currentDay := time.Now().Format("2006-01-02")
		date = currentDay
	}
	return date
}

/*
 Function to make a request to the DB
 @param (string : Represents the schema to be installed in the DB,
		 []byte : Represents the data parsed into JSON format)
 @return (api.Response : Represents DB response,
	 	  error : If an error occurs)
*/
func dbRequest(schema string, jsonData []byte) (*api.Response, error) {

	//------------- Get context -------------
	ctx := context.Background()
	//------------- Create connection to the DB -------------
	conn, connectionErr := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	// Check if an error occurred while connecting to the DB
	if connectionErr != nil {
		return nil, connectionErr
	}
	// Defer close of the connection until is no longer used
	defer conn.Close()

	//------------- Create DGraph client -------------
	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)

	//------------- Schema  -------------
	// Install a schema into dgraph
	schemaError := dg.Alter(context.Background(), &api.Operation{
		Schema: schema,
	})
	// Check if an error occurred while setting schema into DB
	if schemaError != nil {
		return nil, schemaError
	}

	//------------- Mutation -------------
	// Define mutation
	mu := &api.Mutation{CommitNow: true}
	// Convert parameter slice into JSON format
	mu.SetJson = jsonData
	// Make mutation to the DB
	response, err := dg.NewTxn().Mutate(ctx, mu)
	// Check if error occurred while making mutation to the DB
	if err != nil {
		log.Fatal(err)
	}

	//------------- Successfull response -------------
	return response, nil
}

//----------------------- Products -----------------------
// Function to add a list of products to the DB
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

//------------------------ Buyers ------------------------
// Function to add a list of buyers to the DB
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
	message := fmt.Sprintf("PostBuyers: Added [%v] buyers, date [%v]", len(buyersArr), date)
	response.Write([]byte(message))
}

// Function to get all the buyers that have been buy in the plattform
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

	response.Write(resp.GetJson())
}

// Function to get a buyer given an Id
func getBuyerById(response http.ResponseWriter, request *http.Request) {
	buyerId := chi.URLParam(request, "buyerId")

	ctx := context.Background()
	conn, err := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("While trying to dial gRPC")
	}
	defer conn.Close()

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)

	// Query to get buyer by id
	const buyersQuery = `query history($id: string) { 
							history(func: has(buyer)) @filter(uid_in(buyer, $id)) {
								device
								ip
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

	// Query DB
	variables := map[string]string{"$id": buyers[buyerId]}
	// variables := map[string]string{"$id": "0x22c3e"}
	resp, err := dg.NewTxn().QueryWithVars(ctx, buyersQuery, variables)
	if err != nil {
		log.Fatal(err)
	}

	//
	type Response struct {
		SameIp              []Transaction `json:"sameIp"`
		TransactionsHistory []Transaction `json:"history"`
		Suggestions         [][]Product
	}
	var res Response
	// Decode response to evaluate IP addresses
	err = json.Unmarshal(resp.Json, &res)
	if err != nil {
		log.Fatalf("GetBuyerById:Error while unmarshal JSON => %v", err)
	}
	// Get all Ip addresses
	var ipAddresses string
	for _, transaction := range res.TransactionsHistory {
		ipAddresses += transaction.Ip + " "
	}
	// Debug
	// fmt.Printf("Ip Addresses => %v\n", ipAddresses)

	// Query to get buyer by id
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

	// Query DB
	variables = map[string]string{"$addresses": ipAddresses}
	resp, err = dg.NewTxn().QueryWithVars(ctx, sameAddresses, variables)
	if err != nil {
		log.Fatal(err)
	}
	// Debug
	// fmt.Printf("Transactions with the same address => %v", resp.Json)

	// Parse SameIP response into struct
	err = json.Unmarshal(resp.Json, &res)
	if err != nil {
		log.Fatalf("GetBuyerById:Error while unmarshal Same address to JSON => %v", err)
	}

	// Debug
	fmt.Printf("Products => %v\n", productsObj)
	// Add suggestions to the reponse
	responseSuggestions := [][]Product{}
	for i, transaction := range res.TransactionsHistory {
		// Traverse the products in the transactions
		fmt.Printf("Transaction #%v : \n", i)
		for _, product := range transaction.Products {
			productSuggestions := []Product{}
			// Get the first two items in the map
			count := 0
			fmt.Println("")
			fmt.Printf("Suggestions for product %v => %v\n", product.Uid, suggestions[product.Uid])
			for key := range suggestions[product.Uid] {
				if count == 2 {
					break
				}
				fmt.Printf("Product #%v => %v \n", count, key)
				productSuggestions = append(productSuggestions, productsObj[key])
				count++
			}
			responseSuggestions = append(responseSuggestions, productSuggestions)
		}
		fmt.Println("------------------------")
	}
	// Debug
	fmt.Printf("\nSuggestions => %v \n", responseSuggestions)
	// Add slice with suggestions to response
	res.Suggestions = responseSuggestions

	responseJson, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("GetBuyerById: Error while marshal to JSON => %v", err)
	}

	response.Write(responseJson)
}

//--------------------- Transactions ---------------------
// Function to add a list of transactions
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
	router.Get("/buyers", getBuyers)             // Get buyers
	router.Get("/buyer/{buyerId}", getBuyerById) // Get buyer given an id

	//--------------   TRANSACTION endpoints   --------------
	router.Post("/transactions", postTransactions) // Add transactions

	//--------------   Listenning server   --------------
	fmt.Println("Server listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
