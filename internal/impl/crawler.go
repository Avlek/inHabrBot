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
	version       int8
	parserTimeout time.Duration
}

func NewCrawler(server *Server) *Crawler {
	return &Crawler{
		server:        server,
		urls:          server.config.Parser.URLS,
		version:       server.config.Parser.Version,
		parserTimeout: time.Duration(server.config.Parser.Timeout),
	}
}

func (c *Crawler) InitCrawler(ctx context.Context) error {
	var posts []Post
	for _, url := range c.urls {
		p, err := c.GetPosts(c.version, url)
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

func (c *Crawler) GetPosts(version int8, url string) ([]Post, error) {
	var posts []Post

	col := colly.NewCollector()
	if version == int8(1) {
		col.OnHTML("div.posts_list", func(e *colly.HTMLElement) {
			e.ForEach("li.content-list__item_post", func(_ int, a *colly.HTMLElement) {
				post := Post{
					ID:          a.Attr("id"),
					Title:       a.ChildText(".post__title a"),
					Author:      a.ChildText(".post__meta .user-info__nickname"),
					Link:        a.ChildAttr(".post__title a", "href"),
					PublishedAt: a.ChildText(".post__time"),
				}

				posts = append(posts, post)
			})
		})
	} else if version == int8(2) {
		col.OnHTML("div.tm-articles-list", func(e *colly.HTMLElement) {
			e.ForEach("article.tm-articles-list__item", func(_ int, a *colly.HTMLElement) {
				var tags []string
				a.ForEach(".tm-article-snippet__hubs-item-link", func(_ int, p *colly.HTMLElement) {
					if !strings.HasPrefix(p.Text, "Блог") {
						p.Text = strings.Trim(p.Text, "* ")
						tags = append(tags, p.Text)
					}
				})
				post := Post{
					ID:          a.Attr("id"),
					Title:       a.ChildText(".tm-article-snippet__title-link"),
					Author:      a.ChildText(".tm-user-info__user"),
					Tags:        tags,
					Link:        "habr.com" + a.ChildAttr(".tm-article-snippet__title-link", "href"),
					PublishedAt: a.ChildText(".tm-article-snippet__datetime-published"),
				}

				posts = append(posts, post)
			})
		})
	}

	col.SetRequestTimeout(30 * time.Second)
	err := col.Visit(url)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (c *Crawler) Parser(ctx context.Context, url string) error {
	posts, err := c.GetPosts(c.version, url)
	if err != nil {
		return err
	}
	newPosts, err := c.SavePosts(ctx, posts)
	if err != nil {
		return err
	}
	if len(newPosts) > 0 {
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
