package core

// Board represents the military rivalry between the two cities.
// It also holds the Military tokens and the Progress tokens available for
// the current game.
type Board struct {
	Military Military

	activePTokens  []PTokenName
	discardPTokens []PTokenName
}
