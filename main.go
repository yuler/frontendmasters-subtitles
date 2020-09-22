// refs: http://go-colly.org/docs/examples/coursera_courses/

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"cloud.google.com/go/translate"
	"github.com/gocolly/colly"
	"golang.org/x/text/language"
)

func main() {
	// host := "https://frontendmasters.com"
	c := colly.NewCollector(
		colly.AllowedDomains("frontendmasters.com"),
	)

	group := ""
	c.OnHTML(".CourseToc", func(t *colly.HTMLElement) {
		// Get lesson group
		t.ForEach(".LessonList", func(i int, g *colly.HTMLElement) {
			group = g.DOM.Prev().Text()
			// Get Course links
			g.ForEach(".CourseToc a[href]", func(i int, e *colly.HTMLElement) {
				link := e.Attr("href")
				e.Request.Visit(link)
			})
		})
	})

	// Get Lesson Transcript
	c.OnHTML(".LessonTranscript", func(e *colly.HTMLElement) {
		title := e.ChildText(".LessonTranscriptTitle")
		title = regexp.MustCompile("\"(.*?)\"").FindStringSubmatch(title)[1]
		transcripts := []string{}
		e.ForEach(".s-wrap p:not(:first-child)", func(i int, e *colly.HTMLElement) {
			transcripts = append(transcripts, e.Text)
		})
		fmt.Println("====")
		fmt.Printf("Group: %s \n", group)
		fmt.Printf("Title: %s \n", title)
		fmt.Println("Transcripts:")
		os.MkdirAll("output/"+group+"/", 0755)
		f, err := os.Create("output/" + group + "/" + title + ".srt")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		for i := 0; i < len(transcripts); i++ {
			transcript := transcripts[i]
			fmt.Println(transcript)

			start := transcript[1:9]
			end := ""
			if i+1 != len(transcripts) {
				end = transcripts[i+1][1:9]
			} else {
				endTime, _ := time.Parse("15:04:05", start)
				end = endTime.Add(100 * time.Second).Format("15:04:05")
			}

			f.WriteString(strconv.Itoa(i) + "\n")
			f.WriteString(start + ",000  -->  " + end + ",000\n")
			f.WriteString(transcript[10:] + "\n")
			translated := "xxxxxxxxxxxxxxx"
			for i := 0; i < 5; i++ {
				translated, err = translateText(transcript[10:])
				if err == nil {
					break
				}
			}

			fmt.Println(translated)
			f.WriteString(translated + "\n")
		}
		fmt.Println("====")
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.Visit("https://frontendmasters.com/courses/functional-javascript-v3/")
}

func translateText(text string) (string, error) {
	// text := "The Go Gopher is cute"
	ctx := context.Background()

	lang, err := language.Parse("zh")
	if err != nil {
		return "", fmt.Errorf("language.Parse: %v", err)
	}

	client, err := translate.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{text}, lang, nil)
	if err != nil {
		return "", fmt.Errorf("Translate: %v", err)
	}
	if len(resp) == 0 {
		return "", fmt.Errorf("Translate returned empty response to text: %s", text)
	}
	return resp[0].Text, nil
}
