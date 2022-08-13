package impl

import "time"

type PostBody struct {
	PublishedAt time.Time
	Content     string
}

type Post struct {
	ID          string
	Title       string
	Tags        []string
	Author      string
	Link        string
	PublishedAt string
}
