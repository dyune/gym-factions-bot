package store

import (
	"context"
	"github.com/davidwang/factions/internal/exceptions"
	"github.com/uptrace/bun"
	"log"
	"time"
)

type SubmissionType int

const (
	NonGym SubmissionType = iota
	Gym
)

// Change these to actual SQL-compatible timestamps
const weekStart = 0
const weekEnd = 1

const startDate = 0
const endDate = 0

type Week struct {
	Num   int
	Start time.Time
	End   time.Time
}

// Each week number corresponds to a start-date and end-date pair
var weeks map[int]Week

type Submission struct {
	bun.BaseModel `bun:"table:submissions"`

	ID             int       `bun:",pk,autoincrement"`
	SubmittedAt    time.Time `bun:",default:current_timestamp"`
	SubmissionType int
	OwnerID        int `bun:",notnull"`
	Points         int `bun:",nullzero"`

	Owner *Account `bun:"rel:belongs-to,join:owner_id=id"`
}

type ChallengeSubmission struct {
	bun.BaseModel `bun:"table:challenge_submissions"`

	SubmissionID int `bun:",pk"` // Reference to Submission
	Score        int `bun:",nullzero"`
	ChallengeID  int `bun:",notnull"`

	Submission *Submission `bun:"rel:belongs-to,join:submission_id=id"`
	Challenge  *Challenge  `bun:"rel:belongs-to,join:challenge_id=id"`
}

func GetSubsByWeekAndOwnerID() {

}

func SubExistsByDayAndOwnerID(
	time time.Time,
	id int,
	ctx context.Context,
	db *bun.DB,
) (bool, error) {

	exists, err := db.NewSelect().
		Model((*Submission)(nil)).
		Where("owner_id = ?", id).
		Where("DATE(submitted_at) = DATE(?)", time).
		Exists(ctx)

	if err != nil {
		return false, exceptions.DatabaseError{
			Operation: "subExistsByDayAndOwnerID",
			Err:       err,
		}
	}
	return exists, nil
}

func InsertSub(
	time time.Time,
	subType SubmissionType,
	id int,
	pts int,
	ctx context.Context,
	db *bun.DB,
) error {

	sub := &Submission{
		OwnerID:        id,
		Points:         pts,
		SubmissionType: int(subType),
		SubmittedAt:    time,
	}
	_, err := db.NewInsert().Model(sub).Exec(ctx)
	if err != nil {
		log.Printf("[ERROR]: could not properly insert new submission")
		return exceptions.DatabaseError{
			Operation: "insertSub",
			Err:       err,
		}
	}
	return nil
}
