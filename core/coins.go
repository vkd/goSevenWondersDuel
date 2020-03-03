package core

// Coins allow you to construct certain Buildings, and to purchase resources
// through commerce. The Treasury, the accumulated coins, is worth victory
// points at the end of the game.
type Coins uint8

var _ Pricer = Coins(0)
var _ Effect = Coins(0)

func (c Coins) applyPrice(p *Price) {
	p.Coins += c
}

func (c Coins) applyEffect(g *Game, i PlayerIndex) {
	g.player(i).Coins += c
}

func (c *Coins) sub(v Coins) {
	if *c > v {
		*c -= v
	} else {
		*c = 0
	}
}

// Mul - multiply.
func (c Coins) Mul(i uint8) Coins {
	return Coins(uint8(c) * i)
}

// Div - division.
func (c Coins) Div(i uint8) uint8 {
	return uint8(c) / i
}

// --- TODO ---

// CoinsPerWonder - The card is worth x coins per Wonder constructed in your city at the time it is constructed.
type CoinsPerWonder Coins

var _ Effect = CoinsPerWonder(0)

func (c CoinsPerWonder) applyEffect(g *Game, i PlayerIndex) {
	worth := Coins(c).Mul(uint8(len(g.buildWonders[i])))
	g.player(i).Coins += worth
}

// CoinsPerCardColor - This card is worth x coins for each one color card constructed in the player’s city at the moment when it is constructed.
type CoinsPerCardColor struct {
	Colors []CardColor
	Coins  Coins
}

func CoinsPerCard(color CardColor, c Coins) CoinsPerCardColor {
	return CoinsPerCardColor{
		Colors: []CardColor{color},
		Coins:  c,
	}
}

var _ Effect = CoinsPerCardColor{}

func (c CoinsPerCardColor) applyEffect(g *Game, i PlayerIndex) {
	for _, color := range c.Colors {
		g.player(i).Coins += coinsPerColor(c.Coins, color, g, i)
	}
}

func coinsPerColor(coins Coins, color CardColor, g *Game, i PlayerIndex) Coins {
	return coins.Mul(uint8(len(g.builtCards[i][color])))
}

// MaxCoinsPerCardColor - At the time it is constructed, this card grants you 1 coin for each color card in the city which has the most there colors cards.
type MaxCoinsPerCardColor CoinsPerCardColor

var _ Effect = MaxCoinsPerCardColor{}

func (m MaxCoinsPerCardColor) applyEffect(g *Game, i PlayerIndex) {
	var max Coins
	for pi := range g.players {
		var c Coins
		for _, color := range m.Colors {
			c += coinsPerColor(m.Coins, color, g, PlayerIndex(pi))
		}
		if c > max {
			max = c
		}
	}

	g.player(i).Coins += max
}

func MaxOneCoinPerCards(colors ...CardColor) MaxCoinsPerCardColor {
	return MaxCoinsPerCardColor{
		Colors: colors,
		Coins:  1,
	}
}
