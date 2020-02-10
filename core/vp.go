package core

import "fmt"

// VP - victory points
type VP uint8

var _ Effect = VP(0)

func (vp VP) applyEffect(g *Game, i PlayerIndex) {
	panic("Not implemented")
}

// Apply effect
// func (vp VictoryPoint) Apply(g *Game, i PlayerIndex) {
// 	g.player(i).VP[vp.Type].Value += vp.Value
// }

// VPType - type of victory points
type VPType uint8

// VP for stats
const (
	BlueVP VPType = iota
	GreenVP
	YellowVP
	PurpleVP
	WonderVP
	PTokenVP
	CoinsVP
	MilitaryVP
	numVPTypes = iota
)

// VPTypeByColor of card
func VPTypeByColor(c CardColor) VPType {
	switch c {
	case Blue:
		return BlueVP
	case Green:
		return GreenVP
	case Yellow:
		return YellowVP
	case Purple:
		return PurpleVP
	default:
		panic(fmt.Sprintf("not supported color for VP: %s", c.String()))
	}
}

// VPPerWonder - extra VP per every built wonder
type VPPerWonder uint8

// FinalVP - extra VP at the end of game
func (v VPPerWonder) FinalVP(g *Game, i PlayerIndex) VP {
	return VP(len(g.BuildWonders[i])) * VP(v)
}

// MaxFinalVPOfPlayers - finaler return max value of VP of every players
type MaxFinalVPOfPlayers struct {
	Finaler interface{}
}

var _ Effect = MaxFinalVPOfPlayers{}

func (m MaxFinalVPOfPlayers) applyEffect(g *Game, _ PlayerIndex) {
	panic("Not implemented")
	// var out VP
	// for i := range g.players {
	// 	vp := m.FinalVP(g, PlayerIndex(i))
	// 	if vp > out {
	// 		out = vp
	// 	}
	// }
	// return out
}

// BuildersGuild - At the end of the game, this card is worth 2 victory points for each Wonder constructed in the city which has the most wonders.
func BuildersGuild() Effect {
	return MaxFinalVPOfPlayers{VPPerWonder(2)}
}

// VPPerCoins - worth 1 victory point for each set of 3 coins in the city.
type VPPerCoins uint8

// FinalVP - extra VP at the end of game
func (v VPPerCoins) FinalVP(g *Game, i PlayerIndex) VP {
	return VP(g.player(i).Coins.Div(uint(v)))
}

// MoneylendersGuild - At the end of the game, this card is worth 1 victory point for each set of 3 coins in the city.
func MoneylendersGuild() Effect {
	return MaxFinalVPOfPlayers{VPPerCoins(3)}
}
