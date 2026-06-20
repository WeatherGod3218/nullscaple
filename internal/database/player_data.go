package database

import (
	"context"
	"fmt"
	"time"

	t "github.com/WeatherGod3218/nullscaple/internal/nulltypes"
	"github.com/WeatherGod3218/nullscaple/internal/timeutil"
	"github.com/jackc/pgx/v5"
)

func InitPlayerData() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Exec(ctx, `
		DO $$ BEGIN
			CREATE TYPE game_result AS ENUM ('None', 'Win', 'Loss');
		EXCEPTION
			WHEN duplicate_object THEN NULL;
		END $$;
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS player_data (
			player_id 			UUID PRIMARY KEY,

			casual_guesses 		INT NOT NULL DEFAULT 6,
			standard_guesses 	INT NOT NULL DEFAULT 5,
			extreme_guesses 	INT NOT NULL DEFAULT 4,

			casual_result 		game_result NOT NULL DEFAULT 'None',
			standard_result 	game_result NOT NULL DEFAULT 'None',
			extreme_result 		game_result NOT NULL DEFAULT 'None',

			last_played 		VARCHAR(16) NOT NULL DEFAULT '',

			created_at 			TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at			TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	return err
}

func SetPlayerGameResult(id string, result string, difficulty t.GameDifficulties) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var col string

	switch difficulty {
	case t.ModeCasual:
		col = "casual_result"
	case t.ModeStandard:
		col = "standard_result"
	case t.ModeExtreme:
		col = "extreme_result"
	default:
		return fmt.Errorf("invalid mode")
	}

	query := fmt.Sprintf(`
		UPDATE player_data
		SET %s = $1
		WHERE player_id = $2
	`, col)

	_, err := db.Exec(ctx, query, result, id)

	return err
}

func CreatePlayerData(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Exec(ctx, `
		INSERT INTO player_data (player_id, last_played)
		VALUES ($1, $2)
	`, id, timeutil.GetFormattedTime())

	return err
}

func CheckPlayerExists(id string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var exists bool
	err := db.QueryRow(ctx, `
        SELECT EXISTS(
            SELECT 1 FROM player_data
            WHERE player_id = $1
        )
    `, id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func RefreshPlayer(id string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	currentDay := timeutil.GetFormattedTime()
	var lastDay string

	err := db.QueryRow(ctx, `
		SELECT last_played FROM player_data
		WHERE player_id = $1
	`, id).Scan(&lastDay)

	if err != nil {
		return false, err
	}

	if currentDay == lastDay {
		return false, nil
	}

	transaction, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return false, err
	}
	defer transaction.Rollback(ctx)

	_, err = transaction.Exec(ctx, `
		UPDATE player_data
		SET casual_guesses = DEFAULT,
			standard_guesses = DEFAULT,
			extreme_guesses = DEFAULT,

			casual_result = DEFAULT,
			standard_result = DEFAULT,
			extreme_result = DEFAULT,

			last_played = $1
		WHERE player_id = $2
	`, currentDay, id)
	if err != nil {
		return false, err
	}

	_, err = transaction.Exec(ctx, `
		DELETE FROM player_guesses
		WHERE player_id = $1
	`, id)
	if err != nil {
		return false, err
	}

	err = transaction.Commit(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}
