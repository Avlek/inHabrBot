package impl

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gocolly/colly"
)

type Crawler struct {
	server        *Server
	url           string
	version       int8
	parserTimeout time.Duration
}

func NewCrawler(server *Server) *Crawler {
	return &Crawler{
		server:        server,
		url:           server.config.Parser.URL,
		version:       server.config.Parser.Version,
		parserTimeout: time.Duration(server.config.Parser.Timeout),
	}
}

func (c *Crawler) InitCrawler(ctx context.Context) error {
	posts, err := c.GetPosts(c.version)
	if err != nil {
		return err
	}

	_, err = c.SavePosts(ctx, posts)
	return err
}

func (c *Crawler) Run(ctx context.Context) error {
	for {
		err := c.Parser(ctx)
		if err != nil {
			err = c.server.tg.SendMessageToAdmin(err.Error())
			if err != nil {
				log.Printf("crawler error: %s", err)
			}
			time.Sleep(10 * time.Minute)
		}

		time.Sleep(c.parserTimeout * time.Minute)
	}
}

func (c *Crawler) GetPosts(version int8) ([]Post, error) {
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
				post := Post{
					ID:          a.Attr("id"),
					Title:       a.ChildText(".tm-article-snippet__title-link"),
					Author:      a.ChildText(".tm-user-info__user"),
					Link:        "habr.com" + a.ChildAttr(".tm-article-snippet__title-link", "href"),
					PublishedAt: a.ChildText(".tm-article-snippet__datetime-published"),
				}

				posts = append(posts, post)
			})
		})
	}

	err := col.Visit(c.url)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (c *Crawler) Parser(ctx context.Context) error {

	posts, err := c.GetPosts(c.version)
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
