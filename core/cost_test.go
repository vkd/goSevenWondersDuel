package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// from rulebook: Free Construction
func TestCost_FreeConstruction(t *testing.T) {
	c := CardName("Lumber yard").card().Cost
	m := CostByCoins(c, Player{}, TradingPrice{})
	assert.Equal(t, Coins(0), m)
}

// from rulebook: Production
func TestCost_Production(t *testing.T) {
	p := Player{}
	p.Resources[Stone] = 1
	p.Resources[Clay] = 3
	p.Resources[Papyrus] = 1

	tc := NewTradingPrice(Player{})

	c := CardName("Baths").card().Cost
	assert.Equal(t, Coins(0), CostByCoins(c, p, tc))

	c = CardName("Garrison").card().Cost
	assert.Equal(t, Coins(0), CostByCoins(c, p, tc))

	c = CardName("Apothecary").card().Cost
	assert.Equal(t, Coins(2), CostByCoins(c, p, tc))
}

// from rulebook: Trading
func TestCost_TradingCosts(t *testing.T) {
	var bruno, antoine Player

	bruno.Resources[Stone] = 2

	// for Antoine
	costOneStone := NewTradingPrice(bruno)[Stone]
	assert.Equal(t, Coins(4), costOneStone)

	// for Bruno
	costOneStone = NewTradingPrice(antoine)[Stone]
	assert.Equal(t, Coins(2), costOneStone)
}

// from rulebook: Trading
func TestCost_Trading(t *testing.T) {
	var bruno, antoine Player

	bruno.Resources[Stone] = 2
	antoine.Resources[Clay] = 1

	c := CardName("Fortifications").card().Cost
	tc := NewTradingPrice(antoine)
	m := CostByCoins(c, bruno, tc)
	assert.Equal(t, Coins(5), m)

	c = CardName("Aqueduct").card().Cost
	tc = NewTradingPrice(bruno)
	m = CostByCoins(c, antoine, tc)
	assert.Equal(t, Coins(12), m)
}
