package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/davidwang/factions/internal/exceptions"
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

func ExistsByName(name string, ctx context.Context, db *bun.DB) bool {
	exists, err := db.NewSelect().
		Model((*User)(nil)).
		Where("name = ?", name).
		Exists(ctx)

	if err != nil {
		// Log error or handle it depending on your application's error handling strategy
		// For now, we'll assume that any error means the record doesn't exist
		return false
	}

	return exists
}

func InsertFaction(ctx context.Context, db *bun.DB, name string, description string) error {

	faction := &Faction{
		Name:        name,
		Description: description,
	}

	if ExistsByName(name, ctx, db) {
		return exceptions.ValidationError{
			Field: "name",
			Value: name,
			Msg:   "a faction with this name already exists",
		}
	}

	_, err := db.NewInsert().
		Model(faction).
		Exec(ctx)

	if err != nil {
		return exceptions.DatabaseError{
			Operation: "insert",
			Err:       err,
		}
	}

	return nil
}

func LookUpFaction(ctx context.Context, db *bun.DB, name string) (int, error) {
	faction := new(Faction)
	err := db.NewSelect().
		Model(faction).
		Where("name = ?", name).
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return -1, err

	} else if err != nil {
		return -2, exceptions.DatabaseError{
			Operation: "insert",
			Err:       err,
		}

	} else {
		return faction.ID, nil
	}
}
