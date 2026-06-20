package database

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func InitDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newDB, err := pgxpool.New(
		ctx,
		os.Getenv("POSTGRESQL_ADDRESS"),
	)

	if err != nil {
		return err
	}

	err = newDB.Ping(ctx)
	if err != nil {
		return err
	}
	db = newDB

	InitPlayerData()
	InitPlayerGuesses()
	return nil
}
