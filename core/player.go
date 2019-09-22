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

	AvailableWorneds []WonderName
	BuildWonders     []WonderName

	PriceMarkets PriceMarkets
}
