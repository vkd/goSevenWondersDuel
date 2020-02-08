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

	AvailableWonders []WonderID
	BuildWonders     []WonderID

	PriceMarkets PriceMarkets

	BuiltCards [numCardColors][]CardName

	VP VP
}

// NewPlayer of a game
func NewPlayer() Player {
	p := Player{
		Coins: 7,
	}
	return p
}
