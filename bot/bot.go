package bot

import (
	"fmt"
	"log"
	"regexp"
	"reverse-phone-lookup/config"
	"reverse-phone-lookup/internal/scraper"
	"reverse-phone-lookup/utils"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var BotId string

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatalf("[-] Error occured starting bot..\n -> Error: %v\n", err)
	}

	u, err := goBot.User("@me")
	if err != nil {
		log.Fatalf("[-] Error occurred\n -> Error: %v\n", err)
	}

	BotId = u.ID

	goBot.AddHandler(messageHandler)

	err = goBot.Open()
	if err != nil {
		log.Fatalf("[-] Error occurred\n -> Error: %v\n", err)
	}
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotId {
		return
	}

	if strings.HasPrefix(m.Content, "<@"+BotId+">") {
		messageContent := strings.TrimSpace(strings.TrimPrefix(m.Content, "<@"+BotId+">"))

		number, fullName, state, city := parseArgs(messageContent)

		if number == "" && fullName == "" {
			s.ChannelMessageSend(m.ChannelID, "You need to specify at least one of phone or fullname.")
			return
		}

		var usphonebookLink string
		if fullName != "" {
			if state == "" {
				s.ChannelMessageSend(m.ChannelID, "State is required.. retry again and specify the state")
				return
			}
			encodedFullName := utils.EncodeArgs(fullName)
			if city != "" {
				encodedCity := utils.EncodeArgs(city)
				usphonebookLink = fmt.Sprintf("http://usphonebook.com/%v/%v/%v", encodedFullName, state, encodedCity)
			} else {
				usphonebookLink = fmt.Sprintf("https://usphonebook.com/%v/%v", encodedFullName, state)
			}
		} else if number != "" {
			usphonebookLink = fmt.Sprintf("https://usphonebook.com/%v", number)
		}

		targetLinks := scraper.PhonebookSearch(usphonebookLink)
		if len(targetLinks) == 0 {
			s.ChannelMessageSend(m.ChannelID, "No targets found on usphonebook.com")
			return
		} else if len(targetLinks) > 1 {
			targetLinks = targetLinks[1:]
		}

		persons := scraper.PhonebookTarget(targetLinks)

		for _, person := range persons {
			homeAddress := fmt.Sprintf("%s %s %s %s",
				person.HomeLocation.Address.StreetAddress,
				person.HomeLocation.Address.AddressRegion,
				person.HomeLocation.Address.AddressLocality,
				person.HomeLocation.Address.PostalCode)

			embed := &discordgo.MessageEmbed{
				Title:       "Person Information",
				Description: "Here are the details we found:",
				Color:       0x00ff00, // Green color for embed
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Home Location URL",
						Value:  person.HomeLocation.URL,
						Inline: false,
					},
					{
						Name:   "Name",
						Value:  person.Name,
						Inline: true,
					},
					{
						Name:   "Given Name",
						Value:  person.GivenName,
						Inline: true,
					},
					{
						Name:   "Family Name",
						Value:  person.FamilyName,
						Inline: true,
					},
					{
						Name:   "Current Address",
						Value:  homeAddress,
						Inline: false,
					},
					{
						Name:   "Previous Addresses",
						Value:  formatAddresses(person.Address),
						Inline: false,
					},
					{
						Name:   "Related Persons",
						Value:  formatRelatedPersons(person.RelatedTo),
						Inline: false,
					},
					{
						Name:   "Emails",
						Value:  strings.Join(person.Email, ", "),
						Inline: false,
					},
					{
						Name:   "Telephones",
						Value:  strings.Join(person.Telephone, ", "),
						Inline: false,
					},
				},
			}

			s.ChannelMessageSendEmbed(m.ChannelID, embed)
		}
	}
}

func formatAddresses(addresses []struct {
	AddressLocality string `json:"addressLocality"`
	AddressRegion   string `json:"addressRegion"`
	PostalCode      string `json:"postalCode"`
	StreetAddress   string `json:"streetAddress"`
}) string {
	var addrList []string
	for _, addr := range addresses {
		address := fmt.Sprintf("%s %s %s %s", addr.StreetAddress, addr.AddressRegion, addr.AddressLocality, addr.PostalCode)
		addrList = append(addrList, address)
	}
	return strings.Join(addrList, ", ")
}

func formatRelatedPersons(relatedTo []struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}) string {
	var relatedPersons []string
	for _, related := range relatedTo {
		relatedPersons = append(relatedPersons, related.Name)
	}
	return strings.Join(relatedPersons, ", ")
}

func parseArgs(messageContent string) (string, string, string, string) {
	number, fullName, state, city := "", "", "", ""

	// Regular expression to match key=value pairs, allowing spaces in values
	re := regexp.MustCompile(`(\w+)=([^=]+)`)
	matches := re.FindAllStringSubmatch(messageContent, -1)

	for _, match := range matches {
		key, value := match[1], strings.TrimSpace(match[2])
		switch key {
		case "phone":
			number = value
		case "fullname":
			fullName = value
		case "state":
			state = value
		case "city":
			city = value
		}
	}
	return number, fullName, state, city
}
