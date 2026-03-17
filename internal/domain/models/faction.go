package postgres

var FactionNames = map[string]string{
	"vincent-faction": "vincent",
	"david-faction":   "david",
	"jimmy-faction":   "jimmy",
}

type Faction struct {
	ID          int
	Name        string
	Description string

	Admins []Account
}
