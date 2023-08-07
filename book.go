package main

import (
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

var bookQueryType = graphql.NewObject(
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
