package main

import (
	"context"
	"database/sql"
	"github.com/davidwang/factions/bot"
	"github.com/davidwang/factions/internal/config"
	"github.com/davidwang/factions/internal/store"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"log"
	"os"
	"time"
)

func getenv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {

	godotenv.Load()

	dsn := getenv("DATABASE_URL", "postgres://myuser:secret@localhost:5433/mydatabase?sslmode=disable")
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(dsn),
	))

	sqldb.SetMaxOpenConns(25)                 // Maximum open connections
	sqldb.SetMaxIdleConns(10)                 // Maximum idle connections
	sqldb.SetConnMaxLifetime(5 * time.Minute) // Connection lifetime
	sqldb.SetConnMaxIdleTime(5 * time.Minute) // Idle connection timeout

	if err := sqldb.Ping(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	config.DB = bun.NewDB(sqldb, pgdialect.New())

	defer config.DB.Close()

	// Quick health check with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := config.DB.PingContext(ctx); err != nil {
		log.Fatalf("ping db test failure: %v", err)
	}
	log.Printf("Connected to database.")

	// Debug
	config.DB.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
	))

	// Init tables
	config.GlobalCtx = context.Background()
	store.InitTables(config.DB, config.GlobalCtx)

	log.Printf("Running with Bun debug ON.")

	config.GlobalCtx = context.Background()

	bot.Run()

}
