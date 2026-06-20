package types

type PlayerData struct {
	PlayerId string

	CasualGuesses   int
	StandardGuesses int
	ExtremeGuesses  int

	CasualResult   string
	StandardResult string
	ExtremeResult  string

	LastPlayed string
}

type PlayerGuesses struct {
	Id       string
	PlayerId string

	Difficulty string
	EnemyId    string
	Day        string
}
