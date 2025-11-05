package store

import (
	"context"
	"github.com/davidwang/factions/internal/exceptions"
	"github.com/uptrace/bun"
	_ "github.com/uptrace/bun"
	"log"
)

type UserType int

const (
	Admin UserType = iota
	Member
)

var AdminList = map[string]bool{
	"ketchup.383": true,
	"duune":       true,
}

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID             int `bun:",pk"`
	Name           string
	DiscordName    string   `bun:",notnull"`
	UserType       UserType `bun:",default:1"`
	Points         int      `bun:",notnull,default:0"`
	FactionID      int
	AdminFactionID int `bun:"admin_faction_id"`

	Faction *Faction `bun:"rel:belongs-to,join:faction_id=id"`
}

func ExistsByID(id int, ctx context.Context, db *bun.DB) bool {
	exists, err := db.NewSelect().
		Model((*User)(nil)).
		Where("id = ?", id).
		Exists(ctx)

	if err != nil {
		log.Printf("[ERROR] while calling ExistsByID: %v", err)
	}
	if exists {
		return true
	}

	return false
}

func InsertUser(id int, username string, name string, ctx context.Context, db *bun.DB) (int, error) {

	if ExistsByID(id, ctx, db) {
		return -1, exceptions.UserExistsError{
			UserID:   id,
			Username: username,
		}
	}

	var newUser *User
	if AdminList[username] {
		newUser = &User{
			ID:          id,
			Name:        name,
			DiscordName: username,
			UserType:    0,
			Points:      0,
		}
	} else {
		newUser = &User{
			ID:          id,
			Name:        name,
			DiscordName: username,
			UserType:    1,
			Points:      0,
		}
	}

	_, err := db.NewInsert().
		Model(newUser).
		Exec(ctx)

	if err != nil {
		log.Printf("[ERROR]: could not properly insert new user")
		return -1, exceptions.DatabaseError{
			Operation: "insert",
			Err:       err,
		}
	}

	return id, nil
}
