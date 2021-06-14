package impl

type PostPreview struct {
	ID          uint64
	Title       string
	Author      string
	Link        string
	PublishedAt string
}

type Post struct {
	Preview PostPreview
}
