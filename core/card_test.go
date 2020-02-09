package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// from rulebook: Cost of cards
func TestCostOfCards_Cost(t *testing.T) {
	g := &Game{}

	// the Lumber Yard is free
	assert.Equal(t, Coins(0), g.CostName("Lumber yard"))

	// the Stone Pit costs 1 coin
	assert.Equal(t, Coins(1), g.CostName("Stone pit"))

	// the Baths require 1 Stone to be built
	g.players[0] = Player{Resources: NewRes(Stone)}
	assert.Equal(t, Coins(0), g.CostName("Baths"))

	// the Arena requires 1 Clay, 1 Stone and 1 Wood
	g.players[0] = Player{Resources: NewRes(Clay, Stone, Wood)}
	assert.Equal(t, Coins(0), g.CostName("Arena"))
}

// from rulebook: Cost of cards
func TestCostOfCards_Chain(t *testing.T) {
	g := &Game{}

	// the construction of the Horse Breeders requires 1 Clay and 1 Wood
	g.players[0] = Player{Resources: NewRes(Clay, Wood)}
	assert.Equal(t, Coins(0), g.CostName("Horse breeders"))

	// OR the possession of the Stable
	g.players[0] = Player{}
	g.apply("Stable")
	assert.Equal(t, Coins(0), g.CostName("Horse breeders"))
}
