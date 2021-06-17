package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"google.golang.org/grpc"
)

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
