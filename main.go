package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/graphql-go/graphql"
)

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books []Book

func init() {
	// Sample data for demonstration purposes
	books = []Book{
		{ID: "1", Title: "Book 1", Author: "Author 1"},
		{ID: "2", Title: "Book 2", Author: "Author 2"},
		{ID: "3", Title: "Book 3", Author: "Author 3"},
	}
}

var bookType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Book",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"author": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"book": &graphql.Field{
				Type:        bookType,
				Description: "Get a single book by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, _ := params.Args["id"].(string)
					for _, book := range books {
						if book.ID == id {
							return book, nil
						}
					}
					return nil, nil
				},
			},
			"books": &graphql.Field{
				Type:        graphql.NewList(bookType),
				Description: "Get all books",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return books, nil
				},
			},
		},
	},
)

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: queryType,
	},
)

func graphqlHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Parse the JSON request body
	var requestBody map[string]interface{}
	if err := json.Unmarshal(body, &requestBody); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Get the GraphQL query from the request body
	query, ok := requestBody["query"].(string)
	if !ok {
		http.Error(w, "Query not found in request body", http.StatusBadRequest)
		return
	}
	params := graphql.Params{Schema: schema, RequestString: query}
	result := graphql.Do(params)
	json.NewEncoder(w).Encode(result)
}

func main() {
	http.HandleFunc("/graphql", graphqlHandler)
	fmt.Println("GraphQL server is running on http://localhost:8080/graphql")
	http.ListenAndServe(":8080", nil)
}
