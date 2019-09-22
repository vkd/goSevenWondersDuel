package core

// Effect of cards/wonders/ptokens
type Effect interface {
	Apply(g *Game, i PlayerIndex)
}
