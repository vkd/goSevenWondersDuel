package game

import "fmt"

// CardColor ...
type CardColor uint8

// Color of card
const (
	Brown  CardColor = iota // Raw materials
	Grey                    // Manufactured goods
	Blue                    // Civilian Buildings
	Green                   // Scientific Buildings
	Yellow                  // Commercial Buildings
	Red                     // Military Buildings
	Purple                  // Guilds
)

// Card ...
type Card struct {
	// Index     CardIndex
	Name      CardName
	CardColor CardColor
	Effects   []Effect

	Cost                Cost
	FreeCostChainSymbol MaybeChainSym
}

// CardIndex on all cards
// type CardIndex ItemIndex

// CardName - name of card
type CardName string

func (n CardName) find() *Card {
	return mapCards[n]
}

// CardNames - list of name of cards
// type CardNames []CardName

// const (
// 	// InvisibleCard index
// 	InvisibleCard CardIndex = 99999
// )

const (
	numAgeI   = 23
	numAgeII  = 23
	numAgeIII = 20
	numGuilds = 7
)

var (
	ageI = []Card{
		newCard("Stable", Red, Shields(1), Horseshoe, NewCost(Wood)),
		newCard("Garrison", Red, Shields(1), Sword, NewCost(Clay)),
		newCard("Palisade", Red, Shields(1), Wall, NewCost(Money(2))),
		newCard("Guard tower", Red, Shields(1)),

		newCard("Workshop", Green, Tool, VP(1), NewCost(Papyrus)),
		newCard("Scriptorium", Green, Pen, Book, NewCost(Money(2))),
		newCard("Apothecary", Green, Wheel, VP(1), NewCost(Glass)),
		newCard("Pharmacist", Green, Mortar, Gear, NewCost(Money(2))),

		newCard("Tavern", Yellow, Money(4), Bottle),
		newCard("Clay reserve", Yellow, OnePriceMarket{Clay, 1}, NewCost(Money(3))),
		newCard("Stone reserve", Yellow, OnePriceMarket{Stone, 1}, NewCost(Money(3))),
		newCard("Wood reserve", Yellow, OnePriceMarket{Wood, 1}, NewCost(Money(3))),

		newCard("Baths", Blue, VP(3), Water, NewCost(Stone)),
		newCard("Altar", Blue, VP(3), Moon),
		newCard("Theater", Blue, VP(3), Mask),
		newCard("Lumber yard", Brown, Wood),

		newCard("Stone pit", Brown, Stone, NewCost(Money(1))),
		newCard("Clay pool", Brown, Clay),
		newCard("Clay pit", Brown, Clay, NewCost(Money(1))),
		newCard("Quarry", Brown, Stone),

		newCard("Logging camp", Brown, Wood, NewCost(Money(1))),
		newCard("Press", Grey, Papyrus, NewCost(Money(1))),
		newCard("Glassworks", Grey, Glass, NewCost(Money(1))),
	}
	_ = [1]struct{}{}[len(ageI)-numAgeI]

	ageII = []Card{
		newCard("Barracks", Red, Shields(1), NewCost(Money(3)), FreeChainSymbol(Sword)),
		newCard("Horse breeders", Red, Shields(1), NewCost(Clay, Wood), FreeChainSymbol(Horseshoe)),
		newCard("Walls", Red, Shields(2), NewCost(Stone, Stone)),
		newCard("Parade ground", Red, Shields(2), Helm, NewCost(Clay, Clay, Glass)),
		newCard("Archery range", Red, Shields(2), Target, NewCost(Stone, Wood, Papyrus)),

		newCard("Laboratory", Green, Tool, VP(1), Lamp, NewCost(Wood, Glass, Glass)),
		newCard("Dispensary", Green, Mortar, VP(2), NewCost(Clay, Clay, Stone), FreeChainSymbol(Gear)),
		newCard("School", Green, Wheel, VP(1), Harp, NewCost(Wood, Papyrus, Papyrus)),
		newCard("Library", Green, Pen, VP(2), NewCost(Stone, Wood, Glass), FreeChainSymbol(Book)),

		newCard("Customs house", Yellow, OnePriceMarket{Papyrus, 1}, OnePriceMarket{Glass, 1}, NewCost(Money(4))),
		newCard("Brewery", Yellow, Money(6), Barrel),
		newCard("Forum", Yellow, OneOfAnyMarket(manufacturedGoods), NewCost(Money(3), Clay)),
		newCard("Caravansery", Yellow, OneOfAnyMarket(rawMaterials), NewCost(Money(2), Glass, Papyrus)),

		newCard("Temple", Blue, VP(4), Sun, NewCost(Wood, Papyrus), FreeChainSymbol(Moon)),
		newCard("Postrum", Blue, VP(4), Pantheon, NewCost(Stone, Wood)),
		newCard("Aqueduct", Blue, VP(5), NewCost(Stone, Stone, Stone), FreeChainSymbol(Water)),
		newCard("Tribunal", Blue, VP(5), NewCost(Wood, Wood, Glass)),
		newCard("Statue", Blue, VP(4), Column, NewCost(Clay, Clay), FreeChainSymbol(Mask)),

		newCard("Sawmill", Brown, Wood, Wood, NewCost(Money(2))),
		newCard("Shelf quarry", Brown, Stone, Stone, NewCost(Money(2))),
		newCard("Brick yard", Brown, Clay, Clay, NewCost(Money(2))),

		newCard("Drying room", Grey, Papyrus),
		newCard("Glass blower", Grey, Glass),
	}
	_ = [1]struct{}{}[len(ageII)-numAgeII]

	ageIII = []Card{
		newCard("Siege workshop", Red, Shields(2), NewCost(Wood, Wood, Wood, Glass), FreeChainSymbol(Target)),
		newCard("Fortifications", Red, Shields(2), NewCost(Stone, Stone, Clay, Papyrus), FreeChainSymbol(Wall)),
		newCard("Circus", Red, Shields(2), NewCost(Clay, Clay, Stone, Stone), FreeChainSymbol(Helm)),
		newCard("Arsenal", Red, Shields(3), NewCost(Clay, Clay, Clay, Wood, Wood)),
		newCard("Courthouse", Red, Shields(3), NewCost(Money(3))),

		newCard("University", Green, Astronomy, VP(2), NewCost(Clay, Glass, Papyrus), FreeChainSymbol(Harp)),
		newCard("Observatory", Green, Astronomy, VP(2), NewCost(Stone, Papyrus, Papyrus), FreeChainSymbol(Lamp)),
		newCard("Academy", Green, Clock, VP(3), NewCost(Stone, Wood, Glass, Glass)),
		newCard("Study", Green, Clock, VP(3), NewCost(Wood, Wood, Glass, Papyrus)),

		newCard("Chamber of commerce", Yellow, MoneyByCards{Grey, 3}, VP(3), NewCost(Papyrus, Papyrus)),
		newCard("Arena", Yellow, MoneyByWonders{2}, VP(3), NewCost(Clay, Stone, Wood), FreeChainSymbol(Barrel)),
		newCard("Port", Yellow, MoneyByCards{Brown, 2}, VP(3), NewCost(Wood, Glass, Papyrus)),
		newCard("Armory", Yellow, MoneyByCards{Red, 1}, VP(3), NewCost(Stone, Stone, Glass)),
		newCard("Lighthouse", Yellow, MoneyByCards{Yellow, 1}, VP(3), NewCost(Clay, Clay, Glass), FreeChainSymbol(Bottle)),

		newCard("Pantheon", Blue, VP(6), NewCost(Clay, Wood, Papyrus, Papyrus), FreeChainSymbol(Sun)),
		newCard("Palace", Blue, VP(7), NewCost(Clay, Stone, Wood, Glass, Glass)),
		newCard("Gardens", Blue, VP(6), NewCost(Clay, Clay, Wood, Wood), FreeChainSymbol(Column)),
		newCard("Obelisk", Blue, VP(5), NewCost(Stone, Stone, Glass)),
		newCard("Senate", Blue, VP(5), NewCost(Clay, Clay, Stone, Papyrus), FreeChainSymbol(Pantheon)),
		newCard("Town hall", Blue, VP(7), NewCost(Stone, Stone, Stone, Wood, Wood)),
	}
	_ = [1]struct{}{}[len(ageIII)-numAgeIII]

	guilds = []Card{
		// newCard("Merchants guild", Purple, ),
	}
	// _ = [1]struct{}{}[len(guilds)-numGuilds]

	// allCards = append(append(append(append([]Card{}, ageI...), ageII...), ageIII...), guilds...)
	// _        = [1]struct{}{}[len(allCards)-len(ageI)-len(ageII)-len(ageIII)-len(guilds)]

	mapCards = makeMapCardsByName(ageI, ageII, ageIII, guilds)
	_        = [1]struct{}{}[len(mapCards)-len(ageI)-len(ageII)-len(ageIII)-len(guilds)]
)

// var nextCardIndex CardIndex

func newCard(name CardName, ct CardColor, args ...interface{}) (c Card) {
	// c.Index = nextCardIndex
	// nextCardIndex++

	c.Name = name
	c.CardColor = ct
	for _, arg := range args {
		switch arg := arg.(type) {
		case Cost:
			c.Cost = arg
		case Effect:
			c.Effects = append(c.Effects, arg)
		case FreeChainSymbol:
			c.FreeCostChainSymbol.Set(ChainSymbol(arg))
		// default:
		// 	c.Effects = append(c.Effects, arg)
		default:
			panic(fmt.Sprintf("Not implemented: %T", arg))
		}
	}
	return
}

func makeMapCardsByName(cc ...[]Card) map[CardName]*Card {
	m := map[CardName]*Card{}
	for i := range cc {
		for j, c := range cc[i] {
			if _, ok := m[c.Name]; ok {
				panic("%q card already exists in map")
			}
			m[c.Name] = &cc[i][j]
		}
	}
	return m
}
