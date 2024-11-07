package impl

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type Crawler struct {
	server        *Server
	urls          []string
	parserTimeout time.Duration
}

func NewCrawler(server *Server) *Crawler {
	return &Crawler{
		server:        server,
		urls:          server.config.Parser.URLS,
		parserTimeout: time.Duration(server.config.Parser.Timeout),
	}
}

func (c *Crawler) InitCrawler(ctx context.Context) error {
	var posts []Post
	for _, url := range c.urls {
		p, err := c.GetPosts(url)
		if err != nil {
			return err
		}
		posts = append(posts, p...)
	}
	_, err := c.SavePosts(ctx, posts)
	return err
}

func (c *Crawler) Run(ctx context.Context) error {
	i := 0
	errCount := 0
	for {
		err := c.Parser(ctx, c.urls[i])
		if err != nil {
			errCount++
			if errCount > 2 {
				err = c.server.tg.SendMessageToAdmin(err.Error())
				if err != nil {
					log.Printf("crawler error: %s", err)
				}
			}
		} else {
			errCount = 0
		}
		i++
		if i >= len(c.urls) {
			i = 0
		}

		time.Sleep(c.parserTimeout * time.Minute)
	}
}

func (c *Crawler) GetPosts(url string) ([]Post, error) {
	var posts []Post

	col := colly.NewCollector()
	col.OnHTML("div.tm-articles-list", func(e *colly.HTMLElement) {
		e.ForEach("article.tm-articles-list__item", func(_ int, a *colly.HTMLElement) {
			var tags []string
			a.ForEach(".tm-publication-hub__link", func(_ int, p *colly.HTMLElement) {
				if !strings.HasPrefix(p.Text, "Блог") {
					p.Text = strings.Trim(p.Text, "* ")
					tags = append(tags, p.Text)
				}
			})
			post := Post{
				ID:          a.Attr("id"),
				Title:       a.ChildText(".tm-title__link"),
				Author:      a.ChildText(".tm-user-info__username"),
				Tags:        tags,
				Link:        "habr.com" + a.ChildAttr(".tm-title__link", "href"),
				PublishedAt: a.ChildAttr(".tm-article-datetime-published>time", "title"),
			}

			posts = append(posts, post)
		})
	})

	col.SetRequestTimeout(30 * time.Second)
	err := col.Visit(url)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (c *Crawler) Parser(ctx context.Context, url string) error {
	posts, err := c.GetPosts(url)
	if err != nil {
		return err
	}
	newPosts, err := c.SavePosts(ctx, posts)
	if err != nil {
		return err
	}
	if len(newPosts) > 0 {
		//fmt.Println("test")
		c.server.tg.SendPostsToChannel(ctx, newPosts)
	}
	return nil
}

func (c *Crawler) SavePosts(ctx context.Context, posts []Post) ([]Post, error) {
	var newPosts []Post
	for _, post := range posts {
		if post.Link == "" {
			continue
		}
		item, err := c.server.db.Get(ctx, post.ID)
		if err != nil {
			return nil, err
		}
		if item == nil {
			var bPost []byte
			bPost, err = json.Marshal(post)
			if err != nil {
				return nil, err
			}
			err = c.server.db.Set(ctx, post.ID, bPost)
			if err != nil {
				return nil, err
			}
			newPosts = append(newPosts, post)
		}
	}
	return newPosts, nil
}
