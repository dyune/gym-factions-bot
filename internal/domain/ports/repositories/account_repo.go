package repositories

import (
	"context"

	models "github.com/davidwang/factions/internal/domain/models"
)

type AccountRepo interface {
	AccountExistsByID(ctx context.Context, id int) (bool, error)
	InsertAccount(ctx context.Context, id int, name string, factionID int) (int, error)
	GetAccountByID(ctx context.Context, id int) (*models.Account, error)
	GetAccountFactionByID(ctx context.Context, id int) (*models.Faction, error)
	SetAccountDetailsByID(
		ctx context.Context,
		id int,
		name *string,
		points *int,
		faction *models.Faction,
	) error
	GetAccountByDiscordUserID(ctx context.Context, discordUserID string) (*models.Account, error)
	ListAccountsByFactionID(ctx context.Context, factionID int) ([]models.Account, error)
	AssignFaction(ctx context.Context, accountID int, factionID int) error
}
