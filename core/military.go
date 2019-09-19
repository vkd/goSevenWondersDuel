package core

// Military board
type Military struct {
	// Conflict pawn
	shields [numPlayers]Shields

	// Military tokens
	tokens2 [numPlayers]bool
	tokens5 [numPlayers]bool
}

// Shields - military power
type Shields uint8
