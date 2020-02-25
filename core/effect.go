package core

// Effect of cards/wonders/ptokens
type Effect interface {
	applyEffect(g *Game, i PlayerIndex)
}

type Effects []Effect

var _ Effect = Effects(nil)

func (es Effects) applyEffect(g *Game, i PlayerIndex) {
	for _, e := range es {
		e.applyEffect(g, i)
	}
}
