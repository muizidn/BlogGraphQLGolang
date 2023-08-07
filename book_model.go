package main

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	// Add other fields as needed
}

type BookEdge struct {
	Node   Book   `json:"node"`
	Cursor string `json:"cursor"`
}

type BookConnection struct {
	Edges    []BookEdge `json:"edges"`
	PageInfo PageInfo   `json:"pageInfo"`
}
