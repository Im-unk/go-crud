package model

// Post represents a single post
type Post struct {
	ID     int    `json: "id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	Author User   `json:"author"`
}
