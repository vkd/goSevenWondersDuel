package core

import (
	"fmt"
	"math/rand"
)

// Card - In 7 Wonders Duel, all of the Age and Guild cards represent Buildings.
// The Building cards all consist of a name, an effect and a construction cost.
type Card struct {
	ID    CardID
	Name  CardName
	Color CardColor
	Cost  Cost

	// -----------------

	Effects []Effect

	EndEffects []Finaler
}

// CardID - ID of the card
type CardID uint32

func (c CardID) card() *Card {
	return &cards[c]
}

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
	Brown         CardColor = iota // Raw materials
	Grey                           // Manufactured goods
	Blue                           // Civilian Buildings
	Green                          // Scientific Buildings
	Yellow                         // Commercial Buildings
	Red                            // Military Buildings
	Purple                         // Guilds
	numCardColors = iota
)

var (
	namesCardColor = map[CardColor]string{
		Brown:  "Brown",
		Grey:   "Grey",
		Blue:   "Blue",
		Green:  "Green",
		Yellow: "Yellow",
		Red:    "Red",
		Purple: "Purple",
	}
	_ = [1]struct{}{}[len(namesCardColor)-numCardColors]
)

// String representation of card color
func (c CardColor) String() string { return namesCardColor[c] }

// -----------------------------------

const (
	numAgeI   = 23
	numAgeII  = 23
	numAgeIII = 20
	numGuilds = 7

	dropCardsFromEveryAge = 3
	takeGuildsToAgeIII    = 3

	totalNum = numAgeI + numAgeII + numAgeIII + numGuilds

	SizeAge = 20
)

var (
	ageICardIDs   [numAgeI]CardID
	ageIICardIDs  [numAgeII]CardID
	ageIIICardIDs [numAgeIII]CardID
	guildsCardIDs [numGuilds]CardID
)

func init() {
	for i := 0; i < numAgeI; i++ {
		ageICardIDs[i] = CardID(i)
	}
	for i := 0; i < numAgeII; i++ {
		ageIICardIDs[i] = CardID(numAgeI + i)
	}
	for i := 0; i < numAgeIII; i++ {
		ageIIICardIDs[i] = CardID(numAgeI + numAgeII + i)
	}
	for i := 0; i < numGuilds; i++ {
		guildsCardIDs[i] = CardID(numAgeI + numAgeII + numAgeIII + i)
	}
}

func shuffleAgeI(rnd *rand.Rand) []CardID {
	var cards = ageICardIDs
	shuffleCards(rnd, cards[:])
	return cards[:numAgeI-dropCardsFromEveryAge]
}

var _ = [1]struct{}{}[len(shuffleAgeI(zeroRand()))-SizeAge]

func shuffleAgeII(rnd *rand.Rand) []CardID {
	var cards = ageIICardIDs
	shuffleCards(rnd, cards[:])
	return cards[:numAgeII-dropCardsFromEveryAge]
}

var _ = [1]struct{}{}[len(shuffleAgeII(zeroRand()))-SizeAge]

func shuffleAgeIII(rnd *rand.Rand) (out []CardID) {
	var cards = ageIIICardIDs
	shuffleCards(rnd, cards[:])
	out = append(out, cards[:numAgeIII-dropCardsFromEveryAge]...)

	var guilds = guildsCardIDs
	shuffleCards(rnd, guilds[:])
	out = append(out, guilds[:takeGuildsToAgeIII]...)

	shuffleCards(rnd, out)
	return out
}

var _ = [1]struct{}{}[len(shuffleAgeIII(zeroRand()))-SizeAge]

func shuffleCards(rnd *rand.Rand, cards []CardID) {
	rnd.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
}

var (
	listAgeI = []Card{
		newCard("Stable", Red, Shields(1), Horseshoe, NewCost(Wood)),
		newCard("Garrison", Red),    //, Shields(1), Sword, NewCost(Clay)),
		newCard("Palisade", Red),    //, Shields(1), Wall, NewCost(Coins(2))),
		newCard("Guard tower", Red), //, Shields(1)),

		newCard("Workshop", Green),                   //, Tool, VP(1), NewCost(Papyrus)),
		newCard("Scriptorium", Green),                //, Pen, Book, NewCost(Coins(2))),
		newCard("Apothecary", Green, NewCost(Glass)), //, Wheel, VP(1)),
		newCard("Pharmacist", Green),                 //, Mortar, Gear, NewCost(Coins(2))),

		newCard("Tavern", Yellow),        //, Coins(4), Bottle),
		newCard("Clay reserve", Yellow),  //, OneCoinPrice(Clay), NewCost(Coins(3))),
		newCard("Stone reserve", Yellow), //, OneCoinPrice(Stone), NewCost(Coins(3))),
		newCard("Wood reserve", Yellow),  //, OneCoinPrice(Wood), NewCost(Coins(3))),

		newCard("Baths", Blue, VP(3), Water, NewCost(Stone)),
		newCard("Altar", Blue),   //, VP(3), Moon),
		newCard("Theater", Blue), //, VP(3), Mask),
		newCard("Lumber yard", Brown, Wood),

		newCard("Stone pit", Brown, Stone, NewCost(Coins(1))),
		newCard("Clay pool", Brown), //, Clay),
		newCard("Clay pit", Brown),  //, Clay, NewCost(Coint(1))),
		newCard("Quarry", Brown),    //, Stone),

		newCard("Logging camp", Brown), //, Wood, NewCost(Coint(1))),
		newCard("Press", Grey),         //, Papyrus, NewCost(Coint(1))),
		newCard("Glassworks", Grey),    //, Glass, NewCost(Coint(1))),
	}
	_ = [1]struct{}{}[len(listAgeI)-numAgeI]

	listAgeII = []Card{
		newCard("Barracks", Red), //, Shields(1), NewCost(Coint(3), Sword)),
		newCard("Horse breeders", Red, Shields(1), NewCost(Clay, Wood), FreeChain(Horseshoe)),
		newCard("Walls", Red),         //, Shields(2), NewCost(Stone, Stone)),
		newCard("Parade ground", Red), //, Shields(2), Helm, NewCost(Clay, Clay, Glass)),
		newCard("Archery range", Red), //, Shields(2), Target, NewCost(Stone, Wood, Papyrus)),

		newCard("Laboratory", Green), //, Tool, VP(1), Lamp, NewCost(Wood, Glass, Glass)),
		newCard("Dispensary", Green), //, Mortar, VP(2), NewCost(Clay, Clay, Stone, Gear)),
		newCard("School", Green),     //, Wheel, VP(1), Harp, NewCost(Wood, Papyrus, Papyrus)),
		newCard("Library", Green),    //, Pen, VP(2), NewCost(Stone, Wood, Glass, Book)),

		newCard("Customs house", Yellow), //, OnePriceMarket{Papyrus, 1}, OnePriceMarket{Glass, 1}, NewCost(Coint(4))),
		newCard("Brewery", Yellow),       //, Coint(6), Barrel),
		newCard("Forum", Yellow),         //, OneOfAnyMarket(manufacturedGoods), NewCost(Coint(3), Clay)),
		newCard("Caravansery", Yellow),   //, OneOfAnyMarket(rawMaterials), NewCost(Coint(2), Glass, Papyrus)),

		newCard("Temple", Blue),  //, VP(4), Sun, NewCost(Wood, Papyrus, Moon)),
		newCard("Postrum", Blue), //, VP(4), Pantheon, NewCost(Stone, Wood)),
		newCard("Aqueduct", Blue, NewCost(Stone, Stone, Stone), FreeChain(Water)), //, VP(5), ),
		newCard("Tribunal", Blue), //, VP(5), NewCost(Wood, Wood, Glass)),
		newCard("Statue", Blue),   //, VP(4), Column, NewCost(Clay, Clay, Mask)),

		newCard("Sawmill", Brown),      //, Wood, Wood, NewCost(Coint(2))),
		newCard("Shelf quarry", Brown), //, Stone, Stone, NewCost(Coint(2))),
		newCard("Brick yard", Brown),   //, Clay, Clay, NewCost(Coint(2))),

		newCard("Drying room", Grey),  //, Papyrus),
		newCard("Glass blower", Grey), //, Glass),
	}
	_ = [1]struct{}{}[len(listAgeII)-numAgeII]

	listAgeIII = []Card{
		newCard("Siege workshop", Red), //, Shields(2), NewCost(Wood, Wood, Wood, Glass, Target)),
		newCard("Fortifications", Red, NewCost(Stone, Stone, Clay, Papyrus), FreeChain(Wall)), //, Shields(2), ),
		newCard("Circus", Red),     //, Shields(2), NewCost(Clay, Clay, Stone, Stone, Helm)),
		newCard("Arsenal", Red),    //, Shields(3), NewCost(Clay, Clay, Clay, Wood, Wood)),
		newCard("Courthouse", Red), //, Shields(3), NewCost(Coint(3))),

		newCard("University", Green),  //, Astronomy, VP(2), NewCost(Clay, Glass, Papyrus, Harp)),
		newCard("Observatory", Green), //, Astronomy, VP(2), NewCost(Stone, Papyrus, Papyrus, Lamp)),
		newCard("Academy", Green),     //, Clock, VP(3), NewCost(Stone, Wood, Glass, Glass)),
		newCard("Study", Green),       //, Clock, VP(3), NewCost(Wood, Wood, Glass, Papyrus)),

		newCard("Chamber of commerce", Yellow), //, MoneyByCards{Grey, 3}, VP(3), NewCost(Papyrus, Papyrus)),
		newCard("Arena", Yellow, CoinsPerWonder(2), VP(3), NewCost(Clay, Stone, Wood), FreeChain(Barrel)),
		newCard("Port", Yellow),       //, MoneyByCards{Brown, 2}, VP(3), NewCost(Wood, Glass, Papyrus)),
		newCard("Armory", Yellow),     //, MoneyByCards{Red, 1}, VP(3), NewCost(Stone, Stone, Glass)),
		newCard("Lighthouse", Yellow), //, MoneyByCards{Yellow, 1}, VP(3), NewCost(Clay, Clay, Glass, Bottle)),

		newCard("Pantheon", Blue),  //, VP(6), NewCost(Clay, Wood, Papyrus, Papyrus, Sun)),
		newCard("Palace", Blue),    //, VP(7), NewCost(Clay, Stone, Wood, Glass, Glass)),
		newCard("Gardens", Blue),   //, VP(6), NewCost(Clay, Clay, Wood, Wood, Column)),
		newCard("Obelisk", Blue),   //, VP(5), NewCost(Stone, Stone, Glass)),
		newCard("Senate", Blue),    //, VP(5), NewCost(Clay, Clay, Stone, Papyrus, Pantheon)),
		newCard("Town hall", Blue), //, VP(7), NewCost(Stone, Stone, Stone, Wood, Wood)),
	}
	_ = [1]struct{}{}[len(listAgeIII)-numAgeIII]

	listGuilds = []Card{
		newCard("Merchants guild", Purple),  // , ...),
		newCard("Shipowners guild", Purple), // , ...),
		newCard("Builders guild", Purple, BuildersGuild(), NewCost(Stone, Stone, Clay, Wood, Glass)),
		newCard("Magistrates guild", Purple),  // , ...),
		newCard("Scientists guild", Purple),   // , ...),
		newCard("Moneylenders guild", Purple), // , ...),
		newCard("Tacticians guild", Purple),   // , ...),
	}
	_ = [1]struct{}{}[len(listGuilds)-numGuilds]

	cards = appendAll(listAgeI, listAgeII, listAgeIII, listGuilds)
	_     = [1]struct{}{}[len(cards)-totalNum]

	mapCards = makeMapCardsByName(cards)
	_        = [1]struct{}{}[len(mapCards)-totalNum]
)

func init() {
	for i := range cards {
		cards[i].ID = CardID(i)
	}
}

func newCard(name CardName, ct CardColor, args ...interface{}) (c Card) {
	c.Name = name
	c.Color = ct

	for _, arg := range args {
		switch arg := arg.(type) {
		case Cost:
			if c.Cost != nil {
				c.Cost = orCost{c.Cost, arg}
			} else {
				c.Cost = arg
			}
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

func appendAll(lists ...[]Card) []Card {
	var out []Card
	for _, l := range lists {
		out = append(out, l...)
	}
	return out
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
