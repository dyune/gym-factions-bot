package store

import (
	"context"
	"fmt"
	"github.com/davidwang/factions/internal/exceptions"
	"github.com/uptrace/bun"
	_ "github.com/uptrace/bun"
	"log"
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
	bun.BaseModel `bun:"table:accounts,alias:u"`

	ID             int `bun:",pk"`
	Name           string
	AccountType    AccountType `bun:",default:1"`
	Points         int         `bun:",notnull,default:0"`
	FactionID      int         `bun:"faction_id"`
	AdminFactionID int         `bun:"admin_faction_id"`

	Faction *Faction `bun:"rel:belongs-to,join:faction_id=id"`
}

func AccountExistsByID(id int, ctx context.Context, db *bun.DB) (bool, error) {
	exists, err := db.NewSelect().
		Model((*Account)(nil)).
		Where("id = ?", id).
		Exists(ctx)

	if err != nil {
		log.Printf("[ERROR] while calling AccountExistsByID: %v", err)
		return false, exceptions.DatabaseError{
			Operation: "exists",
			Err:       err,
		}
	}
	if exists {
		return true, nil
	}

	return false, nil
}

func InsertAccount(id int, name string, factionID int, ctx context.Context, db *bun.DB) (int, error) {

	if exists, _ := AccountExistsByID(id, ctx, db); exists {
		return -1, exceptions.AccountExistsError{
			AccountID:   id,
			AccountName: name,
		}
	}

	newAccount := &Account{
		ID:          id,
		Name:        name,
		AccountType: 1,
		Points:      0,
		FactionID:   factionID,
	}

	_, err := db.NewInsert().
		Model(newAccount).
		Exec(ctx)

	if err != nil {
		log.Printf("[ERROR]: could not properly insert new Account")
		return -1, exceptions.DatabaseError{
			Operation: "insert",
			Err:       err,
		}
	}

	return id, nil
}

func GetAccountByID(id int, ctx context.Context, db *bun.DB) (*Account, error) {
	account := new(Account)
	err := db.NewSelect().
		Model(account).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		log.Printf("[ERROR]: could not properly insert new Account")
		return nil, exceptions.DatabaseError{
			Operation: "getAccountByID",
			Err:       err,
		}
	}

	return account, nil

}

func GetAccountFactionByID(id int, ctx context.Context, db *bun.DB) (*Faction, error) {
	account := new(Account)

	err := db.NewSelect().
		Model(account).
		Relation("Faction").
		Where("u.id = ?", id). // Account table aliased as 'u' to prevent keyword conflict
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	if account.Faction == nil {
		return nil, fmt.Errorf("account %d has no faction", id)
	}

	return account.Faction, nil
}

func SetAccountDetailsByID(
	id int,
	name *string,
	points *int,
	faction *Faction,
	ctx context.Context,
	db *bun.DB,
) error {

	hasUpdates := false
	update := db.NewUpdate().
		Model((*Account)(nil)).
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

	if faction != nil {
		update = update.Set("faction = ?", *faction)
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
