package main

import (
	"strconv"

	"github.com/graphql-go/graphql"
)

var postType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Post",
		Fields: graphql.Fields{
			"databaseId": &graphql.Field{
				Type: graphql.String,
			},
			"slug": &graphql.Field{
				Type: graphql.String,
			},
			"status": &graphql.Field{
				Type: graphql.String,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"excerpt": &graphql.Field{
				Type: graphql.String,
			},
			"date": &graphql.Field{
				Type: graphql.String,
			},
			"featuredImage": &graphql.Field{
				Type: imageEdgeType,
			},
			"author": &graphql.Field{
				Type: authorEdgeType,
			},
		},
	},
)

var imageType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Image",
		Fields: graphql.Fields{
			"sourceUrl": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var imageEdgeType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "ImageNode",
		Fields: graphql.Fields{
			"node": &graphql.Field{
				Type: imageType,
			},
		},
	},
)

var authorType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Author",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"firstName": &graphql.Field{
				Type: graphql.String,
			},
			"lastName": &graphql.Field{
				Type: graphql.String,
			},
			"avatar": &graphql.Field{
				Type: avatarImageType,
			},
		},
	},
)

var authorEdgeType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "AuthorNode",
		Fields: graphql.Fields{
			"node": &graphql.Field{
				Type: authorType,
			},
		},
	},
)

var avatarImageType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "AvatarImage",
		Fields: graphql.Fields{
			"url": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var postEdgeType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "postEdge",
		Fields: graphql.Fields{
			"node": &graphql.Field{
				Type: postType,
			},
			"cursor": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var postConnectionType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "postConnection",
		Fields: graphql.Fields{
			"edges": &graphql.Field{
				Type: graphql.NewList(postEdgeType),
			},
			"pageInfo": &graphql.Field{
				Type: pageInfoType,
			},
		},
	},
)

var orderByFieldEnum = graphql.NewEnum(graphql.EnumConfig{
	Name: "OrderByField",
	Values: graphql.EnumValueConfigMap{
		"DATE": &graphql.EnumValueConfig{
			Value: "DATE",
		},
		// Add more enum values here if needed
	},
})

// Define enum values for OrderDirection
var orderDirectionEnum = graphql.NewEnum(graphql.EnumConfig{
	Name: "OrderDirection",
	Values: graphql.EnumValueConfigMap{
		"ASC": &graphql.EnumValueConfig{
			Value: "ASC",
		},
		"DESC": &graphql.EnumValueConfig{
			Value: "DESC",
		},
	},
})

var orderByInputType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "OrderByInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"field": &graphql.InputObjectFieldConfig{
				Type: orderByFieldEnum,
			},
			"order": &graphql.InputObjectFieldConfig{
				Type: orderDirectionEnum,
			},
		},
	},
)

var whereInputType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "WhereInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"orderby": &graphql.InputObjectFieldConfig{
				Type: orderByInputType,
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"posts": &graphql.Field{
				Type: postConnectionType,
				Args: graphql.FieldConfigArgument{
					"first": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"after": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"where": &graphql.ArgumentConfig{
						Type: whereInputType,
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

					sortedPosts, _ := orderPostsBy(&params)

					var hasNextPage bool
					var endCursor string
					postEdges := []PostEdge{}

					for i := offset; i < offset+firstArg && i < len(sortedPosts); i++ {
						postEdges = append(postEdges, PostEdge{
							Node:   sortedPosts[i],
							Cursor: strconv.Itoa(i),
						})
					}
					if offset+firstArg < len(sortedPosts) {
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

					return PostConnection{
						Edges:    postEdges,
						PageInfo: pageInfo,
					}, nil
				},
			},
		},
	},
)

func orderPostsBy(params *graphql.ResolveParams) ([]Post, error) {
	// Extract the "where" argument
	where, _ := params.Args["where"].(map[string]interface{})

	// Extract the "orderby" argument
	orderBy, _ := where["orderby"].(map[string]interface{})
	field, _ := orderBy["field"]
	order, _ := orderBy["order"]

	var sortedPosts []Post
	sortedPosts = append(sortedPosts, posts...)

	if field == "DATE" && order == "DESC" {
		// Implement sorting logic here
	}
	return sortedPosts, nil
}

// Sample data for demonstration purposes
var posts = []Post{
	{
		ID:      "1",
		Slug:    "post-1",
		Status:  "published",
		Title:   "Post 1",
		Excerpt: "This is the excerpt of Post 1",
		Date:    "2023-08-15",
		FeaturedImage: ImageEdge{
			Node: Image{SourceURL: "https://example.com/image1.jpg"},
		},
		Author: AuthorEdge{
			Node: Author{
				Name:      "Author 1",
				FirstName: "John",
				LastName:  "Doe",
				Avatar:    AvatarImage{URL: "https://example.com/avatar1.jpg"},
			},
		},
	},
	// Add more sample posts here
}
