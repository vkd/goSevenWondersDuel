package core

// Finaler - extra VP at the end of a game
type Finaler interface {
	FinalVP(g *Game, i PlayerIndex) VP
}
