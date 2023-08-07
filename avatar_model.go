package main

type AuthorEdge struct {
	Node Author `json:"node"`
}

type Author struct {
	Name      string      `json:"name"`
	FirstName string      `json:"firstName"`
	LastName  string      `json:"lastName"`
	Avatar    AvatarImage `json:"avatar"`
}
