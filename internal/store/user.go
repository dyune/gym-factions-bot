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
	UserType       UserType `bun:",default:1"`
	Points         int      `bun:",notnull,default:0"`
	FactionID      int
	AdminFactionID int `bun:"admin_faction_id"`

	Faction *Faction `bun:"rel:belongs-to,join:faction_id=id"`
}

func ExistsByID(id int, ctx context.Context, db *bun.DB) (bool, error) {
	exists, err := db.NewSelect().
		Model((*User)(nil)).
		Where("id = ?", id).
		Exists(ctx)

	if err != nil {
		log.Printf("[ERROR] while calling ExistsByID: %v", err)
		return false, err
	}
	if exists {
		return true, err
	}

	return false, err
}

func InsertUser(id int, name string, ctx context.Context, db *bun.DB) (int, error) {
	if exists, _ := ExistsByID(id, ctx, db); exists {
		return -1, exceptions.UserExistsError{
			UserID:   id,
			Username: name,
		}
	}

	var newUser *User
	if AdminList[name] {
		newUser = &User{
			ID:       id,
			Name:     name,
			UserType: 0,
			Points:   0,
		}
	} else {
		newUser = &User{
			ID:       id,
			Name:     name,
			UserType: 1,
			Points:   0,
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

func LookUpUserByID(id int, ctx context.Context, db *bun.DB) (*User, error) {
	user := new(User)
	err := db.NewSelect().
		Model(user).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return user, nil

}

func SetUserDetailsByID(
	id int,
	name *string,
	points *int,
	factionID *int,
	ctx context.Context,
	db *bun.DB,
) error {

	hasUpdates := false
	update := db.NewUpdate().
		Model((*User)(nil)).
		Where("id = ?", id)

	// Explicit deref since Go won't auto-deref here
	if name != nil {
		update = update.Set("name = ?", *name)
		hasUpdates = true
	}

	if points != nil {
		update = update.Set("points = ?", *points)
		hasUpdates = true
	}

	if factionID != nil {
		update = update.Set("faction_id = ?", *factionID)
		hasUpdates = true
	}

	if hasUpdates {
		_, err := update.Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
