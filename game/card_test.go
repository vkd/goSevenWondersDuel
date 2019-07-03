package game

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func findCard(name CardName) *Card {
	c := name.find()
	if c == nil {
		panic(fmt.Sprintf("card %q not found", name))
	}
	return c
}

// from rulebook: Cost of cards
func TestCostOfCards(t *testing.T) {
	var g Game
	assert.True(t, g.checkAndBuyCard(findCard("Lumber yard")))
	assert.Equal(t, Res(Wood, 1), g.player().Resources)

	g = Game{}
	g.applyEffect(Money(1))
	assert.True(t, g.checkAndBuyCard(findCard("Stone pit")))
	assert.Equal(t, Res(Stone, 1), g.player().Resources)

	g = Game{}
	g.applyEffect(Stone)
	assert.True(t, g.checkAndBuyCard(findCard("Baths")))
	assert.Equal(t, VP(3), g.player().VP)
	assert.True(t, g.player().ChainSymbols.Exists(Water))

	g = Game{}
	g.applyEffect(Clay, Stone, Wood)
	assert.True(t, g.checkAndBuyCard(findCard("Arena")))
	assert.Equal(t, VP(3), g.player().VP)

	g = Game{}
	g.applyEffect(Clay, Wood)
	assert.True(t, g.checkAndBuyCard(findCard("Horse breeders")))
	assert.Equal(t, Shields(1), g.war.Shields[g.activePlayer])

	g = Game{}
	g.applyEffect(Horseshoe)
	assert.True(t, g.checkAndBuyCard(findCard("Horse breeders")))
	assert.Equal(t, Shields(1), g.war.Shields[g.activePlayer])
}
