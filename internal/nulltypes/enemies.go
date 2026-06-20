package types

import "fmt"

type ComparisonResults string

const (
	Comparison_Higher      ComparisonResults = "higher"
	Comparison_High_Middle ComparisonResults = "high-mid"
	Comparison_Equal       ComparisonResults = "equal"
	Comparison_Low_Middle  ComparisonResults = "low-mid"
	Comparison_Lower       ComparisonResults = "lower"
)

type EnemyData struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	CurseAmount int    `json:"curse_amount"`
	StartLevel  int    `json:"start_level"`
	KillMethod  string `json:"kill_method"`
	Filename    string `json:"filename"`
}

type EnemyComparison struct {
	Id          bool              `json:"id"`
	Name        bool              `json:"name"`
	KillMethod  bool              `json:"kill_method"`
	CurseAmount ComparisonResults `json:"curse_amount"`
	StartLevel  ComparisonResults `json:"start_level"`
}

type EnemyRequest struct {
	ID   string `json:"enemy_id"`
	Mode string `json:"gameplay_mode"`
}

type EnemyGuess struct {
	Enemy            *EnemyData      `json:"enemy"`
	ComparisonResult EnemyComparison `json:"result"`
}

func ParseDifficulty(s string) (GameDifficulties, error) {
	switch GameDifficulties(s) {
	case ModeCasual, ModeStandard, ModeExtreme:
		return GameDifficulties(s), nil
	}
	return "", fmt.Errorf("invalid difficulty: %s", s)
}
