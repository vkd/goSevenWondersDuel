package core

// PlayerIndex - index of player
type PlayerIndex int

// Next player index
func (i PlayerIndex) Next() PlayerIndex {
	return (i + 1) % numPlayers
}

// Player of a game
type Player struct {
	Coins     Coins
	Resources Resources

	Chains Chains

	ScientificSymbols ScientificSymbols

	AvailableWorneds []WonderName
	BuildWonders     []WonderName

	PriceMarkets PriceMarkets

	BuiltCards [numOfCardColors][]CardName

	VP VP
}

// NewPlayer of a game
func NewPlayer() Player {
	p := Player{
		Coins: 7,
	}
	return p
}
