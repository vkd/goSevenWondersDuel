package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// from rulebook: Preparation
func TestGame_Preparation(t *testing.T) {
	g := NewGame()
	assert.Len(t, g.AvailableProgressTokens(), 5)

	assert.Equal(t, Money(7), g.players[0].Money)
	assert.Equal(t, Money(7), g.players[1].Money)
}

// from rulebook: Wonders Selection Phase
func TestGame_WondersSelectionPhase(t *testing.T) {
	g := NewGame()
	assert.Equal(t, WondersSelectionPhase, g.state)
	startPlayer := g.ActiveIndex()

	// part 1 choose by 1 player 1 wonder
	// -----------------------------------
	ws := g.AvailableWonders()
	assert.Len(t, ws, 4)

	// control
	var fPlayer WonderNames
	fPlayer.Append(ws[2], ws[1])

	g.TakeWonders(ws[2])

	// part 1 choose by 2 player 2 wonders
	// -----------------------------------
	ws3 := g.AvailableWonders()
	assert.Len(t, ws3, 3)
	assert.Equal(t, ws[0], ws3[0])
	assert.Equal(t, ws[1], ws3[1])
	assert.Equal(t, ws[3], ws3[2])

	// control
	var sPlayer WonderNames
	sPlayer.Append(ws[0], ws[3])

	g.TakeWonders(ws[0], ws[3])

	// part 2 choose by 2 player 1 wonder
	// -----------------------------------
	ws = g.AvailableWonders()
	assert.Len(t, ws, 4)

	// control
	sPlayer.Append(ws[1], ws[0])

	g.TakeWonders(ws[1])

	// part 2 choose by 1 player 2 wonders
	// -----------------------------------
	ws3 = g.AvailableWonders()
	assert.Len(t, ws3, 3)
	assert.Equal(t, ws[0], ws3[0])
	assert.Equal(t, ws[2], ws3[1])
	assert.Equal(t, ws[3], ws3[2])

	// control
	fPlayer.Append(ws[2], ws[3])

	g.TakeWonders(ws[2], ws[3])

	// game state
	// -----------------------------------
	assert.Equal(t, GameState, g.state)
	assert.Equal(t, startPlayer, g.ActiveIndex())
	assert.Len(t, g.activeWonders, 0)

	assert.Equal(t, fPlayer, g.player().Wonders)
	assert.Equal(t, sPlayer, g.opponent().Wonders)
}

func TestGame_Military(t *testing.T) {
	var g Game
	g.applyEffect(Shields(1))
	assert.Equal(t, Shields(1), g.war.Shields[g.activePlayer])

	g = Game{}
	g.applyEffect(Money(7), Opponent(Money(7)), Shields(3))
	assert.Equal(t, Money(7), g.player().Money)
	assert.Equal(t, Money(5), g.opponent().Money)
	assert.Equal(t, Shields(3), g.war.Shields[g.activePlayer])

	g = Game{}
	g.applyEffect(Money(7), Opponent(Money(7)), Shields(6))
	assert.Equal(t, Money(7), g.player().Money)
	assert.Equal(t, Money(0), g.opponent().Money)

	g = Game{}
	g.applyEffect(Money(7), Opponent(Money(7)), Shields(9))
	assert.Equal(t, Money(7), g.player().Money)
	assert.Equal(t, Money(0), g.opponent().Money)
	assert.Equal(t, VictoryState, g.state)
}

func TestGameGetCardCostByIndex(t *testing.T) {
	addCard := func(g *Game, c Card) { g.cards = append(g.cards, &c) }

	g := &Game{}
	addCard(g, newCard("Test", Red, Cost(Wood, Stone)))

	money, ok := g.GetCardCostByIndex(0)
	assert.True(t, ok)
	assert.Equal(t, Money(4), money)

	g = &Game{}
	addCard(g, newCard("Test", Red, Cost(Wood, Stone)))
	g.applyEffect(OnePriceMarket{Stone, 1})

	money, ok = g.GetCardCostByIndex(0)
	assert.True(t, ok)
	assert.Equal(t, Money(3), money)

	g = &Game{}
	g.applyEffect(Wood)
	g.nextTurn()
	addCard(g, newCard("Text", Red, Cost(Wood, Stone))) // Cost: 3 2 2 2 2

	money, ok = g.GetCardCostByIndex(0)
	assert.True(t, ok)
	assert.Equal(t, Money(5), money)

	g = &Game{}
	g.applyEffect(Wood)
	addCard(g, newCard("Text", Red, Cost(Wood, Stone)))

	money, ok = g.GetCardCostByIndex(0)
	assert.True(t, ok)
	assert.Equal(t, Money(2), money)

	g = &Game{}
	g.applyEffect(Wood)
	g.nextTurn()
	g.applyEffect(OnePriceMarket{Clay, 1})
	addCard(g, newCard("Text", Red, Cost(Wood, Stone, Clay))) // costs: 3 2 1

	money, ok = g.GetCardCostByIndex(0)
	assert.True(t, ok)
	assert.Equal(t, Money(6), money)

	g.applyEffect(OneOfAnyMarket([]Resource{Wood, Stone}))
	g.applyEffect(OneOfAnyMarket([]Resource{Wood, Clay}))

	money, ok = g.GetCardCostByIndex(0)
	assert.True(t, ok)
	assert.Equal(t, Money(1), money)
}
