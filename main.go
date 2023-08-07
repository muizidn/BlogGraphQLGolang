package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/graphql-go/graphql"
)

var books []Book

func init() {
	// Sample data for demonstration purposes
	books = []Book{
		{ID: "1", Title: "Book 1", Author: "Author 1"},
		{ID: "2", Title: "Book 2", Author: "Author 2"},
		{ID: "3", Title: "Book 3", Author: "Author 3"},
		{ID: "4", Title: "Book 4", Author: "Author 4"},
		{ID: "5", Title: "Book 5", Author: "Author 5"},
		// Add more sample books here
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
			// Add other fields as needed
		},
	},
)

var bookEdgeType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "BookEdge",
		Fields: graphql.Fields{
			"node": &graphql.Field{
				Type: bookType,
			},
			"cursor": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var bookConnectionType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "BookConnection",
		Fields: graphql.Fields{
			"edges": &graphql.Field{
				Type: graphql.NewList(bookEdgeType),
			},
			"pageInfo": &graphql.Field{
				Type: pageInfoType,
			},
		},
	},
)

var pageInfoType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "PageInfo",
		Fields: graphql.Fields{
			"hasNextPage": &graphql.Field{
				Type: graphql.Boolean,
			},
			"endCursor": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"books": &graphql.Field{
				Type: bookConnectionType,
				Args: graphql.FieldConfigArgument{
					"first": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"after": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					firstArg, _ := params.Args["first"].(int)
					afterArg, _ := params.Args["after"].(string)

					// Implement pagination logic here based on 'firstArg' and 'afterArg'
					// For simplicity, we'll use a simple offset-based pagination
					var offset int
					if afterArg != "" {
						offset, _ = strconv.Atoi(afterArg)
					}

					var hasNextPage bool
					var endCursor string
					bookEdges := []BookEdge{}
					for i := offset; i < offset+firstArg && i < len(books); i++ {
						bookEdges = append(bookEdges, BookEdge{
							Node:   books[i],
							Cursor: strconv.Itoa(i),
						})
					}
					if offset+firstArg < len(books) {
						hasNextPage = true
						endCursor = strconv.Itoa(offset + firstArg)
					} else {
						hasNextPage = false
						endCursor = ""
					}

					pageInfo := PageInfo{
						HasNextPage: hasNextPage,
						EndCursor:   endCursor,
					}

					return BookConnection{
						Edges:    bookEdges,
						PageInfo: pageInfo,
					}, nil
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
