package core

import (
	"fmt"
	"math/rand"
)

// Wonder from the Age of Antiquity
type Wonder struct {
	Name WonderName

	// ----
	Cost    Cost
	Effects []Effect
}

// WonderID - ID of the wonder
type WonderID uint32

func (id WonderID) wonder() *Wonder {
	return &listWonders[id]
}

// WonderName - name of a wonder
type WonderName string

// ----

const (
	numWonders = 12
)

var (
	wonderIDs [numWonders]WonderID
)

func init() {
	for i := 0; i < numWonders; i++ {
		wonderIDs[i] = WonderID(i)
	}
}

func shuffleWonders(rnd *rand.Rand) []WonderID {
	var wonders = wonderIDs
	rnd.Shuffle(len(wonders), func(i, j int) {
		wonders[i], wonders[j] = wonders[j], wonders[i]
	})
	return wonders[:]
}

var _ = [1]struct{}{}[len(shuffleWonders(zeroRand()))-numWonders]

var (
	listWonders = []Wonder{
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
	_ = [1]struct{}{}[numWonders-len(listWonders)]

	mapWonders = makeMapWondersByName(listWonders)
	_          = [1]struct{}{}[numWonders-len(mapWonders)]
)

func newWonder(name WonderName, args ...interface{}) (w Wonder) {
	w.Name = name
	for i := range args {
		switch a := args[i].(type) {
		case Cost:
			w.Cost = a
		case VP:
			w.Effects = append(w.Effects, typedVP{a, WonderVP})
		case Effect:
			w.Effects = append(w.Effects, a)
		default:
			panic(fmt.Sprintf("Not allow for PToken builder: %T", a))
		}
	}
	return w
}

func makeMapWondersByName(list []Wonder) map[WonderName]*Wonder {
	m := map[WonderName]*Wonder{}
	for i, w := range list {
		m[w.Name] = &list[i]
	}
	return m
}

// WonderNames - list on wonder's names
type WonderNames []WonderName

// IsExists - the name is exists in current list
func (ws WonderNames) IsExists(name WonderName) bool {
	for _, w := range ws {
		if w == name {
			return true
		}
	}
	return false
}

// IsExistsAll - the names are exist in current list
func (ws WonderNames) IsExistsAll(checkedNames WonderNames) bool {
	for _, w := range ws {
		if !checkedNames.IsExists(w) {
			return false
		}
	}
	return true
}
