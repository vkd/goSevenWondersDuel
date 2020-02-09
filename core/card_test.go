package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// from rulebook: Cost of cards
func TestCostOfCards_Cost(t *testing.T) {
	g := &Game{}

	// the Lumber Yard is free
	cost, ok := g.CostName("Lumber yard")
	assert.True(t, ok)
	assert.Equal(t, Coins(0), cost)

	// the Stone Pit costs 1 coin
	cost, ok = g.CostName("Stone pit")
	assert.True(t, ok)
	assert.Equal(t, Coins(1), cost)

	// the Baths require 1 Stone to be built
	g.players[0] = Player{Resources: NewRes(Stone)}
	cost, ok = g.CostName("Baths")
	assert.True(t, ok)
	assert.Equal(t, Coins(0), cost)

	// the Arena requires 1 Clay, 1 Stone and 1 Wood
	g.players[0] = Player{Resources: NewRes(Clay, Stone, Wood)}
	cost, ok = g.CostName("Arena")
	assert.True(t, ok)
	assert.Equal(t, Coins(0), cost)
}

// from rulebook: Cost of cards
func TestCostOfCards_Chain(t *testing.T) {
	g := &Game{}

	// the construction of the Horse Breeders requires 1 Clay and 1 Wood
	g.players[0] = Player{Resources: NewRes(Clay, Wood)}
	cost, ok := g.CostName("Horse breeders")
	assert.True(t, ok)
	assert.Equal(t, Coins(0), cost)

	// OR the possession of the Stable
	g.players[0] = Player{}
	g.apply("Stable")
	cost, ok = g.CostName("Horse breeders")
	assert.True(t, ok)
	assert.Equal(t, Coins(0), cost)
}
