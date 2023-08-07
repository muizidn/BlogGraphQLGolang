package main

type PostEdge struct {
	Node   Post   `json:"node"`
	Cursor string `json:"cursor"`
}

type PostConnection struct {
	Edges    []PostEdge `json:"edges"`
	PageInfo PageInfo   `json:"pageInfo"`
}

type Post struct {
	ID            string     `json:"databaseId"`
	Slug          string     `json:"slug"`
	Status        string     `json:"status"`
	Title         string     `json:"title"`
	Excerpt       string     `json:"excerpt"`
	Date          string     `json:"date"`
	FeaturedImage ImageEdge  `json:"featuredImage"`
	Author        AuthorEdge `json:"author"`
}
