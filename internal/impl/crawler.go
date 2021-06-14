package impl

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

const crawlerTimeout = time.Minute

type Crawler struct {
	server *Server
}

func NewCrawler() *Crawler {
	return &Crawler{}
}

func (crawler *Crawler) Run() error {
	for {
		err := crawler.Parse()
		if err != nil {
			log.Printf("crawler error: %s", err)
			return err
		}

		time.Sleep(crawlerTimeout)
	}
}

func (crawler *Crawler) Parse() error {
	c := colly.NewCollector()

	c.OnHTML("div.posts_list", func(e *colly.HTMLElement) {
		e.ForEachWithBreak("article.post", func(_ int, a *colly.HTMLElement) bool {
			link := a.ChildAttr(".post__title a", "href")
			parsedLink := strings.Split(strings.TrimRight(link, "/"), "/")
			id, _ := strconv.Atoi(parsedLink[len(parsedLink)-1])
			post := Post{
				Preview: PostPreview{
					ID:          uint64(id),
					Title:       a.ChildText(".post__title a"),
					Author:      a.ChildText(".post__meta .user-info__nickname"),
					Link:        link,
					PublishedAt: a.ChildText(".post__time"),
				},
			}
			message := fmt.Sprintf("%s\n%s", post.Preview.Link, post.Preview.PublishedAt)
			_ = crawler.server.tg.SendMessageToAdmin(message)
			return false
		})
	})

	return c.Visit("https://habr.com/ru/hub/controllers/")
}
