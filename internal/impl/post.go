package impl

import "time"

type PostPreview struct {
	ID          uint64
	Title       string
	Author      string
	Link        string
	PublishedAt string
}

type PostBody struct {
	PublishedAt time.Time
}

type Post struct {
	Preview PostPreview
	Body    PostBody
}
