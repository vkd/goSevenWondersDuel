package core

// Coins allow you to construct certain Buildings, and to purchase resources
// through commerce. The Treasury, the accumulated coins, is worth victory
// points at the end of the game.
type Coins uint

func (c Coins) applyPrice(p *Price) {
	p.Coins += c
}

// Mul - multiply
func (c Coins) Mul(i uint) Coins {
	return Coins(uint(c) * i)
}

// Div - division
func (c Coins) Div(i uint) uint {
	return uint(c) / i
}

// Cost based on coins
// func (c Coins) Cost() Cost {
// 	return Cost{Coins: c}
// }

// Apply effect
func (c Coins) Apply(g *Game, i PlayerIndex) {
	g.players[i].Coins += c
}

// CoinsPerWonder - The card is worth x coins per Wonder constructed in your city at the time it is constructed.
type CoinsPerWonder Coins

// Apply effect
func (c CoinsPerWonder) Apply(g *Game, i PlayerIndex) {
	worth := Coins(c).Mul(uint(len(g.player(i).BuildWonders)))
	g.player(i).Coins += worth
}

// CoinsPerCardColor - This card is worth x coins for each one color card constructed in the player’s city at the moment when it is constructed.
type CoinsPerCardColor struct {
	Colors []CardColor
	Coins  Coins
}

// Apply effect
func (c CoinsPerCardColor) Apply(g *Game, i PlayerIndex) {
	for _, color := range c.Colors {
		g.player(i).Coins += coinsPerColor(c.Coins, color, *g.player(i))
	}
}

func coinsPerColor(coins Coins, color CardColor, p Player) Coins {
	return coins.Mul(uint(len(p.BuiltCards[color])))
}

// MaxCoinsPerCardColor - At the time it is constructed, this card grants you 1 coin for each color card in the city which has the most there colors cards.
type MaxCoinsPerCardColor CoinsPerCardColor

// Apply effect
func (m MaxCoinsPerCardColor) Apply(g *Game, i PlayerIndex) {
	var max Coins
	for _, p := range g.players {
		var c Coins
		for _, color := range m.Colors {
			c += coinsPerColor(m.Coins, color, p)
		}
		if c > max {
			max = c
		}
	}

	g.player(i).Coins += max
}
