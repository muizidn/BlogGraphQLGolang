package main

type ImageEdge struct {
	Node Image `json:"node"`
}

type Image struct {
	SourceURL string `json:"sourceUrl"`
}
