package core

import (
	"fmt"
	"math/rand"
)

// Card - In 7 Wonders Duel, all of the Age and Guild cards represent Buildings.
// The Building cards all consist of a name, an effect and a construction cost.
type Card struct {
	Name  CardName
	Color CardColor
	Cost  Cost

	// -----------------

	Effects []Effect
}

func (c Card) discard(g *Game, i PlayerIndex) {
	for _, e := range c.Effects {
		switch e := e.(type) {
		case Resource:
			g.players[i].Resources.reduceOne(e)
		default:
			panic(fmt.Sprintf("Unknown effect type for discard: %T", e))
		}
	}
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
		newCard("Garrison", Red, Shields(1), Sword, NewCost(Clay)),
		newCard("Palisade", Red, Shields(1), Wall, NewCost(Coins(2))),
		newCard("Guard tower", Red, Shields(1)),

		newCard("Workshop", Green, Tool, VP(1), NewCost(Papyrus)),
		newCard("Scriptorium", Green, Pen, Book, NewCost(Coins(2))),
		newCard("Apothecary", Green, Wheel, VP(1), NewCost(Glass)),
		newCard("Pharmacist", Green, Mortar, Gear, NewCost(Coins(2))),

		newCard("Tavern", Yellow, Coins(4), Bottle),
		newCard("Clay reserve", Yellow, OneCoinPrice(Clay), NewCost(Coins(3))),
		newCard("Stone reserve", Yellow, OneCoinPrice(Stone), NewCost(Coins(3))),
		newCard("Wood reserve", Yellow, OneCoinPrice(Wood), NewCost(Coins(3))),

		newCard("Baths", Blue, VP(3), Water, NewCost(Stone)),
		newCard("Altar", Blue, VP(3), Moon),
		newCard("Theater", Blue, VP(3), Mask),
		newCard("Lumber yard", Brown, Wood),

		newCard("Stone pit", Brown, Stone, NewCost(Coins(1))),
		newCard("Clay pool", Brown, Clay),
		newCard("Clay pit", Brown, Clay, NewCost(Coins(1))),
		newCard("Quarry", Brown, Stone),

		newCard("Logging camp", Brown, Wood, NewCost(Coins(1))),
		newCard("Press", Grey, Papyrus, NewCost(Coins(1))),
		newCard("Glassworks", Grey, Glass, NewCost(Coins(1))),
	}
	_ = [1]struct{}{}[len(listAgeI)-numAgeI]

	listAgeII = []Card{
		newCard("Barracks", Red, Shields(1), NewCost(Coins(3)), FreeChain(Sword)),

		newCard("Horse breeders", Red, Shields(1), NewCost(Clay, Wood), FreeChain(Horseshoe)),
		newCard("Walls", Red, Shields(2), NewCost(Stone, Stone)),
		newCard("Parade ground", Red, Shields(2), Helm, NewCost(Clay, Clay, Glass)),
		newCard("Archery range", Red, Shields(2), Target, NewCost(Stone, Wood, Papyrus)),

		newCard("Laboratory", Green, Tool, VP(1), Lamp, NewCost(Wood, Glass, Glass)),
		newCard("Dispensary", Green, Mortar, VP(2), NewCost(Clay, Clay, Stone), FreeChain(Gear)),
		newCard("School", Green, Wheel, VP(1), Harp, NewCost(Wood, Papyrus, Papyrus)),
		newCard("Library", Green, Pen, VP(2), NewCost(Stone, Wood, Glass), FreeChain(Book)),

		newCard("Customs house", Yellow, OneCoinPrice(Papyrus), OneCoinPrice(Glass), NewCost(Coins(4))),
		newCard("Brewery", Yellow, Coins(6), Barrel),
		newCard("Forum", Yellow, OneManufacturedMarket(), NewCost(Coins(3), Clay)),
		newCard("Caravansery", Yellow, OneRawMarket(), NewCost(Coins(2), Glass, Papyrus)),

		newCard("Temple", Blue, VP(4), Sun, NewCost(Wood, Papyrus), FreeChain(Moon)),
		newCard("Postrum", Blue, VP(4), Pantheon, NewCost(Stone, Wood)),
		newCard("Aqueduct", Blue, VP(5), NewCost(Stone, Stone, Stone), FreeChain(Water)),
		newCard("Tribunal", Blue, VP(5), NewCost(Wood, Wood, Glass)),
		newCard("Statue", Blue, VP(4), Column, NewCost(Clay, Clay), FreeChain(Mask)),

		newCard("Sawmill", Brown, Wood, Wood, NewCost(Coins(2))),
		newCard("Shelf quarry", Brown, Stone, Stone, NewCost(Coins(2))),
		newCard("Brick yard", Brown, Clay, Clay, NewCost(Coins(2))),

		newCard("Drying room", Grey, Papyrus),
		newCard("Glass blower", Grey, Glass),
	}
	_ = [1]struct{}{}[len(listAgeII)-numAgeII]

	listAgeIII = []Card{
		newCard("Siege workshop", Red, Shields(2), NewCost(Wood, Wood, Wood, Glass), FreeChain(Target)),
		newCard("Fortifications", Red, Shields(2), NewCost(Stone, Stone, Clay, Papyrus), FreeChain(Wall)),
		newCard("Circus", Red, Shields(2), NewCost(Clay, Clay, Stone, Stone), FreeChain(Helm)),
		newCard("Arsenal", Red, Shields(3), NewCost(Clay, Clay, Clay, Wood, Wood)),
		newCard("Courthouse", Red, Shields(3), NewCost(Coins(8))),

		newCard("University", Green, Astronomy, VP(2), NewCost(Clay, Glass, Papyrus), FreeChain(Harp)),
		newCard("Observatory", Green, Astronomy, VP(2), NewCost(Stone, Papyrus, Papyrus), FreeChain(Lamp)),
		newCard("Academy", Green, Clock, VP(3), NewCost(Stone, Wood, Glass, Glass)),
		newCard("Study", Green, Clock, VP(3), NewCost(Wood, Wood, Glass, Papyrus)),

		newCard("Chamber of commerce", Yellow, CoinsPerCard(Grey, 3), VP(3), NewCost(Papyrus, Papyrus)),
		newCard("Arena", Yellow, CoinsPerWonder(2), VP(3), NewCost(Clay, Stone, Wood), FreeChain(Barrel)),
		newCard("Port", Yellow, CoinsPerCard(Brown, 2), VP(3), NewCost(Wood, Glass, Papyrus)),
		newCard("Armory", Yellow, CoinsPerCard(Red, 1), VP(3), NewCost(Stone, Stone, Glass)),
		newCard("Lighthouse", Yellow, CoinsPerCard(Yellow, 1), VP(3), NewCost(Clay, Clay, Glass), FreeChain(Bottle)),

		newCard("Pantheon", Blue, VP(6), NewCost(Clay, Wood, Papyrus, Papyrus), FreeChain(Sun)),
		newCard("Palace", Blue, VP(7), NewCost(Clay, Stone, Wood, Glass, Glass)),
		newCard("Gardens", Blue, VP(6), NewCost(Clay, Clay, Wood, Wood), FreeChain(Column)),
		newCard("Obelisk", Blue, VP(5), NewCost(Stone, Stone, Glass)),
		newCard("Senate", Blue, VP(5), NewCost(Clay, Clay, Stone, Papyrus), FreeChain(Pantheon)),
		newCard("Town hall", Blue, VP(7), NewCost(Stone, Stone, Stone, Wood, Wood)),
	}
	_ = [1]struct{}{}[len(listAgeIII)-numAgeIII]

	listGuilds = []Card{
		newCard("Merchants guild", Purple, MaxOneCoinPerCards(Yellow), MaxOneVPPerCards(Yellow), NewCost(Clay, Wood, Glass, Papyrus)),
		newCard("Shipowners guild", Purple, MaxOneCoinPerCards(Brown, Grey), MaxOneVPPerCards(Brown, Grey), NewCost(Clay, Stone, Glass, Papyrus)),
		newCard("Builders guild", Purple, BuildersGuild(), NewCost(Stone, Stone, Clay, Wood, Glass)),
		newCard("Tacticians guild", Purple, MaxOneCoinPerCards(Red), MaxOneVPPerCards(Red), NewCost(Stone, Stone, Clay, Papyrus)),
		newCard("Moneylenders guild", Purple, MoneylendersGuild(), NewCost(Stone, Stone, Wood, Wood)),
		newCard("Scientists guild", Purple, MaxOneCoinPerCards(Green), MaxOneVPPerCards(Green), NewCost(Clay, Clay, Wood, Wood)),
		newCard("Magistrates guild", Purple, MaxOneCoinPerCards(Blue), MaxOneVPPerCards(Blue), NewCost(Wood, Wood, Clay, Papyrus)),
	}
	_ = [1]struct{}{}[len(listGuilds)-numGuilds]

	cards = appendAll(listAgeI, listAgeII, listAgeIII, listGuilds)
	_     = [1]struct{}{}[len(cards)-totalNum]

	mapCards = makeMapCardsByName(cards)
	_        = [1]struct{}{}[len(mapCards)-totalNum]
)

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
		case VP:
			c.Effects = append(c.Effects, typedVP{arg, VPTypeByColor(ct)})
		case maxVPsPerCards:
			arg.Type = VPTypeByColor(ct)
			c.Effects = append(c.Effects, arg)
		case vPsPerWonder:
			arg.Type = VPTypeByColor(ct)
			c.Effects = append(c.Effects, arg)
		case vPPerCoins:
			arg.Type = VPTypeByColor(ct)
			c.Effects = append(c.Effects, arg)
		case Effect:
			c.Effects = append(c.Effects, arg)
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
