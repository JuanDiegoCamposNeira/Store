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
 param (httpRequest)
 return (string) date given in the url or the current day by default
*/
func checkDate(request *http.Request) string {
	date := request.URL.Query().Get("date")
	if date == "" {
		currentDay := time.Now().Format("2006-01-02")
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
	products_arr := []Product{}

	// Decode request body into products slice
	err := json.NewDecoder(request.Body).Decode(&products_arr)
	// Check whether an error occurred while parsing data or not
	if err != nil {
		message := fmt.Sprintf("PostProducts:Error ... %v", err)
		response.Write([]byte(message))
		return
	}
	//------- Modification
	// Create and save Objects of products
	// products_quantity := map[string]int{}
	// for _, product := range products_arr {
	// 	// Create Product to store in map
	// 	p := Product{
	// 		Type:  product.Type,
	// 		Uid:   product.Uid[2:],
	// 		Name:  product.Name,
	// 		Price: product.Price,
	// 	}
	// 	// Save product in map with the object as a key
	// 	productsObj[p.Uid] = p
	// 	//
	// 	products_quantity[p.Uid] = 0
	// }
	//----- Modification
	// for index := range suggestions {
	// 	suggestions[index] = products_quantity
	// }

	// Debug
	// fmt.Printf("Suggestions Map => %v \n", suggestions)

	// Add products to DB
	//------------------------------
	// Connection to DB
	ctx := context.Background()
	conn, err := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("While trying to dial gRPC")
	}
	defer conn.Close()

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)

	// Schema
	// Install a schema into dgraph. Accounts have a `name` and a `balance`.
	err1 := dg.Alter(context.Background(), &api.Operation{
		Schema: `
			type:  string @index(exact) . 
			name:  string @index(term) .
			price: string @index(term) .
		`,
	})
	if err1 != nil {
		log.Fatal("postProducts:Error setting schema")
	}

	// Mutation
	mu := &api.Mutation{CommitNow: true}

	pb, err := json.Marshal(products_arr)
	if err != nil {
		log.Fatal(err)
	}
	mu.SetJson = pb
	assigned, err := dg.NewTxn().Mutate(ctx, mu)
	if err != nil {
		log.Fatal(err)
	}
	// Debug
	// fmt.Printf("Response => %v \n", assigned)
	//------------------------------

	// Save products assigned UIDs
	products = assigned.Uids
	// Debug
	// fmt.Printf("Saved Products => %v \n", products)
	// Fill suggestions map
	for _, productUid := range products {
		suggestions[productUid] = make(map[string]int)
	}

	//-------- Modification
	// Debug
	// fmt.Printf("Assigned UIDS => %v\n", products)
	// Create and save Objects of products
	for _, product := range products_arr {
		// Save product in map with the object as a key
		productsObj[products[product.Uid[2:]]] = product // product.Uid = _:abcdef
	}
	//--------

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
	buyers_arr := []Person{}

	// Decode body request into buyers slice
	err := json.NewDecoder(request.Body).Decode(&buyers_arr)
	// Check whether an error occurred while parsing data or not
	if err != nil {
		message := fmt.Sprintf("PostBuyers: Error ... %v", err)
		response.Write([]byte(message))
		return
	}

	// person := Person{
	// 	Type: "Person",
	// 	Uid:  "_:1c5a5873",
	// 	Name: "Kilroy",
	// 	Age:  47,
	// }

	// Add buyers to DB
	//------------------------------
	// Connection to DB
	ctx := context.Background()
	conn, err := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("While trying to dial gRPC")
	}
	defer conn.Close()

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)

	// Schema
	// Install a schema into dgraph. Accounts have a `name` and a `balance`.
	err1 := dg.Alter(context.Background(), &api.Operation{
		Schema: `
			type:  string @index(exact) . 
			name:  string @index(term) .
			age:   int .
		`,
	})
	if err1 != nil {
		errorMessage := fmt.Sprintf("postProducts:Error setting schema -> %v", err1)
		log.Fatal(errorMessage)
	}

	// Mutation
	mu := &api.Mutation{CommitNow: true}

	pb, err := json.Marshal(buyers_arr)
	if err != nil {
		log.Fatal(err)
	}
	mu.SetJson = pb
	assigned, err := dg.NewTxn().Mutate(ctx, mu)
	if err != nil {
		log.Fatal(err)
	}
	//------------------------------

	buyers = assigned.Uids
	// Debug
	// fmt.Println("Assigned => ", buyers)

	// Send succsessfull message
	message := fmt.Sprintf("PostBuyers: Added [%v] buyers, date [%v]", len(buyers_arr), date)
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

	// Aux struct to decode body in the request
	type TransactionAux struct {
		Type     string
		Uid      string
		BuyerId  string
		Ip       string
		Device   string
		Products []string
	}
	// Create slice to store transactions
	transactions_encoded := []TransactionAux{} // Transactions without the reference to the buyer or the products

	// Decode request body into transactions slice
	err := json.NewDecoder(request.Body).Decode(&transactions_encoded)
	// Check whether an error occurred while parsing data or not
	if err != nil {
		message := fmt.Sprintf("PostTransactions:Error ... %v", err)
		response.Write([]byte(message))
		return
	}

	//
	transactions_decoded := []Transaction{}

	// Debug
	// fmt.Printf("Products : %v \n", products)
	for _, transaction := range transactions_encoded {
		// Get products UIDs
		products_arr := []Product{}
		for i, productId := range transaction.Products {
			products_arr = append(products_arr, Product{Uid: products[productId]})

			// Fill suggestions map
			for j, j_productId := range transaction.Products {
				if i == j {
					continue
				}

				suggestions[products[productId]][products[j_productId]] += 1
			}
		}

		// Create transaction with links to Products and Buyer
		transaction_temp := Transaction{
			Type:     "Transaction",
			Uid:      "_:" + transaction.Uid,
			Buyer:    Person{Uid: buyers[transaction.BuyerId]},
			Ip:       transaction.Ip,
			Device:   transaction.Device,
			Products: products_arr,
		}
		// Add decoded transaction
		transactions_decoded = append(transactions_decoded, transaction_temp)
	}

	// transaction := Transaction{
	// 	Type:   "Transaction",
	// 	Uid:    "_:000060c3f900",
	// 	Buyer:  Person{Uid: buyers["a"]},
	// 	Ip:     "12.3.2.1",
	// 	Device: "mac",
	// }

	// Add transactions to DB
	//------------------------------
	// Connection to DB
	ctx := context.Background()
	conn, err := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("While trying to dial gRPC")
	}
	defer conn.Close()

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)
	// Schema
	// Install a schema into dgraph.
	err1 := dg.Alter(context.Background(), &api.Operation{
		Schema: `
			type:  		string @index(exact) . 
			buyer:  	uid .
			ip:  		string @index(term) .
			device: 	string . 
			products: 	[uid]  .
		`,
	})
	if err1 != nil {
		errorMessage := fmt.Sprintf("postProducts:Error setting schema => %v", err1)
		log.Fatal(errorMessage)
	}

	// Mutation
	mu := &api.Mutation{CommitNow: true}

	pb, err := json.Marshal(transactions_decoded)
	if err != nil {
		log.Fatal(err)
	}
	mu.SetJson = pb
	assigned, err := dg.NewTxn().Mutate(ctx, mu)
	if err != nil {
		log.Fatal(err)
	}
	//------------------------------

	// Save transactions Uids
	transactions = assigned.Uids
	// Debug
	// fmt.Printf("Assigned %v \n", transactions)

	// Debug
	// fmt.Printf("Suggestions => %v", suggestions)

	// Send succsessfull message
	message := fmt.Sprintf("PostTransactions: Added [%v] transactions, date [%v]", len(transactions_decoded), date)
	response.Write([]byte(message))
}

//--------------------- Transactions ---------------------
func test(response http.ResponseWriter, request *http.Request) {

	ctx := context.Background()
	conn, err := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("While trying to dial gRPC")
	}
	defer conn.Close()

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)

	// Create query
	const q = `	{
					t(func: eq(type, "Transaction")) {
						type
						name
						age
						
						price
						
						uid
							device
						ip
						buyer {
								uid
						name
						}
						products {
						uid
								name
						price
						}
					}
				}`

	// Ask for the type of name and age.
	resp, err := dg.NewTxn().Query(ctx, q)
	if err != nil {
		log.Fatal(err)
	}

	response.Write(resp.GetJson())
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

	//--------------   TRANSACTION endpoints   --------------
	router.Get("/test", test)

	//--------------   Listenning server   --------------
	fmt.Println("Server listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
