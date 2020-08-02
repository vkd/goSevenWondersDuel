package core

import (
	"fmt"
)

const (
	// WondersCount - amount of different wonders on the game.
	WondersCount = 12
)

// Wonder card.
// Each large card represents a Wonder from the Age of Antiquity.
// Each Wonder consists of a name, a construction cost, and an effect.
type Wonder struct {
	Name   WonderName
	Cost   Cost
	Effect Effect
}

// WonderName - name of a wonder.
type WonderName string

// WonderID - ID of a wonder.
type WonderID uint32

func (id WonderID) wonder() *Wonder {
	return &allWonders[id]
}

func wonderID(name WonderName) WonderID {
	for i := range allWonders {
		if allWonders[i].Name == name {
			return WonderID(i)
		}
	}
	panic(fmt.Sprintf("cannot find %q wonder", name))
}

var (
	allWonders = []Wonder{
		newWonder("Temple of Artemis", Coins(12), RepeatTurn(), NewCost(Wood, Stone, Glass, Papyrus)),
		newWonder("The Great Lighthouse", OneRawMarket(), VP(4), NewCost(Wood, Stone, Papyrus, Papyrus)),
		newWonder("The Colossus", Shields(2), VP(3), NewCost(Clay, Clay, Clay, Glass)),
		newWonder("The Pyramids", VP(9), NewCost(Stone, Stone, Stone, Papyrus)),
		newWonder("The Mausoleum", BuildFreeDiscardedCard(), VP(2), NewCost(Clay, Clay, Glass, Glass, Papyrus)),
		newWonder("The Statue of Zeus", DiscardOpponentBuild(Brown), Shields(1), VP(3), NewCost(Stone, Wood, Clay, Papyrus, Papyrus)),
		newWonder("The Appian Way", Coins(3), DiscardOpponentCoins(3), RepeatTurn(), VP(3), NewCost(Stone, Stone, Clay, Clay, Papyrus)),
		newWonder("Circus Maximus", DiscardOpponentBuild(Grey), Shields(1), VP(3), NewCost(Stone, Stone, Wood, Glass)),
		newWonder("The Great Library", PlayOneOf3DiscardedPToken(), VP(4), NewCost(Wood, Wood, Wood, Glass, Papyrus)),
		newWonder("Piraeus", OneManufacturedMarket(), RepeatTurn(), VP(2), NewCost(Wood, Wood, Stone, Clay)),
		newWonder("The Hanging Gardens", Coins(6), RepeatTurn(), VP(3), NewCost(Wood, Wood, Glass, Papyrus)),
		newWonder("The Sphinx", RepeatTurn(), VP(6), NewCost(Stone, Clay, Glass, Glass)),
	}
	_ = [1]struct{}{}[WondersCount-len(allWonders)]
)

func newWonder(name WonderName, args ...interface{}) (w Wonder) {
	w.Name = name

	var es Effects
	for i := range args {
		switch a := args[i].(type) {
		case Cost:
			w.Cost = a
		case VP:
			es = append(es, typedVP{a, WonderVP})
		case Effect:
			es = append(es, a)
		default:
			panic(fmt.Sprintf("Not allowed for the Wonder constructor: %T", a))
		}
	}
	w.Effect = es
	return w
}
