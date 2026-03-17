package repositories

import (
	"context"
	"time"

	models "github.com/davidwang/factions/internal/domain/models"
)

type ChallengeRepo interface {
	CreateChallenge(ctx context.Context, challenge *models.Challenge) (int, error)
	GetChallengeByID(ctx context.Context, challengeID int) (*models.Challenge, error)
	ListChallengesActiveAt(ctx context.Context, at time.Time) ([]models.Challenge, error)
	AttachSubmissionToChallenge(
		ctx context.Context,
		submissionID int,
		challengeID int,
		score int,
	) error
}
