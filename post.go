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
					postEdges := []PostEdge{}
					for i := offset; i < offset+firstArg && i < len(posts); i++ {
						postEdges = append(postEdges, PostEdge{
							Node:   posts[i],
							Cursor: strconv.Itoa(i),
						})
					}
					if offset+firstArg < len(posts) {
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
