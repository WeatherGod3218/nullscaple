package types

type GameResults string

const (
	None GameResults = "None"
	Win  GameResults = "Win"
	Loss GameResults = "Loss"
)

type GameDifficulties string

const (
	ModeCasual   GameDifficulties = "casual"
	ModeStandard GameDifficulties = "standard"
	ModeExtreme  GameDifficulties = "extreme"
)
