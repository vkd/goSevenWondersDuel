package core

import (
	"fmt"
)

// Card - In 7 Wonders Duel, all of the Age and Guild cards represent Buildings.
// The Building cards all consist of a name, an effect and a construction cost.
type Card struct {
	Name    CardName
	Color   CardColor
	Effects []Effect
	Cost    Cost

	EndEffects []Finaler
}

type CardID uint32

// CardName - name of card
type CardName string

func (n CardName) card() *Card {
	c, ok := mapCards[n]
	if !ok {
		panic(fmt.Sprintf("Cannot find %q card", n))
	}
	return c
}

// CardColor - There are 7 different types of Buildings, easily identifiable by their colored border.
type CardColor uint8

// Color of card
const (
	Brown           CardColor = iota // Raw materials
	Grey                             // Manufactured goods
	Blue                             // Civilian Buildings
	Green                            // Scientific Buildings
	Yellow                           // Commercial Buildings
	Red                              // Military Buildings
	Purple                           // Guilds
	numOfCardColors = iota
)

var (
	nameCardColor = map[CardColor]string{
		Brown:  "Brown",
		Grey:   "Grey",
		Blue:   "Blue",
		Green:  "Green",
		Yellow: "Yellow",
		Red:    "Red",
		Purple: "Purple",
	}
	_ = [1]struct{}{}[len(nameCardColor)-numOfCardColors]
)

// String representation of card color
func (c CardColor) String() string { return nameCardColor[c] }

const (
	numAgeI   = 23
	numAgeII  = 23
	numAgeIII = 20
	numGuilds = 7

	totalNum = numAgeI + numAgeII + numAgeIII + numGuilds

	SizeAge = 20
)

var (
	ageI = []Card{
		newCard("Stable", Red, Shields(1), Horseshoe, NewCost(Wood)),

		newCard("Baths", Blue, VP(3), Water, NewCost(Stone)),
		newCard("Lumber yard", Brown, Wood),

		newCard("Stone pit", Brown, Stone, NewCost(Coins(1))),
	}
	// _ = [1]struct{}{}[len(ageI)-numAgeI]

	ageII = []Card{
		newCard("Horse breeders", Red, Shields(1), NewCost(Clay, Wood, Horseshoe)),
	}
	// _ = [1]struct{}{}[len(ageII)-numAgeII]

	ageIII = []Card{
		newCard("Arena", Yellow, CoinsPerWonder(2), VP(3), NewCost(Clay, Stone, Wood, Barrel)),
	}
	// _ = [1]struct{}{}[len(ageIII)-numAgeIII]

	guilds = []Card{
	}
	// _ = [1]struct{}{}[len(guilds)-numGuilds]

	mapCards = makeMapCardsByName(ageI, ageII, ageIII, guilds)
	// _        = [1]struct{}{}[len(mapCards)-totalNum]
)

func newCard(name CardName, ct CardColor, args ...interface{}) (c Card) {
	c.Name = name
	c.Color = ct

	for _, arg := range args {
		switch arg := arg.(type) {
		case Cost:
			c.Cost = arg
		case Effect:
			c.Effects = append(c.Effects, arg)
		case Finaler:
			c.EndEffects = append(c.EndEffects, arg)
		default:
			panic(fmt.Sprintf("Not allow for card builder: %T", arg))
		}
	}
	return
}

func makeMapCardsByName(cc ...[]Card) map[CardName]*Card {
	m := map[CardName]*Card{}
	for i := range cc {
		for j, c := range cc[i] {
			if _, ok := m[c.Name]; ok {
				panic(fmt.Sprintf("%q card already exists in map", c.Name))
			}
			m[c.Name] = &cc[i][j]
		}
	}
	return m
}
