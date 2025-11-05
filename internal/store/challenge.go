package store

import (
	"github.com/uptrace/bun"
	"time"
)

type Challenge struct {
	bun.BaseModel `bun:"table:challenges"`
	ID            int `bun:",pk,autoincrement"`
	Name          string
	Description   string
	StartDate     time.Time `bun:",notnull"`
	EndDate       time.Time `bun:",notnull"`
}
