package postgres

import (
	"time"
)

type Challenge struct {
	ID          int
	Name        string
	Description string
	StartDate   time.Time
	EndDate     time.Time
}
