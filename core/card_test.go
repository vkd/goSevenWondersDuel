package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// from rulebook: Cost of cards
func TestCostOfCards_Cost(t *testing.T) {
	g := &Game{}

	// the Lumber Yard is free
	assert.Equal(t, Coins(0), g.Cost("Lumber yard"))

	// the Stone Pit costs 1 coin
	assert.Equal(t, Coins(1), g.Cost("Stone pit"))

	// the Baths require 1 Stone to be built
	g.players[0] = Player{Resources: NewRes(Stone)}
	assert.Equal(t, Coins(0), g.Cost("Baths"))

	// the Arena requires 1 Clay, 1 Stone and 1 Wood
	g.players[0] = Player{Resources: NewRes(Clay, Stone, Wood)}
	assert.Equal(t, Coins(0), g.Cost("Arena"))
}

// from rulebook: Cost of cards
func TestCostOfCards_Chain(t *testing.T) {
	g := &Game{}

	// the construction of the Horse Breeders requires 1 Clay and 1 Wood
	g.players[0] = Player{Resources: NewRes(Clay, Wood)}
	assert.Equal(t, Coins(0), g.Cost("Horse breeders"))

	// OR the possession of the Stable
	g.players[0] = Player{}
	g.apply("Stable")
	assert.Equal(t, Coins(0), g.Cost("Horse breeders"))
}
