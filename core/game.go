package core

const (
	numPlayers     = 2
)

// Game - game state
type Game struct {
	players       [numPlayers]Player
	currentPlayer PlayerIndex

	military Military

}

// Cost card by coins
func (g *Game) Cost(card CardName) Coins {
	return card.card().Cost.ByCoins(g, g.currentPlayer)
}

func (g *Game) player(i PlayerIndex) *Player {
	return &g.players[i]
}

func (g *Game) apply(card CardName) {
	for _, e := range card.card().Effects {
		e.Apply(g, g.currentPlayer)
	}
}
