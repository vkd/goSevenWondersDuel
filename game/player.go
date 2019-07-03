package game

const (
	numPlayers = 2
)

// Player - game player
type Player struct {
	Money             Money
	Resources         Resources // Brown and Grey cards
	VP                VP        // VP
	ScientificSymbols ScientificSymbols
	ChainSymbols      ChainSymbols

	// // City
	Buildings []*Card
	Wonders   WonderNames

	OnePriceMarkets []OnePriceMarket
	OneOfAnyMarkets []OneOfAnyMarket
}

// PlayerIndex ...
type PlayerIndex int

// Next player index
func (i PlayerIndex) Next() PlayerIndex {
	return (i + 1) % numPlayers
}

// AddScience ...
func AddScience(g *Game, ss ScientificSymbol) {
	player := g.player()
	if player.ScientificSymbols == nil {
		player.ScientificSymbols = make(ScientificSymbols)
	}
	player.ScientificSymbols[ss]++
	if player.ScientificSymbols[ss]%2 == 0 {
		g.canSelectActiveProgressToken()
	}
}
