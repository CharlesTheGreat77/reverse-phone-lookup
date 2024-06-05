package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"flag"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/imroc/req/v3"
)


// struct for each Person found in results
type Person struct {
	Name         string   `json:"name"`
	GivenName    string   `json:"givenName"`
	FamilyName   string   `json:"familyName"`
	HomeLocation struct {
		URL          string `json:"url"`
		Description  string `json:"description"`
		Address      struct {
			AddressLocality string `json:"addressLocality"`
			AddressRegion   string `json:"addressRegion"`
			PostalCode      string `json:"postalCode"`
			StreetAddress   string `json:"streetAddress"`
		} `json:"address"`
	} `json:"homeLocation"`
	Address   []struct {
		AddressLocality string `json:"addressLocality"`
		AddressRegion   string `json:"addressRegion"`
		PostalCode      string `json:"postalCode"`
		StreetAddress   string `json:"streetAddress"`
	} `json:"address"`
	RelatedTo []struct {
		URL  string `json:"url"`
		Name string `json:"name"`
	} `json:"relatedTo"`
	Email     []string `json:"email"`
	Telephone []string `json:"telephone"`
}


// function to create a collector that impersonates a chrome browser
func createCollector() *colly.Collector {
	fakeChrome := req.DefaultClient().ImpersonateChrome()
	c := colly.NewCollector(
		colly.UserAgent(fakeChrome.Headers.Get("user-agent")),
	   )
	c.SetClient(&http.Client{Transport: fakeChrome.Transport,})
	return c
}


// function to encode args for search query
func encode_args(args string) string {
	parts := strings.Split(args, " ")
	return strings.Join(parts, "-")
}


// main function to scrape and output result
func main() {
	number := flag.String("phone", "", "specify a phone number [777-999-0000]")
	fullName := flag.String("fullname", "", "specify the targets full name [John Doe]")
	state := flag.String("state", "", "specify the state the target resides [California]")
	city := flag.String("city", "", "specify the city the target resides [Los Angelos]")
	help := flag.Bool("h", false, "show usage")
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}

	if *number == "" && *fullName == "" {
		flag.Usage()
		return
	}

	var usphonebookLink string
	var targetLinks []string
	var persons []Person

	if *fullName != "" {
		encodedFullName := encode_args(*fullName)
		if *state == "" {
			fmt.Printf("[-] State is required.. retry again and specify the state\n")
			flag.Usage()
			return
		}
		if *city != "" {
			encodedCity := encode_args(*city)
			usphonebookLink = fmt.Sprintf("http://usphonebook.com/%v/%v/%v", encodedFullName, *state, encodedCity)
		} else {
			usphonebookLink = fmt.Sprintf("https://usphonebook.com/%v/%v", encodedFullName, *state)
		}
	}

	if *number != "" {
		usphonebookLink = fmt.Sprintf("https://usphonebook.com/%v", *number)
	}

	targetLinks = phonebookSearch(usphonebookLink)
	if len(targetLinks) == 0 {
		fmt.Printf("[-] No targets found on usphonebook.com\n")
		return
	} else if len(targetLinks) > 1 {
		targetLinks = targetLinks[1:] // first one is usually the same link we initially scraped (not always the case)
	}

	persons = phonebookTarget(targetLinks)

	for _, person := range persons {
		
		fmt.Println("Home Location URL:", person.HomeLocation.URL)
		fmt.Println("Name:", person.Name)
		fmt.Println("Given Name:", person.GivenName)
		fmt.Println("Family Name:", person.FamilyName)
		homeAddress := fmt.Sprintf("%s %s %s %s", person.HomeLocation.Address.StreetAddress, person.HomeLocation.Address.AddressRegion, person.HomeLocation.Address.AddressLocality, person.HomeLocation.Address.PostalCode)
		fmt.Println("Current Address:", homeAddress)

		var addresses []string
		for _, addr := range person.Address {
			address := fmt.Sprintf("%s %s %s %s", addr.StreetAddress, addr.AddressRegion, addr.AddressLocality, addr.PostalCode)
			addresses = append(addresses, address)
		}
		fmt.Println("Previous Addresses:", strings.Join(addresses, ", "))

		var relatedPersons []string
		for _, related := range person.RelatedTo {
			relatedPersons = append(relatedPersons, related.Name)
		}
		fmt.Println("Related Persons:", strings.Join(relatedPersons, ", "))
		fmt.Println("Emails:", strings.Join(person.Email, ", "))
		fmt.Printf("Telephones: %v\n\n", strings.Join(person.Telephone, ", "))
	}
}


// function to search usphonebook and scrape urls of the results (some are older results)
func phonebookSearch(usphonebookLink string) []string {
	targetLinks := []string{}
	c := createCollector()

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
		fmt.Println("[*] Sending request to:", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Printf("[*] Response Code: %v\n", r.StatusCode)
	})

	c.OnHTML("input[name='link']", func(e *colly.HTMLElement) {
		link := e.Attr("value")
		targetLinks = append(targetLinks, e.Request.AbsoluteURL(link))
	})

	c.OnHTML("a.ls_contacts-btn[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		targetLinks = append(targetLinks, e.Request.AbsoluteURL(link))
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Printf("[-] Error Requesting URL: %v\n", err)
	})
	
	c.Visit(usphonebookLink)

	return targetLinks
}


// function to scrape each url for a given query
func phonebookTarget(targetLinks []string) []Person {
	var persons []Person
	c := createCollector()
	
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
		fmt.Println("[*] Sending request to:", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Printf("[*] Response Code: %v\n", r.StatusCode)
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Printf("[-] Error Requesting URL: %v\n", err)
	})

	c.OnHTML("script[type='application/ld+json']", func(e *colly.HTMLElement) {
		scriptText := e.Text
		if strings.Contains(scriptText, "@type") {
			var person Person
			err := json.Unmarshal([]byte(scriptText), &person) // love this
			if err != nil {
				fmt.Println("[-] Error parsing JSON:", err)
				return // exit if error occurs
			}
			if person.GivenName != "" {
				persons = append(persons, person)

			}
		}
	})
	// loop through each link
	for _, url := range targetLinks {
		c.Visit(url)
	}
	return persons
}