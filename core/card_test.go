package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func costOf(name CardName) Cost {
	return name.card().Cost
}

// from rulebook: Cost of cards
func TestCostOfCards_Cost(t *testing.T) {
	var tp = NewTradingPrice(Player{})
	var p = Player{}
	var c Cost

	// the Lumber Yard is free
	c = costOf("Lumber yard")
	assert.Equal(t, Coins(0), CostByCoins(c, p, tp))

	// the Stone Pit costs 1 coin
	c = costOf("Stone pit")
	assert.Equal(t, Coins(1), CostByCoins(c, p, tp))

	// the Baths require 1 Stone to be built
	p = Player{Resources: NewRes(Stone)}
	c = costOf("Baths")
	assert.Equal(t, Coins(0), CostByCoins(c, p, tp))

	// the Arena requires 1 Clay, 1 Stone and 1 Wood
	p = Player{Resources: NewRes(Clay, Stone, Wood)}
	c = costOf("Arena")
	assert.Equal(t, Coins(0), CostByCoins(c, p, tp))
}

// from rulebook: Cost of cards
func TestCostOfCards_Chain(t *testing.T) {
	g := &Game{}
	var tp = NewTradingPrice(Player{})

	// the construction of the Horse Breeders requires 1 Clay and 1 Wood
	g.players[0] = Player{Resources: NewRes(Clay, Wood)}
	assert.Equal(t, Coins(0), CostByCoins(costOf("Horse breeders"), g.players[0], tp))

	// OR the possession of the Stable
	g.players[0] = Player{}
	g.apply("Stable")
	assert.Equal(t, Coins(0), g.CardCostCoins(getID("Horse breeders")))
}

func getID(name CardName) CardID {
	for i, c := range cards {
		if c.Name == name {
			return CardID(i)
		}
	}
	panic("unknown name")
}
