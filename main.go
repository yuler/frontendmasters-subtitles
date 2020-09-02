// refs: http://go-colly.org/docs/examples/coursera_courses/

package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/gocolly/colly"
)

func main() {
	// host := "https://frontendmasters.com"
	c := colly.NewCollector(
		colly.AllowedDomains("frontendmasters.com"),
	)

	// Get Course links
	c.OnHTML(".CourseToc a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		e.Request.Visit(link)
	})

	// Get Lesson Transcript
	c.OnHTML(".LessonTranscript", func(e *colly.HTMLElement) {
		title := e.ChildText(".LessonTranscriptTitle")
		title = regexp.MustCompile("\"(.*?)\"").FindStringSubmatch(title)[1]
		transcripts := e.ChildTexts(".s-wrap p:not(:first-child)")
		fmt.Println("====")
		fmt.Printf("Title: %s \n", title)
		fmt.Println("Transcripts:")
		f, err := os.Create("output/" + title + ".txt")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		for i := 0; i < len(transcripts); i++ {
			fmt.Println(transcripts[i])
			f.WriteString(transcripts[i] + "\n")
		}
		fmt.Println("====")
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.Visit("https://frontendmasters.com/courses/deep-javascript-v3/")
}
