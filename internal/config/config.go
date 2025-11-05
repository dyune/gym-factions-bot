package config

import (
	"context"
	"github.com/uptrace/bun"
)

var (
	GlobalCtx context.Context
	DB        *bun.DB
)
