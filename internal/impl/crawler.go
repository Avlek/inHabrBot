package impl

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

type Crawler struct {
}

func NewCrawler() *Crawler {
	return &Crawler{}
}

func (crawler *Crawler) Run(quit chan error) {
	for {
		time.After(time.Second)

		err := crawler.Parse()
		if err != nil {
			log.Errorf("crawler error: %s", err)
			quit <- err
		}

		quit <- nil
	}
}

func (crawler *Crawler) Parse() error {
	c := colly.NewCollector()

	c.OnHTML("div.posts_list", func(e *colly.HTMLElement) {
		e.ForEach("a.post__title_link", func(_ int, a *colly.HTMLElement) {
			fmt.Println(a.Attr("href"))
		})
	})

	return c.Visit("https://habr.com/ru/all/")
}
