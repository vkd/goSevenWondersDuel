package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestZeroGame(t *testing.T) {
	game, err := NewGame(WithSeed(0))
	require.NoError(t, err)

	wonders := game.GetAvailableWonders()

	// [0:The Appian Way 1:The Statue of Zeus 2:The Great Library 3:Temple of Artemis
	// 4:The Hanging Gardens 5:The Great Lighthouse 6:The Mausoleum 7:The Sphinx]
	assert.Len(t, wonders, initialWonders)

	err = game.SelectWonders(
		// Temple of Artemis, The Great Library, The Hanging Gardens, The Sphinx
		[...]WonderID{wonders[3], wonders[2], wonders[4], wonders[7]},
		// The Appian Way, The Statue of Zeus, The Great Lighthouse, The Mausoleum
		[...]WonderID{wonders[0], wonders[1], wonders[5], wonders[6]},
	)
	require.NoError(t, err)
	assert.Equal(t, game.GetState(), StateGameTurn)

	// 0
	_, err = game.ConstructBuilding(cardID("Logging camp"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Quarry"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Workshop"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Guard tower"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Stone pit"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Apothecary"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Clay pool"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Baths"))
	require.NoError(t, err)

	// 0
	_, err = game.DiscardCard(cardID("Pharmacist"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Garrison"))
	require.NoError(t, err)

	// 0
	_, err = game.DiscardCard(cardID("Clay pit"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Press"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Glassworks"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Lumber yard"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Palisade"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Stable"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructWonder(cardID("Scriptorium"), wonderID("Temple of Artemis"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Wood reserve"))
	require.NoError(t, err)

	_, err = game.ConstructBuilding(cardID("Theater"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructWonder(cardID("Altar"), wonderID("The Great Library"))
	require.NoError(t, err)
	err = game.PlayDiscardedPToken(pTokenID("Agriculture"))
	require.NoError(t, err)

	// === End Age I ===
	assert.Equal(t, [numPlayers]Shields{1, 3}, game.Military().Shields)
	assert.Equal(t, 2, len(game.discardedCards))
	assert.Equal(t, 2, len(game.buildWonders[0]))
	assert.Equal(t, 0, len(game.buildWonders[1]))
	assert.Equal(t, 7, countBuiltCards(game, 0))
	assert.Equal(t, 9, countBuiltCards(game, 1))
	assert.Equal(t, VP(9), countVPs(game, 0))
	assert.Equal(t, VP(7), countVPs(game, 1))
	assert.Equal(t, Coins(11), game.Player(0).Coins)
	assert.Equal(t, Coins(1), game.Player(1).Coins)
	assert.Equal(t, uint8(1), game.currentAge)
	assert.Equal(t, game.GetState(), StateChooseFirstPlayer)

	// === Age II ===
	err = game.ChooseFirstPlayer(1)
	require.NoError(t, err)
	assert.Equal(t, game.GetState(), StateGameTurn)

	// 1
	_, err = game.ConstructBuilding(cardID("Horse breeders"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructWonder(cardID("Brick yard"), wonderID("The Hanging Gardens"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Drying room"))
	require.NoError(t, err)

	_, err = game.DiscardCard(cardID("Tribunal"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructWonder(cardID("Library"), wonderID("The Sphinx"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Forum"))
	require.NoError(t, err)

	_, err = game.ConstructWonder(cardID("Laboratory"), wonderID("The Great Lighthouse"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Caravansery"))
	require.NoError(t, err)
	_, err = game.DiscardCard(cardID("Barracks"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Postrum"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Aqueduct"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Brewery"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Shelf quarry"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("School"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Walls"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Dispensary"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Archery range"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Parade ground"))
	require.NoError(t, err)
	_, err = game.DiscardCard(cardID("Customs house"))
	require.NoError(t, err)

	// 0
	_, err = game.DiscardCard(cardID("Sawmill"))
	require.NoError(t, err)

	// === End Age II ===
	assert.Equal(t, [numPlayers]Shields{1 + 2, 3 + 5}, game.Military().Shields)
	assert.Equal(t, 2+4, len(game.discardedCards))
	assert.Equal(t, 2+2, len(game.buildWonders[0]))
	assert.Equal(t, 0+1, len(game.buildWonders[1]))
	assert.Equal(t, 7+8, countBuiltCards(game, 0))
	assert.Equal(t, 9+5, countBuiltCards(game, 1))
	assert.Equal(t, VP(9+16), countVPs(game, 0))
	assert.Equal(t, VP(7+9), countVPs(game, 1))
	assert.Equal(t, Coins(11), game.Player(0).Coins)
	assert.Equal(t, Coins(2), game.Player(1).Coins)
	assert.Equal(t, uint8(2), game.currentAge)
	assert.Equal(t, game.GetState(), StateChooseFirstPlayer)

	// === Age III ===
	err = game.ChooseFirstPlayer(1)
	require.NoError(t, err)
	assert.Equal(t, game.GetState(), StateGameTurn)

	// 1
	_, err = game.DiscardCard(cardID("Shipowners guild"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("University"))
	require.NoError(t, err)

	_, err = game.ConstructWonder(cardID("Circus"), wonderID("The Statue of Zeus"))
	require.NoError(t, err)
	err = game.DiscardOpponentBuild(cardID("Clay pool"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Port"))
	require.NoError(t, err)
	_, err = game.DiscardCard(cardID("Academy"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Armory"))
	require.NoError(t, err)

	_, err = game.ConstructWonder(cardID("Arsenal"), wonderID("The Appian Way"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Builders guild"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Pantheon"))
	require.NoError(t, err)
	_, err = game.DiscardCard(cardID("Observatory"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Senate"))
	require.NoError(t, err)
	_, err = game.DiscardCard(cardID("Obelisk"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Arena"))
	require.NoError(t, err)
	_, err = game.DiscardCard(cardID("Chamber of commerce"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Fortifications"))
	require.NoError(t, err)
	_, err = game.DiscardCard(cardID("Courthouse"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Siege workshop"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Gardens"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Tacticians guild"))
	require.NoError(t, err)
	_, err = game.DiscardCard(cardID("Lighthouse"))
	require.NoError(t, err)

	// === End Age III ===
	assert.Equal(t, game.GetState(), StateVictory)
	assert.Equal(t, Winner1Player, game.winner)
	assert.Equal(t, WinCivilian, game.winReason)

	assert.Equal(t, [numPlayers]Shields{1 + 2 + 4, 3 + 5 + 1}, game.Military().Shields)
	assert.Equal(t, 2+4+8, len(game.discardedCards))
	assert.Equal(t, 2+2, len(game.buildWonders[0]))
	assert.Equal(t, 0+1+2, len(game.buildWonders[1]))
	assert.Equal(t, 7+8+8, countBuiltCards(game, 0))
	assert.Equal(t, 9+5+2, countBuiltCards(game, 1))
	assert.Equal(t, VP(9+16+36), countVPs(game, 0))
	assert.Equal(t, VP(7+9+24), countVPs(game, 1))
	assert.Equal(t, Coins(25), game.Player(0).Coins)
	assert.Equal(t, Coins(7), game.Player(1).Coins)

	assert.Equal(t, [numVPTypes]VP{
		BlueVP:     15,
		GreenVP:    6,
		YellowVP:   9,
		PurpleVP:   6,
		WonderVP:   13,
		PTokenVP:   4,
		CoinsVP:    8,
		MilitaryVP: 0,
	}, game.vps[0])
	assert.Equal(t, [numVPTypes]VP{
		BlueVP:     17,
		GreenVP:    1,
		YellowVP:   0,
		PurpleVP:   8,
		WonderVP:   10,
		PTokenVP:   0,
		CoinsVP:    2,
		MilitaryVP: 2,
	}, game.vps[1])
}

func countBuiltCards(g *Game, i PlayerIndex) (count int) {
	for _, cc := range g.builtCards[i] {
		count += len(cc)
	}
	return
}

func countVPs(g *Game, i PlayerIndex) (count VP) {
	for _, vps := range g.vps[i] {
		count += vps
	}
	return
}
