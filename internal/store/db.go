package store

import (
	"context"
	"errors"
	"github.com/davidwang/factions/internal/config"
	"github.com/davidwang/factions/internal/exceptions"
	"github.com/uptrace/bun"
	"log"
)

func InitTables(db *bun.DB, ctx context.Context) {
	// Create table
	_, err := db.
		NewCreateTable().
		Model((*Account)(nil)). // (*Account)(nil) is a null pointer to Account
		IfNotExists().
		Exec(ctx)

	if err != nil {
		log.Fatalf("Failed to create Accounts table: %v", err)
	}

	_, err = db.
		NewCreateTable().
		Model((*Submission)(nil)).
		IfNotExists().
		Exec(ctx)

	if err != nil {
		log.Fatalf("Failed to create submissions table: %v", err)
	}

	_, err = db.
		NewCreateTable().
		Model((*Faction)(nil)).
		IfNotExists().
		Exec(ctx)

	if err != nil {
		log.Fatalf("Failed to create factions table: %v", err)
	}

	_, err = db.
		NewCreateTable().
		Model((*ChallengeSubmission)(nil)).
		IfNotExists().
		Exec(ctx)

	if err != nil {
		log.Fatalf("Failed to create challenge submissions table: %v", err)
	}

	log.Printf("All tables created successfully.")

	for _, v := range FactionNames {
		err = InsertFaction(config.GlobalCtx, config.DB, v, "")
		if errors.Is(err, exceptions.ErrInvalidInput) {
			log.Printf("Skipped duplicate faction '%s'", v)

		} else if err != nil {
			log.Fatalf("Failed to create faction '%s'", v)
		}
	}

	log.Printf("Factions initialized.")

}
