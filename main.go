package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("frontendmasters.com"),
	)

	c.OnHTML(".CourseToc a[href]", func (e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
	})

	c.OnRequest(func(r *colly.Request) {
      fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("https://frontendmasters.com/courses/deep-javascript-v3/")
}
