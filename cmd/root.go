package cmd

import (
	"flag"
	"fmt"
	"reverse-phone-lookup/internal/scraper"
	"reverse-phone-lookup/internal/parser"
	"reverse-phone-lookup/utils"
	"strings"
)

func Execute() {
	number := flag.String("phone", "", "specify a phone number [777-999-0000]")
	fullName := flag.String("fullname", "", "specify the target's full name [John Doe]")
	state := flag.String("state", "", "specify the state the target resides [California]")
	city := flag.String("city", "", "specify the city the target resides [Los Angeles]")
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
	var persons []parser.Person

	if *fullName != "" {
		encodedFullName := utils.EncodeArgs(*fullName)
		if *state == "" {
			fmt.Printf("[-] State is required.. retry again and specify the state\n")
			flag.Usage()
			return
		}
		if *city != "" {
			encodedCity := utils.EncodeArgs(*city)
			usphonebookLink = fmt.Sprintf("http://usphonebook.com/%v/%v/%v", encodedFullName, *state, encodedCity)
		} else {
			usphonebookLink = fmt.Sprintf("https://usphonebook.com/%v/%v", encodedFullName, *state)
		}
	}

	if *number != "" {
		usphonebookLink = fmt.Sprintf("https://usphonebook.com/%v", *number)
	}

	targetLinks = scraper.PhonebookSearch(usphonebookLink)
	if len(targetLinks) == 0 {
		fmt.Printf("[-] No targets found on usphonebook.com\n")
		return
	} else if len(targetLinks) > 1 {
		targetLinks = targetLinks[1:]
	}

	persons = scraper.PhonebookTarget(targetLinks)

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