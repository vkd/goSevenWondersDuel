package core

// Coins allow you to construct certain Buildings, and to purchase resources
// through commerce. The Treasury, the accumulated coins, is worth victory
// points at the end of the game.
type Coins uint

// Mul - multiply
func (c Coins) Mul(i uint) Coins {
	return Coins(uint(c) * i)
}

// CoinsPerWonder - The card is worth x coins per Wonder constructed in your city at the time it is constructed.
type CoinsPerWonder Coins

// Apply effect
func (c CoinsPerWonder) Apply(g *Game, i PlayerIndex) {
	worth := Coins(c).Mul(uint(len(g.player(i).BuildWonders)))
	g.player(i).Coins += worth
}
