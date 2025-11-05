package store

import (
	"github.com/uptrace/bun"
	"time"
)

type SubmissionType int

const (
	NonGym SubmissionType = iota
	Gym
)

type Submission struct {
	bun.BaseModel `bun:"table:submissions"`

	ID             int       `bun:",pk,autoincrement"`
	SubmittedAt    time.Time `bun:",default:current_timestamp"`
	SubmissionType int
	OwnerID        int `bun:",notnull"`
	Points         int `bun:",nullzero"`

	Owner *User `bun:"rel:belongs-to,join:owner_id=id"`
}

type ChallengeSubmission struct {
	bun.BaseModel `bun:"table:challenge_submissions"`

	SubmissionID int `bun:",pk"` // Reference to Submission
	Score        int `bun:",nullzero"`
	ChallengeID  int `bun:",notnull"`

	Submission *Submission `bun:"rel:belongs-to,join:submission_id=id"`
	Challenge  *Challenge  `bun:"rel:belongs-to,join:challenge_id=id"`
}
