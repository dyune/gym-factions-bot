package postgres

import (
	_ "github.com/uptrace/bun"
)

type AccountType int

const (
	Admin AccountType = iota
	Member
)

var AdminList = map[string]bool{
	"ketchup.383": true,
	"duune":       true,
}

type Account struct {
	ID             int
	Name           string
	AccountType    AccountType
	Points         int
	FactionID      int
	AdminFactionID int

	Faction *Faction
}
