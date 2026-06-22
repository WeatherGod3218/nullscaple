package database

import (
	"context"
	"fmt"
	"time"

	"github.com/WeatherGod3218/nullscaple/internal/idgen"
	t "github.com/WeatherGod3218/nullscaple/internal/nulltypes"
	"github.com/WeatherGod3218/nullscaple/internal/timeutil"
	"github.com/jackc/pgx/v5"
)

func AddPlayerGuess(playerId string, enemyId string, difficulty t.GameDifficulties) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rowId, err := idgen.GenerateNewId()
	if err != nil {
		return 0, err
	}

	var col string
	switch difficulty {
	case "casual":
		col = "casual_guesses"
	case "standard":
		col = "standard_guesses"
	case "extreme":
		col = "extreme_guesses"
	default:
		return 0, fmt.Errorf("invalid mode")
	}

	query := fmt.Sprintf(`
		UPDATE player_data
		SET %s = %s - 1
		WHERE player_id = $1
		RETURNING %s
	`, col, col, col)

	transaction, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, err
	}

	defer transaction.Rollback(ctx)

	var remaining int
	err = transaction.QueryRow(ctx, query, playerId).Scan(&remaining)
	if err != nil {
		return 0, err
	}
	if remaining < 0 {
		return remaining, fmt.Errorf("no guesses remaining")
	}

	_, err = transaction.Exec(ctx, `
		INSERT INTO player_guesses (id, player_id, difficulty, enemy_id, day)
		VALUES ($1, $2, $3, $4, $5)
	`, rowId, playerId, difficulty, enemyId, timeutil.GetFormattedTime())
	if err != nil {
		return 0, err
	}

	err = transaction.Commit(ctx)
	if err != nil {
		return 0, err
	}

	return remaining, nil
}

func GetPlayerGuesses(playerId string, difficulty t.GameDifficulties) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch difficulty {
	case "casual", "standard", "extreme":
	default:
		return nil, fmt.Errorf("invalid mode: %s", difficulty)
	}

	rows, err := db.Query(ctx, `
		SELECT enemy_id FROM player_guesses
		WHERE player_id = $1 AND difficulty = $2
		ORDER BY guessed_at ASC
	`, playerId, difficulty)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enemyIds []string
	for rows.Next() {
		var enemyId string

		err := rows.Scan(&enemyId)
		if err != nil {
			return nil, err
		}

		enemyIds = append(enemyIds, enemyId)
	}

	return enemyIds, nil
}

func InitPlayerGuesses() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS player_guesses (
			id 			UUID PRIMARY KEY,
			player_id 	UUID NOT NULL REFERENCES player_data(player_id),

			difficulty 	VARCHAR(32) NOT NULL,
			enemy_id	INT NOT NULL,
			day 		VARCHAR(16) NOT NULL DEFAULT '',
			guessed_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, `
		CREATE INDEX IF NOT EXISTS idx_player_guesses_lookup 
			ON player_guesses(player_id, mode, day);
	`)

	return err
}
