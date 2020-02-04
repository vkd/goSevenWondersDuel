package game

import (
	"log"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	numStartPTokens           = 5
	numActiveWondersPerPlayer = 4
	numActiveWonders          = numActiveWondersPerPlayer * numPlayers
)

// Game 7 Wonders
type Game struct {
	players [numPlayers]Player
	active  PlayerIndex

	state                State
	wonderSelectionState WonderSelectionState

	activeWonders, discardWonders WonderNames
	// -----------------------

	// currentPlayerIndex int
	currentEpoch int
	cards        []*Card

	war Military

	victoryType VictoryType
	// victoryPlayer PlayerIndex

	log logrus.Logger
}

// NewGame - create new initialized game instance
func NewGame() *Game {
	var g Game
	g.Initialize()
	return &g
}

// Initialize new game state
func (g *Game) Initialize() {
	if g.state != NewState {
		return
	}

	// take shuffled progress tokens
	ptokens := NewAllPTokenNames().Shuffle()
	g.activePTokens, g.discardPTokens = ptokens[:numStartPTokens], ptokens[numStartPTokens:]

	for i := range g.players {
		g.players[i].Money = 7
	}

	// random first player
	g.active = PlayerIndex(rnd.Int() % numPlayers)

	// g.startAge(ageI, 20)

	// take shuffled wonders
	wnds := NewAllWonderNames().Shuffle()
	g.activeWonders, g.discardWonders = wnds[:numActiveWonders], wnds[numActiveWonders:]

	g.state = WondersSelectionPhase
}

// AvailableProgressTokens - return available to take a progress tokens
func (g *Game) AvailableProgressTokens() []PTokenName {
	var out []PTokenName
	for _, n := range g.activePTokens {
		out = append(out, n)
	}
	return out
}

// AvailableWonders - return available to take a wonders
func (g *Game) AvailableWonders() WonderNames {
	if g.state != WondersSelectionPhase {
		return nil
	}
	var out WonderNames
	for _, n := range g.activeWonders[:4] {
		// if n == choosenWonderName {
		// 	continue
		// }
		out = append(out, n)
	}
	return out
}

const (
	choosenWonderName WonderName = "---"
)

// TakeWonders ...
func (g *Game) TakeWonders(ww ...WonderName) {
	if g.state != WondersSelectionPhase {
		return
	}
	if len(ww) < 1 {
		return
	}
	for _, w := range ww {
		if w == choosenWonderName {
			return
		}
	}

	switch g.wonderSelectionState {
	case Part1Choose1WonderBy1Player, Part2Choose1WonderBy2Player:
		ws := g.activeWonders[:4]
		if !checkChoosenWonders(ws, 0) {
			return
		}
		i := ws.Index(ww[0])
		if i == -1 {
			return
		}

		g.current().Wonders.Append(ww[0])
		g.activeWonders[i] = choosenWonderName

		g.setNextPlayer()
		switch g.wonderSelectionState {
		case Part1Choose1WonderBy1Player:
			g.wonderSelectionState = Part1Choose2WonderBy2Plaeyr
		case Part2Choose1WonderBy2Player:
			g.wonderSelectionState = Part2Choose2WonderBy1Player
		}
	case Part1Choose2WonderBy2Plaeyr, Part2Choose2WonderBy1Player:
		if len(ww) < 2 {
			return
		}
		ws := g.activeWonders[:4]
		if !checkChoosenWonders(ws, 1) {
			return
		}

		i := ws.Index(ww[0])
		if i == -1 {
			return
		}
		j := ws.Index(ww[1])
		if j == -1 || i == j {
			return
		}
		g.current().Wonders.Append(ww[:2]...)
		g.activeWonders[i] = choosenWonderName
		g.activeWonders[j] = choosenWonderName

		var lastIndex int
		for i := range ws {
			if ws[i] != choosenWonderName {
				lastIndex = i
				break
			}
		}
		g.opponent().Wonders.Append(ws[lastIndex])
		g.activeWonders = g.activeWonders[4:]

		switch g.wonderSelectionState {
		case Part1Choose2WonderBy2Plaeyr:
			g.wonderSelectionState = Part2Choose1WonderBy2Player
		case Part2Choose2WonderBy1Player:
			g.state = GameState
		}
	default:
		return
	}
}

func checkChoosenWonders(ws WonderNames, expect int) bool {
	if c := ws.Count(choosenWonderName); c != expect {
		log.Printf("WondersNames wrong: active wonders have not exactly %d choosen name: %d", expect, c)
		return false
	}
	return true
}

// ConstructBuilding - construct one building
func (g *Game) ConstructBuilding(name CardName) {
	panic("Not implemented")
	g.nextTurn()
}

// ConstructWonder by Building card
func (g *Game) ConstructWonder(wname WonderName, cname CardName) {
	panic("Not implemented")
}

// // ActiveIndex of current game state
// func (g *Game) ActiveIndex() PlayerIndex {
// 	return g.activePlayer
// }

// Shields of current game state
func (g *Game) Shields() Shields {
	return g.war.Shields[0] - g.war.Shields[1]
}

// // Desk of active cards
// func (g *Game) Desk() []CardIndex {
// 	panic("Not implemented")
// }

// // Construct a Building
// func (g *Game) Construct(i CardIndex) error {
// 	panic("Not implemented")
// }

// func (g *Game) applyEffect(ee ...Effect) {
// 	g.applyEffectByPlayer(g.activePlayer, ee...)
// }

// func (g *Game) applyEffectByPlayer(pIndex PlayerIndex, ee ...Effect) {
// 	player := g.player(pIndex)
// 	for _, e := range ee {
// 		switch e := e.(type) {
// 		// player
// 		case Money:
// 			player.Money += e
// 		case DiscardMoney:
// 			player.Money.Sub(e)
// 		case Shields:
// 			AddShields(g, e)
// 		case Resource:
// 			player.Resources[e]++
// 		case VP:
// 			player.VP += e

// 		// opponent
// 		case opponent:
// 			g.applyEffectByPlayer(pIndex.Next(), e.e)

// 		// science
// 		case ScientificSymbol:
// 			AddScience(g, pIndex, e)

// 		// chains
// 		case ChainSymbol:
// 			player.ChainSymbols = append(player.ChainSymbols, e)

// 		// markets
// 		case OnePriceMarket:
// 			player.OnePriceMarkets = append(player.OnePriceMarkets, e)
// 		case OneOfAnyMarket:
// 			player.OneOfAnyMarkets = append(player.OneOfAnyMarkets, e)

// 		// money by ...
// 		case MoneyByWonders:
// 			player.Money += e.Value.Mul(len(player.Wonders))
// 		default:
// 			panic(fmt.Sprintf("Not implemented: %T", e))
// 		}
// 	}
// }

func (g *Game) nextTurn() {
	g.setNextPlayer()
}

func (g *Game) victory(playerIndex PlayerIndex, vt VictoryType) {
	g.state = VictoryState
	g.victoryType = vt
	// g.victoryPlayer = playerIndex
}

func (g *Game) canSelectActiveProgressToken() {
	g.state = SelectActiveProgressTokenState
}

// func (g *Game) playerI(i PlayerIndex) *Player {
// 	return &g.players[i]
// }

// // Player ...
// func (g *Game) Player(i PlayerIndex) Player {
// 	return g.players[i]
// }

func (g *Game) player(i PlayerIndex) *Player {
	return &g.players[i]
}

func (g *Game) current() *Player {
	return g.player(g.active)
}

func (g *Game) opponent() *Player {
	return g.player(g.active.Next())
}

func (g *Game) setNextPlayer() {
	g.active = g.active.Next()
}

// func (g *Game) nextPlayerI(i PlayerIndex) *Player {
// 	return g.playerI(i.Next())
// }

// --------------------------------------

var (
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// State ...
type State uint8

// Game states
const (
	NewState State = iota
	WondersSelectionPhase
	GameState
	SelectActiveProgressTokenState
	VictoryState
	// numGameStates int = iota
)

// WonderSelectionState - states of WondersSelectionPhase
type WonderSelectionState = uint8

// Wonder selection states
const (
	Part1Choose1WonderBy1Player WonderSelectionState = iota
	Part1Choose2WonderBy2Plaeyr
	Part2Choose1WonderBy2Player
	Part2Choose2WonderBy1Player
)

// func (g *Game) startAge(age []Card, num int) {
// 	var cards = make([]CardIndex, len(age))
// 	for i := range cards {
// 		cards[i] = CardIndex(i)
// 	}
// 	rnd.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })

// 	// first num cards
// 	g.cards = make([]*Card, num)
// 	for i := 0; i < num; i++ {
// 		var c = ageI[cards[i]]
// 		g.cards[i] = &c
// 	}
// }

// Cards ...
// func (g *Game) Cards() []CardIndex {
// 	out := make([]CardIndex, len(g.cards))
// 	for i := range out {
// 		out[i] = g.cards[i].Index
// 	}

// 	return out
// }

// ConstructByIndex ...
func (g *Game) ConstructByIndex(index int) bool {
	c := g.getCardByIndex(index)
	if c == nil {
		log.Printf("Construct is denied: card not found")
		return false
	}

	ok := g.checkAndBuyCard(c)
	if !ok {
		log.Printf("Construct is failed: cannot check card price")
		return false
	}

	g.cards[index] = nil
	return true
}

func (g *Game) getCardByIndex(index int) *Card {
	if index < 0 || index >= len(g.cards) {
		log.Printf("Construct is deny: wrong card index")
		return nil
	}

	return g.cards[index]
}

// func (g *Game) getIncome(effects ...Effect) {
// 	g.getIncomeByPlayerIndex(g.activePlayer, effects...)
// }

// func (g *Game) getIncomeByPlayerIndex(index PlayerIndex, effects ...Effect) {
// 	player := &g.players[index]

// 	for _, e := range effects {
// 		switch e := e.(type) {
// 		case Money:
// 			player.Money += e
// 			// case Shields:
// 		case Resource:
// 			player.Resources[e]++
// 		case VP:
// 			player.VP += e
// 		case OnePriceMarket:
// 			player.OnePriceMarkets = append(player.OnePriceMarkets, e)
// 		case OneOfAnyMarket:
// 			player.OneOfAnyMarkets = append(player.OneOfAnyMarkets, e)
// 		case ChainSymbol:
// 			player.ChainSymbols = append(player.ChainSymbols, e)
// 		default:
// 			panic(fmt.Sprintf("Not implemented: %T", e))
// 		}
// 	}
// 	return
// }

// func (g *Game) currentPlayer() *Player {
// 	return g.player()
// }

// func (g *Game) oppositePlayer() *Player {
// 	return g.opponent()
// }

// GetCardCostByIndex ...
// func (g *Game) GetCardCostByIndex(index int) (Money, bool) {
// 	c := g.getCardByIndex(index)
// 	if c == nil {
// 		log.Printf("Get card cost is denied: card not found")
// 		return 0, false
// 	}

// 	return g.costCardOfMoney(c, g.activePlayer), true
// }

func (g *Game) costCardOfMoney(c *Card, pi PlayerIndex) Money {
	player := g.player(pi)
	if c.FreeCostChainSymbol.Exists && player.ChainSymbols.Exists(c.FreeCostChainSymbol.ChainSymbol) {
		return 0
	}

	tradingCosts := g.getTradingCosts(pi)
	debug("Trading costs: %v", tradingCosts)
	checkMoney := CostByMoney(player, c.Cost, tradingCosts)
	return checkMoney
}

func (g *Game) checkAndBuyCard(c *Card) bool {
	checkMoney := g.costCardOfMoney(c, g.active)

	player := g.current()
	if checkMoney > player.Money {
		log.Printf("Error on checking of a trading: not enough money")
		return false
	}

	player.Money -= checkMoney
	for _, e := range c.Effects {
		e.effect(g, g.active)
	}
	return true
}

func (g *Game) getTradingCosts(pi PlayerIndex) (out TradingCosts) {
	opponent := g.player(pi.Next())
	tc := NewTradingCosts(opponent.Resources)
	return tc.ApplyMarkets(g.player(pi).OnePriceMarkets)
}
