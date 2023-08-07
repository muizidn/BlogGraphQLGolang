package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAllPostsQuery(t *testing.T) {
	// Create a test request with the GraphQL query
	query := `
	query AllPosts {
		posts(first: 20, where: { orderby: { field: DATE, order: DESC } }) {
			edges {
				node {
					title
					excerpt
					slug
					date
					featuredImage {
						node {
							sourceUrl
						}
					}
					author {
						node {
							name
							firstName
							lastName
							avatar {
								url
							}
						}
					}
				}
			}
		}
	}
	`
	reqBody := map[string]interface{}{
		"query": query,
	}
	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/graphql", bytes.NewBuffer(reqJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a recorder to capture the response
	rr := httptest.NewRecorder()
	http.HandlerFunc(http.HandlerFunc(graphqlHandler)).ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, but got %v", rr.Code)
	}

	// Decode and check the response JSON
	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	fmt.Println(response)

	data, dataExists := response["data"].(map[string]interface{})
	if !dataExists {
		t.Errorf("Expected 'data' field in response, but not found")
		return
	}

	posts, postsExists := data["posts"].(map[string]interface{})
	if !postsExists {
		t.Errorf("Expected 'posts' field in response data, but not found")
		return
	}

	edges, edgesExists := posts["edges"].([]interface{})
	if !edgesExists {
		t.Errorf("Expected 'edges' field in 'posts' data, but not found")
		return
	}

	if len(edges) == 0 {
		t.Errorf("Expected at least one post, but found none")
		return
	}

	// Check fields in each post
	for _, edge := range edges {
		post, postExists := edge.(map[string]interface{})["node"].(map[string]interface{})
		if !postExists {
			t.Errorf("Expected 'node' field in post edge, but not found")
			continue
		}

		// Check author field
		author, authorExists := post["author"].(map[string]interface{})["node"].(map[string]interface{})
		if !authorExists {
			t.Errorf("Expected 'author' field in post, but not found")
			continue
		}
		_, avatarExists := author["avatar"].(map[string]interface{})["url"].(string)
		if !avatarExists {
			t.Errorf("Expected 'avatar' field in author, but not found")
		}

		firstName, firstNameExists := author["firstName"].(string)
		if !firstNameExists {
			t.Errorf("Expected 'firstName' field in author, but not found")
		}

		lastName, lastNameExists := author["lastName"].(string)
		if !lastNameExists {
			t.Errorf("Expected 'lastName' field in author, but not found")
		}

		_, nameExists := author["name"].(string)
		if !nameExists {
			t.Errorf("Expected 'name' field in author, but not found")
		}

		// Check other fields
		date, dateExists := post["date"].(string)
		if !dateExists {
			t.Errorf("Expected 'date' field in post, but not found")
		}

		excerpt, excerptExists := post["excerpt"].(string)
		if !excerptExists {
			t.Errorf("Expected 'excerpt' field in post, but not found")
		}

		featuredImage, featuredImageExists := post["featuredImage"].(map[string]interface{})["node"].(map[string]interface{})
		if !featuredImageExists {
			t.Errorf("Expected 'featuredImage' field in post, but not found")
		}
		sourceUrl, sourceUrlExists := featuredImage["sourceUrl"].(string)
		if !sourceUrlExists {
			t.Errorf("Expected 'sourceUrl' field in featuredImage, but not found")
		}

		slug, slugExists := post["slug"].(string)
		if !slugExists {
			t.Errorf("Expected 'slug' field in post, but not found")
		}

		title, titleExists := post["title"].(string)
		if !titleExists {
			t.Errorf("Expected 'title' field in post, but not found")
		}

		// Print the checked fields for each post
		t.Logf("Post - Author: %s %s, Date: %s, Excerpt: %s, SourceUrl: %s, Slug: %s, Title: %s",
			firstName, lastName, date, excerpt, sourceUrl, slug, title)
	}
}
