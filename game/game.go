package game

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

const (
	numStartPTokens           = 5
	numActiveWondersPerPlayer = 4
	numActiveWonders          = numActiveWondersPerPlayer * numPlayers
)

// Game 7 Wonders
type Game struct {
	players      [numPlayers]Player
	activePlayer PlayerIndex

	state                State
	wonderSelectionState WonderSelectionState

	activePTokens  []PTokenName
	discardPTokens []PTokenName

	activeWonders, discardWonders WonderNames
	// -----------------------

	// currentPlayerIndex int
	currentEpoch int
	cards        []*Card

	war Military

	victoryType VictoryType
	// victoryPlayer PlayerIndex

	log *log.Logger
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
	g.activePlayer = PlayerIndex(rnd.Int() % numPlayers)

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
		if n == choosenWonderName {
			continue
		}
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
	case Part1ChooseBy1Player1Wonder, Part2ChooseBy2Player1Wonder:
		ws := g.activeWonders[:4]
		if !checkChoosenWonders(ws, 0) {
			return
		}
		i := ws.Index(ww[0])
		if i == -1 {
			return
		}

		g.player().Wonders.Append(ww[0])
		g.activeWonders[i] = choosenWonderName

		g.setNextPlayer()
		switch g.wonderSelectionState {
		case Part1ChooseBy1Player1Wonder:
			g.wonderSelectionState = Part1ChooseBy2Plaeyr2Wonder
		case Part2ChooseBy2Player1Wonder:
			g.wonderSelectionState = Part2ChooseBy1Player2Wonder
		}
	case Part1ChooseBy2Plaeyr2Wonder, Part2ChooseBy1Player2Wonder:
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
		g.player().Wonders.Append(ww[:2]...)
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
		case Part1ChooseBy2Plaeyr2Wonder:
			g.wonderSelectionState = Part2ChooseBy2Player1Wonder
		case Part2ChooseBy1Player2Wonder:
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

// ActiveIndex of current game state
func (g *Game) ActiveIndex() PlayerIndex {
	return g.activePlayer
}

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

func (g *Game) applyEffect(ee ...Effect) {
	g.applyEffectByPlayer(g.activePlayer, ee...)
}

func (g *Game) applyEffectByPlayer(pIndex PlayerIndex, ee ...Effect) {
	player := g.playerI(pIndex)
	for _, e := range ee {
		switch e := e.(type) {
		// player
		case Money:
			player.Money += e
		case DiscardMoney:
			player.Money.Sub(e)
		case Shields:
			AddShields(g, e)
		case Resource:
			player.Resources[e]++
		case VP:
			player.VP += e

		// opponent
		case opponent:
			g.applyEffectByPlayer(pIndex.Next(), e.e)

		// science
		case ScientificSymbol:
			AddScience(g, e)

		// chains
		case ChainSymbol:
			player.ChainSymbols = append(player.ChainSymbols, e)

		// markets
		case OnePriceMarket:
			player.OnePriceMarkets = append(player.OnePriceMarkets, e)
		case OneOfAnyMarket:
			player.OneOfAnyMarkets = append(player.OneOfAnyMarkets, e)

		// money by ...
		case MoneyByWonders:
			player.Money += e.Value.Mul(len(player.Wonders))
		default:
			panic(fmt.Sprintf("Not implemented: %T", e))
		}
	}
}

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

func (g *Game) playerI(i PlayerIndex) *Player {
	return &g.players[i]
}

func (g *Game) player() *Player {
	return g.playerI(g.activePlayer)
}

func (g *Game) opponent() *Player {
	return g.nextPlayerI(g.activePlayer)
}

func (g *Game) setNextPlayer() {
	g.activePlayer = g.activePlayer.Next()
}

func (g *Game) nextPlayerI(i PlayerIndex) *Player {
	return g.playerI(i.Next())
}

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
	Part1ChooseBy1Player1Wonder WonderSelectionState = iota
	Part1ChooseBy2Plaeyr2Wonder
	Part2ChooseBy2Player1Wonder
	Part2ChooseBy1Player2Wonder
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

func (g *Game) currentPlayer() *Player {
	return g.player()
}

func (g *Game) oppositePlayer() *Player {
	return g.opponent()
}

// GetCardCostByIndex ...
func (g *Game) GetCardCostByIndex(index int) (Money, bool) {
	c := g.getCardByIndex(index)
	if c == nil {
		log.Printf("Get card cost is denied: card not found")
		return 0, false
	}

	return g.costCardOfMoney(c, g.activePlayer), true
}

func (g *Game) costCardOfMoney(c *Card, pi PlayerIndex) Money {
	player := g.playerI(pi)
	if c.FreeCostChainSymbol.Exists && player.ChainSymbols.Exists(c.FreeCostChainSymbol.ChainSymbol) {
		return 0
	}

	tradingCosts := g.getTradingCosts(pi)
	debug("Trading costs: %v", tradingCosts)
	checkMoney := CostByMoney(player, c.Cost, tradingCosts)
	return checkMoney
}

func (g *Game) checkAndBuyCard(c *Card) bool {
	checkMoney := g.costCardOfMoney(c, g.activePlayer)

	player := g.currentPlayer()
	if checkMoney > player.Money {
		log.Printf("Error on checking of a trading: not enough money")
		return false
	}

	player.Money -= checkMoney
	g.applyEffect(c.Effects...)
	return true
}

func (g *Game) getTradingCosts(pi PlayerIndex) (out TradingCosts) {
	opponent := g.nextPlayerI(pi)
	tc := NewTradingCosts(opponent.Resources)
	return tc.ApplyMarkets(g.playerI(pi).OnePriceMarkets)
}
