package store

import (
	"github.com/uptrace/bun"
)

var FactionNames = map[string]string{
	"vincent-faction": "vincent",
	"david-faction":   "david",
	"jimmy-faction":   "jimmy",
}

type Faction struct {
	bun.BaseModel `bun:"table:factions"`
	ID            int    `bun:",pk,autoincrement"`
	Name          string `bun:",notnull"`
	Description   string

	Admins []User `bun:"rel:has-many,join:id=admin_faction_id"`
}
