package repositories

import (
	"context"

	models "github.com/davidwang/factions/internal/domain/models"
)

type FactionRepo interface {
	ExistsByName(ctx context.Context, name string) (bool, error)
	InsertFaction(ctx context.Context, name string, description string) error
	LookUpFaction(ctx context.Context, name string) (*models.Faction, error)

	// Additional faction operations needed for assignment, management, and
	// leaderboard-style reads.
	GetFactionByID(ctx context.Context, id int) (*models.Faction, error)
	ListActiveFactions(ctx context.Context) ([]models.Faction, error)
	UpdateFaction(ctx context.Context, id int, update FactionUpdate) error
	ArchiveFaction(ctx context.Context, id int) error
}
