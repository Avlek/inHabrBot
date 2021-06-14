package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

const crawlerTimeout = 3 * time.Minute

type Crawler struct {
	server *Server
}

func NewCrawler(server *Server) *Crawler {
	return &Crawler{
		server: server,
	}
}

func (crawler *Crawler) Run(ctx context.Context) error {
	for {
		err := crawler.ParseList(ctx)
		if err != nil {
			log.Printf("crawler error: %s", err)
			return err
		}

		time.Sleep(crawlerTimeout)
	}
}

func (crawler *Crawler) ParseList(ctx context.Context) error {
	c := colly.NewCollector()

	c.OnHTML("div.posts_list", func(e *colly.HTMLElement) {
		e.ForEach("article.post", func(_ int, a *colly.HTMLElement) {
			link := a.ChildAttr(".post__title a", "href")
			parsedLink := strings.Split(strings.TrimRight(link, "/"), "/")
			sid := parsedLink[len(parsedLink)-1]
			id, _ := strconv.Atoi(sid)

			post := Post{
				Preview: PostPreview{
					ID:          uint64(id),
					Title:       a.ChildText(".post__title a"),
					Author:      a.ChildText(".post__meta .user-info__nickname"),
					Link:        link,
					PublishedAt: a.ChildText(".post__time"),
				},
			}
			item, err := crawler.server.db.Get(ctx, sid)
			if err != nil {
				log.Println("redis get error:", err)
			}
			if item == nil {
				body := crawler.ParsePost(link)
				if body.PublishedAt.Before(time.Now().Add(-5 * time.Hour)) {
					return
				}

				bPost, err := json.Marshal(post)
				if err != nil {
					log.Println("redis set error:", err)
				}
				err = crawler.server.db.Set(ctx, sid, bPost)
				if err != nil {
					log.Println("redis set error:", err)
				}
				message := fmt.Sprintf("%s\n%s", post.Preview.Link, post.Preview.PublishedAt)
				_ = crawler.server.tg.SendMessageToAdmin(message)
			}
		})
	})

	return c.Visit("https://habr.com/ru/all/top25/")
}

func (crawler *Crawler) ParsePost(url string) *PostBody {
	c := colly.NewCollector()
	body := &PostBody{}
	c.OnHTML("article.post", func(e *colly.HTMLElement) {
		str := e.ChildAttr(".post__time", "data-time_published")
		body.PublishedAt, _ = time.Parse("2006-01-02T15:04Z", str)
	})

	_ = c.Visit(url)
	return body
}
