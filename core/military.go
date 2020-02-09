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

var _ Effect = Shields(0)

func (s Shields) applyEffect(g *Game, i PlayerIndex) {
	g.military.shields[i] += s
}
