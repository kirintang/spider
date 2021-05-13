package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

func main() {
	t := time.Now()
	number := 1
	c := colly.NewCollector(func(c *colly.Collector) {
		extensions.RandomUserAgent(c)
		c.Async = true
	}, colly.URLFilters(regexp.MustCompile("^(https://movie\\.douban\\.com/top250)\\?start=[0-9].*&filter=")))
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(link))
	})
	c.OnHTML("div.info", func(e *colly.HTMLElement) {
		e.DOM.Each(func(i int, selection *goquery.Selection) {
			movies := selection.Find("span.title").First().Text()
			director := strings.Join(strings.Fields(selection.Find("div.db p").First().Text()), " ")
			quote := selection.Find("p.quote span.inq").Text()
			fmt.Printf("%d -> %s:%s %s\n", number, movies, director, quote)
			number += 1
		})
	})
	c.OnError(func(response *colly.Response, err error) {
		fmt.Println(err)
	})
	c.Visit("https://movie.douban.com/top250?start=0&filter=")
	c.Wait()
	fmt.Printf("cost: %s", time.Since(t))
}
