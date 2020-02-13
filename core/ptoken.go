package core

import (
	"fmt"
	"math/rand"
)

// PToken - progress token
type PToken struct {
	Name    PTokenName
	Effects []Effect
}

type PTokenID uint8

// PTokenName - name of a progress token
type PTokenName string

const (
	numPTokens = 10
)

var (
	listPTokens = []PToken{
		newPToken("Agriculture", Coins(6), VP(4)),
		newPToken("Architecture", Architecture()),
		newPToken("Economy", Economy()),
		newPToken("Law", Scales),
		newPToken("Masonry", Masonry()),
		newPToken("Mathematics", Mathematics()),
		newPToken("Philosophy", VP(7)),
		newPToken("Strategy", Strategy()),
		newPToken("Theology", Theology()),
		newPToken("Urbanism", Coins(6), Urbanism()),
	}
	_ = [1]struct{}{}[len(listPTokens)-numPTokens]

	mapPTokens = makeMapPTokensByName(listPTokens)
	_          = [1]struct{}{}[len(mapPTokens)-numPTokens]

	listPTokensIDs [numPTokens]PTokenID
)

func init() {
	for i := range listPTokensIDs {
		listPTokensIDs[i] = PTokenID(i)
	}
}

func newPToken(name PTokenName, ee ...interface{}) PToken {
	var pt = PToken{
		Name: name,
	}
	for _, e := range ee {
		switch e := e.(type) {
		case VP:
			pt.Effects = append(pt.Effects, typedVP{e, PTokenVP})
		case Effect:
			pt.Effects = append(pt.Effects, e)
		default:
			panic(fmt.Sprintf("Not allow for ptoken builder: %T", e))
		}
	}
	return pt
}

func makeMapPTokensByName(list []PToken) map[PTokenName]*PToken {
	m := map[PTokenName]*PToken{}
	for i, t := range list {
		m[t.Name] = &list[i]
	}
	return m
}

func shufflePTokens(rnd *rand.Rand) []PTokenID {
	var ptokens = listPTokensIDs
	rnd.Shuffle(len(ptokens), func(i, j int) {
		ptokens[i], ptokens[j] = ptokens[j], ptokens[i]
	})
	return ptokens[:]
}

func Architecture() Effect {
	return effertFunc(func(g *Game, i PlayerIndex) {
		g.player(i).IsArchitecture = true
	})
}

func Economy() Effect {
	return effertFunc(func(g *Game, i PlayerIndex) {
		g.player(i).IsEconomy = true
	})
}

func Masonry() Effect {
	return effertFunc(func(g *Game, i PlayerIndex) {
		g.player(i).IsMasonry = true
	})
}

func Mathematics() Effect {
	return mathematics{}
}

type mathematics struct{}

var _ Effect = mathematics{}
var _ Finaler = mathematics{}

func (m mathematics) applyEffect(g *Game, i PlayerIndex) {
	g.endEffects[i] = append(g.endEffects[i], m)
}

func (mathematics) finalVP(g *Game, i PlayerIndex) typedVP {
	return typedVP{
		v: VP(3).Mul(len(g.builtPTokens[i])),
		t: PTokenVP,
	}
}

func Strategy() Effect {
	return effertFunc(func(g *Game, i PlayerIndex) {
		g.player(i).IsStrategy = true
	})
}

func Theology() Effect {
	return effertFunc(func(g *Game, i PlayerIndex) {
		g.player(i).IsTheology = true
	})
}

func Urbanism() Effect {
	return effertFunc(func(g *Game, i PlayerIndex) {
		g.player(i).IsUrbanism = true
	})
}

type effertFunc func(*Game, PlayerIndex)

func (fn effertFunc) applyEffect(g *Game, i PlayerIndex) {
	fn(g, i)
}
