package game

// Money ...
type Money int

func (m Money) effect(g *Game, i PlayerIndex) {
	g.player(i).Money += m
}

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

func (s ScientificSymbol) effect(g *Game, i PlayerIndex) {
	player := g.player(i)
	if player.ScientificSymbols == nil {
		player.ScientificSymbols = make(ScientificSymbols)
	}
	player.ScientificSymbols[s]++
	if player.ScientificSymbols[s]%2 == 0 {
		g.canSelectActiveProgressToken()
	}
}

// ScientificSymbols ...
type ScientificSymbols map[ScientificSymbol]int

func (s ChainSymbol) effect(g *Game, i PlayerIndex) {
	player := g.player(i)
	player.ChainSymbols = append(player.ChainSymbols, s)
}

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
	effect(g *Game, i PlayerIndex)
}

type opponent struct {
	e Effect
}

func (o opponent) effect(g *Game, i PlayerIndex) {
	o.e.effect(g, i.Next())
}

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

func (v VP) effect(g *Game, i PlayerIndex) {
	g.player(i).VP += v
}

// Market ...
type Market interface{}

// PriceMarket ...
// type PriceMarket interface {
// 	Price(r Resource) Money
// }

// OneOfAnyMarket ...
type OneOfAnyMarket []Resource

func (m OneOfAnyMarket) effect(g *Game, i PlayerIndex) {
	player := g.player(i)
	player.OneOfAnyMarkets = append(player.OneOfAnyMarkets, m)
}

// OnePriceMarket ...
type OnePriceMarket struct {
	Res   Resource
	Price Money
}

func (m OnePriceMarket) effect(g *Game, i PlayerIndex) {
	player := g.player(i)
	player.OnePriceMarkets = append(player.OnePriceMarkets, m)
}

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

func (d DiscardMoney) effect(g *Game, i PlayerIndex) {
	g.player(i).Money.Sub(d)
}

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

func (MoneyByCards) effect(g *Game, i PlayerIndex) {
	panic("Not implemented")
}

// MoneyByWonders ...
type MoneyByWonders struct {
	Value Money
}

func (m MoneyByWonders) effect(g *Game, i PlayerIndex) {
	m.Value.Mul(len(g.player(i).Wonders)).effect(g, i)
}

// RepeatTurn ...
func RepeatTurn() Effect {
	return repeatTurn{}
}

type repeatTurn struct{}

func (repeatTurn) effect(g *Game, i PlayerIndex) {
	panic("Not implemented")
}
