package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/graphql-go/graphql"
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
