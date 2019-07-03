package game

// Money ...
type Money int

func (m Money) effect() {}

// var _ Price = Money(0)

// func (m Money) pricer() {}

// SetMin ...
// func (m *Money) SetMin(new Money) {
// 	if new < *m {
// 		*m = new
// 	}
// }

// Sub ...
func (m *Money) Sub(money DiscardMoney) {
	if *m < Money(money) {
		*m = 0
	} else {
		*m -= Money(money)
	}
}

// Mul ...
func (m Money) Mul(x int) Money {
	return Money(x) * m
}

// type Building interface{}

// Score ...
type Score int

// ScientificSymbol ...
type ScientificSymbol string

func (ScientificSymbol) effect() {}

// ScientificSymbols ...
type ScientificSymbols map[ScientificSymbol]int

// ScientificSymbols
const (
	Wheel     ScientificSymbol = "Wheel"
	Mortar    ScientificSymbol = "Mortar"
	Clock     ScientificSymbol = "Clock"
	Tool      ScientificSymbol = "Tool"
	Pen       ScientificSymbol = "Pen"
	Astronomy ScientificSymbol = "Astronomy"
	Scales    ScientificSymbol = "Scales"
)

// ChainSymbol ...
type ChainSymbol string

func (ChainSymbol) effect() {}

// ChainSymbols
const (
	Horseshoe ChainSymbol = "Horseshoe"
	Sword     ChainSymbol = "Sword"
	Wall      ChainSymbol = "Wall"
	Target    ChainSymbol = "Target"
	Helm      ChainSymbol = "Helm"
	Book      ChainSymbol = "Book"
	Gear      ChainSymbol = "Gear"
	Harp      ChainSymbol = "Harp"
	Lamp      ChainSymbol = "Lamp"
	Mask      ChainSymbol = "Mask"
	Column    ChainSymbol = "Column"
	Moon      ChainSymbol = "Moon"
	Sun       ChainSymbol = "Sun"
	Water     ChainSymbol = "Water"
	Pantheon  ChainSymbol = "Pantheon"
	Bottle    ChainSymbol = "Bottle"
	Barrel    ChainSymbol = "Barrel"
)

// ChainSymbols ...
type ChainSymbols []ChainSymbol

// Exists ...
func (c ChainSymbols) Exists(cs ChainSymbol) bool {
	for _, cc := range c {
		if cc == cs {
			return true
		}
	}
	return false
}

// FreeChainSymbol - free cost by that chain symbol
type FreeChainSymbol ChainSymbol

// Effect ...
type Effect interface {
	effect()
}

type opponent struct {
	e Effect
}

func (opponent) effect() {}

// Opponent ...
func Opponent(e Effect) Effect {
	return opponent{e}
}

// Effects ...
func Effects(args ...Effect) []Effect {
	return args
}

// Cost ...
// type Cost []Price

// VP ...
type VP uint

func (VP) effect() {}

// Market ...
type Market interface{}

// PriceMarket ...
// type PriceMarket interface {
// 	Price(r Resource) Money
// }

// OneOfAnyMarket ...
type OneOfAnyMarket []Resource

func (OneOfAnyMarket) effect() {}

// OnePriceMarket ...
type OnePriceMarket struct {
	Res   Resource
	Price Money
}

func (OnePriceMarket) effect() {}

// func OnePriceMarket(r Resource, price Money) PriceMarket {
// 	panic("Not implemented")
// }

// type payment struct {
// 	// resources Resources

// 	cardIndex int
// 	money     int
// 	// moneyRequirements int
// }

// type payments []payment

// MaybeChainSym ...
type MaybeChainSym struct {
	ChainSymbol ChainSymbol
	Exists      bool
}

// Set ...
func (m *MaybeChainSym) Set(c ChainSymbol) {
	m.ChainSymbol = c
	m.Exists = true
}

// DiscardMoney ...
type DiscardMoney Money

func (DiscardMoney) effect() {}

// VictoryType ...
type VictoryType uint8

// Types of victory
const (
	CivilianVictory VictoryType = iota
	MilitarySupremacy
	ScientificSupremacy
)

// ItemIndex ...
type ItemIndex int

// GameCost ...
// type GameCost CostOfCard

// MoneyByCards ...
type MoneyByCards struct {
	Color CardColor
	Value Money
}

func (MoneyByCards) effect() {}

// MoneyByWonders ...
type MoneyByWonders struct {
	Value Money
}

func (MoneyByWonders) effect() {}

// RepeatTurn ...
func RepeatTurn() Effect {
	return repeatTurn{}
}

type repeatTurn struct{}

func (repeatTurn) effect() {}
