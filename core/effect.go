package core

// Effect of cards/wonders/ptokens
type Effect interface {
	applyEffect(g *Game, i PlayerIndex)
}
