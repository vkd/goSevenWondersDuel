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

	Cost                CostOfCard
	FreeCostChainSymbol MaybeChainSym
}

// CardIndex on all cards
// type CardIndex ItemIndex

// CardName - name of card
type CardName string

func (n CardName) find() *Card {
	return mapCards[n]
}

// const (
// 	// InvisibleCard index
// 	InvisibleCard CardIndex = 99999
// )

var (
	ageI = []Card{
		newCard("Stable", Red, Shields(1), Horseshoe, Cost(Wood)),
		newCard("Garrison", Red, Shields(1), Sword, Cost(Clay)),
		newCard("Palisade", Red, Shields(1), Wall, Cost(Money(2))),
		newCard("Guard tower", Red, Shields(1)),

		newCard("Workshop", Green, Tool, VP(1), Cost(Papyrus)),
		newCard("Scriptorium", Green, Pen, Book, Cost(Money(2))),
		newCard("Apothecary", Green, Wheel, VP(1), Cost(Glass)),
		newCard("Pharmacist", Green, Mortar, Gear, Cost(Money(2))),

		newCard("Tavern", Yellow, Money(4), Bottle),
		newCard("Clay reserve", Yellow, OnePriceMarket{Clay, 1}, Cost(Money(3))),
		newCard("Stone reserve", Yellow, OnePriceMarket{Stone, 1}, Cost(Money(3))),
		newCard("Wood reserve", Yellow, OnePriceMarket{Wood, 1}, Cost(Money(3))),

		newCard("Baths", Blue, VP(3), Water, Cost(Stone)),
		newCard("Altar", Blue, VP(3), Moon),
		newCard("Theater", Blue, VP(3), Mask),
		newCard("Lumber yard", Brown, Wood),

		newCard("Stone pit", Brown, Stone, Cost(Money(1))),
		newCard("Clay pool", Brown, Clay),
		newCard("Clay pit", Brown, Clay, Cost(Money(1))),
		newCard("Quarry", Brown, Stone),

		newCard("Logging camp", Brown, Wood, Cost(Money(1))),
		newCard("Press", Grey, Papyrus, Cost(Money(1))),
		newCard("Glassworks", Grey, Glass, Cost(Money(1))),
	}
	_ = [1]struct{}{}[len(ageI)-23]

	ageII = []Card{
		newCard("Barracks", Red, Shields(1), Cost(Money(3)), Cost(Sword)),
		newCard("Horse breeders", Red, Shields(1), Cost(Clay, Wood), Cost(Horseshoe)),
		newCard("Walls", Red, Shields(2), Cost(Stone, Stone)),
		newCard("Parade ground", Red, Shields(2), Helm, Cost(Clay, Clay, Glass)),
		newCard("Archery range", Red, Shields(2), Target, Cost(Stone, Wood, Papyrus)),

		newCard("Laboratory", Green, Tool, VP(1), Lamp, Cost(Wood, Glass, Glass)),
		newCard("Dispensary", Green, Mortar, VP(2), Cost(Clay, Clay, Stone), Cost(Gear)),
		newCard("School", Green, Wheel, VP(1), Harp, Cost(Wood, Papyrus, Papyrus)),
		newCard("Library", Green, Pen, VP(2), Cost(Stone, Wood, Glass), Cost(Book)),

		newCard("Customs house", Yellow, OnePriceMarket{Papyrus, 1}, OnePriceMarket{Glass, 1}, Cost(Money(4))),
		newCard("Brewery", Yellow, Money(6), Barrel),
		newCard("Forum", Yellow, OneOfAnyMarket(manufacturedGoods), Cost(Money(3), Clay)),
		newCard("Caravansery", Yellow, OneOfAnyMarket(rawMaterials), Cost(Money(2), Glass, Papyrus)),

		newCard("Temple", Blue, VP(4), Sun, Cost(Wood, Papyrus), Cost(Moon)),
		newCard("Postrum", Blue, VP(4), Pantheon, Cost(Stone, Wood)),
		newCard("Aqueduct", Blue, VP(5), Cost(Stone, Stone, Stone), Cost(Water)),
		newCard("Tribunal", Blue, VP(5), Cost(Wood, Wood, Glass)),
		newCard("Statue", Blue, VP(4), Column, Cost(Clay, Clay), Cost(Mask)),

		newCard("Sawmill", Brown, Wood, Wood, Cost(Money(2))),
		newCard("Shelf quarry", Brown, Stone, Stone, Cost(Money(2))),
		newCard("Brick yard", Brown, Clay, Clay, Cost(Money(2))),

		newCard("Drying room", Grey, Papyrus),
		newCard("Glass blower", Grey, Glass),
	}
	_ = [1]struct{}{}[len(ageII)-23]

	ageIII = []Card{
		newCard("Siege workshop", Red, Shields(2), Cost(Wood, Wood, Wood, Glass), Cost(Target)),
		newCard("Fortifications", Red, Shields(2), Cost(Stone, Stone, Clay, Papyrus), Cost(Wall)),
		newCard("Circus", Red, Shields(2), Cost(Clay, Clay, Stone, Stone), Cost(Helm)),
		newCard("Arsenal", Red, Shields(3), Cost(Clay, Clay, Clay, Wood, Wood)),
		newCard("Courthouse", Red, Shields(3), Cost(Money(3))),

		newCard("University", Green, Astronomy, VP(2), Cost(Clay, Glass, Papyrus), Cost(Harp)),
		newCard("Observatory", Green, Astronomy, VP(2), Cost(Stone, Papyrus, Papyrus), Cost(Lamp)),
		newCard("Academy", Green, Clock, VP(3), Cost(Stone, Wood, Glass, Glass)),
		newCard("Study", Green, Clock, VP(3), Cost(Wood, Wood, Glass, Papyrus)),

		newCard("Chamber of commerce", Yellow, MoneyByCards{Grey, 3}, VP(3), Cost(Papyrus, Papyrus)),
		newCard("Arena", Yellow, MoneyByWonders{2}, VP(3), Cost(Clay, Stone, Wood), Cost(Barrel)),
		newCard("Port", Yellow, MoneyByCards{Brown, 2}, VP(3), Cost(Wood, Glass, Papyrus)),
		newCard("Armory", Yellow, MoneyByCards{Red, 1}, VP(3), Cost(Stone, Stone, Glass)),
		newCard("Lighthouse", Yellow, MoneyByCards{Yellow, 1}, VP(3), Cost(Clay, Clay, Glass), Cost(Bottle)),

		newCard("Pantheon", Blue, VP(6), Cost(Clay, Wood, Papyrus, Papyrus), Cost(Sun)),
		newCard("Palace", Blue, VP(7), Cost(Clay, Stone, Wood, Glass, Glass)),
		newCard("Gardens", Blue, VP(6), Cost(Clay, Clay, Wood, Wood), Cost(Column)),
		newCard("Obelisk", Blue, VP(5), Cost(Stone, Stone, Glass)),
		newCard("Senate", Blue, VP(5), Cost(Clay, Clay, Stone, Papyrus), Cost(Pantheon)),
		newCard("Town hall", Blue, VP(7), Cost(Stone, Stone, Stone, Wood, Wood)),
	}
	_ = [1]struct{}{}[len(ageIII)-20]

	guilds = []Card{
		// newCard("Merchants guild", Purple, ),
	}
	// _ = [1]struct{}{}[len(guilds)-7]

	allCards = append(append(append([]Card{}, ageI...), ageII...), ageIII...)

	mapCards = map[CardName]*Card{}
	// _ = [1]struct{}{}[len(mapCards)-23-23-20-7]
)

// var nextCardIndex CardIndex

func newCard(name CardName, ct CardColor, args ...interface{}) (c Card) {
	// c.Index = nextCardIndex
	// nextCardIndex++

	c.Name = name
	c.CardColor = ct
	for _, arg := range args {
		switch arg := arg.(type) {
		case CostOfCard:
			c.Cost = arg
		case Effect:
			c.Effects = append(c.Effects, arg)
		// default:
		// 	c.Effects = append(c.Effects, arg)
		default:
			panic(fmt.Sprintf("Not implemented: %T", arg))
		}
	}
	return
}
