package scraper

import (
	"fmt"
	"reverse-phone-lookup/internal/parser"
	"strings"
	"encoding/json"

	"github.com/gocolly/colly/v2"
)

// function to search usphonebook and scrape urls of the results (some are older results)
func PhonebookSearch(usphonebookLink string) []string {
	targetLinks := []string{}
	c := createCollector()
	setCollyBehavior(c)

	c.OnHTML("input[name='link']", func(e *colly.HTMLElement) {
		link := e.Attr("value")
		targetLinks = append(targetLinks, e.Request.AbsoluteURL(link))
	})

	c.OnHTML("a.ls_contacts-btn[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		targetLinks = append(targetLinks, e.Request.AbsoluteURL(link))
	})

	c.Visit(usphonebookLink)

	return targetLinks
}

// function to scrape each url for a given query
func PhonebookTarget(targetLinks []string) []parser.Person {
	var persons []parser.Person
	c := createCollector()
	setCollyBehavior(c)

	c.OnHTML("script[type='application/ld+json']", func(e *colly.HTMLElement) {
		scriptText := e.Text
		if strings.Contains(scriptText, "@type") {
			var person parser.Person
			err := json.Unmarshal([]byte(scriptText), &person)
			if err != nil {
				fmt.Printf("[-] Error parsing JSON: %v\n", err)
				return
			}
			if person.GivenName != "" {
				persons = append(persons, person)
			}
		}
	})

	for _, url := range targetLinks {
		c.Visit(url)
	}
	return persons
}