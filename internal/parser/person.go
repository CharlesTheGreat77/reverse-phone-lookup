package parser

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