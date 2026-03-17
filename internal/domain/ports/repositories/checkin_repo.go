package repositories

import (
	"context"
	"time"
)

type CheckinRepo interface {
	InsertSub(
		ctx context.Context,
		submittedAt time.Time,
		submissionType int,
		ownerID int,
		points int,
	) error
	SubExistsByWeekAndOwnerID(ctx context.Context, submittedAt time.Time, ownerID int) (bool, error)
	SubExistsByDayAndOwnerID(ctx context.Context, submittedAt time.Time, ownerID int) (bool, error)

	// These methods round out the workout log use cases in AGENTS.md.
	GetSubmissionByID(ctx context.Context, submissionID int) (*SubmissionRecord, error)
	ListSubmissionsByOwnerID(ctx context.Context, ownerID int, limit int) ([]SubmissionRecord, error)
	ListSubmissionsBySeason(ctx context.Context, seasonID int, limit int) ([]SubmissionRecord, error)
	UpdateSubmission(ctx context.Context, submissionID int, update SubmissionUpdate) error
	DeleteSubmission(ctx context.Context, submissionID int) error
}
