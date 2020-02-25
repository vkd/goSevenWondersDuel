package core

import (
	"fmt"
	"math/rand"
)

const (
	// PTokensCount - amount of different PTokens of the game.
	PTokensCount = 10
)

// PToken - Progress Tokens.
// The progress tokens represent effects which you can obtain by collecting
// identical pairs of scientific symbols.
type PToken struct {
	Name   PTokenName
	Effect Effect
}

// PTokenName - name of a progress token.
type PTokenName string

// PTokenID - ID of a PToken.
type PTokenID uint8

func (p PTokenID) pToken() *PToken {
	return &allPTokens[p]
}

func pTokenID(name PTokenName) PTokenID {
	for i, p := range allPTokens {
		if p.Name == name {
			return PTokenID(i)
		}
	}
	panic(fmt.Sprintf("cannot find %q ptoken", name))
}

var (
	pTokenIDs [PTokensCount]PTokenID
)

func init() {
	for i := range pTokenIDs {
		pTokenIDs[i] = PTokenID(i)
	}
}

func shufflePTokens(rnd *rand.Rand) []PTokenID {
	var ptokens = pTokenIDs
	rnd.Shuffle(len(ptokens), func(i, j int) {
		ptokens[i], ptokens[j] = ptokens[j], ptokens[i]
	})
	return ptokens[:]
}

var _ = [1]struct{}{}[len(shufflePTokens(zeroRand()))-PTokensCount]

var (
	allPTokens = []PToken{
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
	_ = [1]struct{}{}[len(allPTokens)-PTokensCount]
)

func newPToken(name PTokenName, ee ...interface{}) PToken {
	var pt = PToken{
		Name: name,
	}

	var es Effects
	for _, e := range ee {
		switch e := e.(type) {
		case VP:
			es = append(es, typedVP{e, PTokenVP})
		case Effect:
			es = append(es, e)
		default:
			panic(fmt.Sprintf("Not allowed for the PToken constructor: %T", e))
		}
	}
	pt.Effect = es
	return pt
}

// --- TODO ---

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
		v: VP(3).Mul(uint8(len(g.builtPTokens[i]))),
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
