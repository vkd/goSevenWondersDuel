package core

import (
	"fmt"
)

// VP - victory points
type VP uint8

func (v VP) Mul(i uint8) VP {
	return v * VP(i)
}

type typedVP struct {
	v VP
	t VPType
}

var _ Effect = typedVP{}

func (v typedVP) applyEffect(g *Game, i PlayerIndex) {
	g.vps[i][v.t] += v.v
}

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
	SumVP
	numVPTypes = iota
)

var (
	namesVPType = map[VPType]string{
		BlueVP:     "Blue",
		GreenVP:    "Green",
		YellowVP:   "Yellow",
		PurpleVP:   "Purple",
		WonderVP:   "Wonder",
		PTokenVP:   "PToken",
		CoinsVP:    "Coins",
		MilitaryVP: "Military",
		SumVP:      "Sum",
	}
	_ = [1]struct{}{}[len(namesVPType)-numVPTypes]
)

func (t VPType) String() string {
	return namesVPType[t]
}

// VPTypeByColor of card
func VPTypeByColor(c CardColor) VPType { //nolint: interfacer
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

type maxVPsPerCards struct {
	VP     VP
	Type   VPType
	Colors []CardColor
}

var _ Effect = maxVPsPerCards{}
var _ Finaler = maxVPsPerCards{}

func MaxOneVPPerCards(colors ...CardColor) Effect {
	return maxVPsPerCards{
		VP:     VP(1),
		Colors: colors,
	}
}

func (m maxVPsPerCards) applyEffect(g *Game, i PlayerIndex) {
	g.endEffects[i] = append(g.endEffects[i], m)
}

func (m maxVPsPerCards) finalVP(g *Game, _ PlayerIndex) typedVP {
	var max VP
	for pi := range g.players {
		var vp VP
		for _, color := range m.Colors {
			vp += m.VP.Mul(uint8(g.CardsState.NumByColor(color, PlayerIndex(pi))))
		}
		if vp > max {
			max = vp
		}
	}

	return typedVP{max, m.Type}
}

type vPsPerWonder struct {
	VP   VP
	Type VPType
}

var _ Effect = vPsPerWonder{}
var _ Finaler = vPsPerWonder{}

func (v vPsPerWonder) applyEffect(g *Game, i PlayerIndex) {
	g.endEffects[i] = append(g.endEffects[i], v)
}

func (v vPsPerWonder) finalVP(g *Game, i PlayerIndex) typedVP {
	var max VP
	for pi := range g.players {
		vp := v.VP.Mul(uint8(g.WondersState.CountBuiltByPlayer(PlayerIndex(pi))))
		if vp > max {
			max = vp
		}
	}
	return typedVP{max, v.Type}
}

// BuildersGuild - At the end of the game, this card is worth 2 victory points for each Wonder constructed in the city which has the most wonders.
func BuildersGuild() Effect {
	return vPsPerWonder{VP: 2}
}

// VPPerCoins - worth 1 victory point for each set of n coins in the city.
type vPPerCoins struct {
	Coins Coins
	Type  VPType
}

var _ Effect = vPPerCoins{}
var _ Finaler = vPPerCoins{}

func (v vPPerCoins) applyEffect(g *Game, i PlayerIndex) {
	g.endEffects[i] = append(g.endEffects[i], v)
}

func (v vPPerCoins) finalVP(g *Game, i PlayerIndex) typedVP {
	var max VP
	for pi := range g.players {
		vp := VP(1).Mul(g.players[pi].Coins.Div(3))
		if vp > max {
			max = vp
		}
	}
	return typedVP{max, v.Type}
}

// MoneylendersGuild - At the end of the game, this card is worth 1 victory point for each set of 3 coins in the city.
func MoneylendersGuild() Effect {
	return vPPerCoins{Coins: 3}
}

type Finaler interface {
	finalVP(*Game, PlayerIndex) typedVP
}
