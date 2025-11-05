package store

import (
	"context"
	"github.com/uptrace/bun"
	"log"
)

func InitTables(db *bun.DB, ctx context.Context) {
	// Create table
	_, err := db.
		NewCreateTable().
		Model((*User)(nil)). // (*User)(nil) is a null pointer to User
		IfNotExists().
		Exec(ctx)

	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
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
}
