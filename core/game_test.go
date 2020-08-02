package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestZeroGame(t *testing.T) { //nolint: funlen
	game, err := NewGame(WithSeed(0))
	require.NoError(t, err)

	err = game.SelectWonders(
		[...]WonderID{wonderID("Temple of Artemis"), wonderID("The Great Library"), wonderID("The Hanging Gardens"), wonderID("The Sphinx")},
		//
		[...]WonderID{wonderID("The Appian Way"), wonderID("The Statue of Zeus"), wonderID("The Great Lighthouse"), wonderID("The Mausoleum")},
	)
	require.NoError(t, err)

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
	assert.Equal(t, [numPlayers]Shields{1, 3}, game.Military().ConflictPawn.Shields)
	assert.Equal(t, 2, len(game.discardedCards))
	assert.Equal(t, 2, game.WondersState.CountBuiltByPlayer(0))
	assert.Equal(t, 0, game.WondersState.CountBuiltByPlayer(1))
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
	assert.Equal(t, [numPlayers]Shields{1 + 2, 3 + 5}, game.Military().ConflictPawn.Shields)
	assert.Equal(t, 2+4, len(game.discardedCards))
	assert.Equal(t, 2+2, game.WondersState.CountBuiltByPlayer(0))
	assert.Equal(t, 0+1, game.WondersState.CountBuiltByPlayer(1))
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
	assert.Equal(t, CivilianVictory, game.victoryType)

	assert.Equal(t, [numPlayers]Shields{1 + 2 + 4, 3 + 5 + 1}, game.Military().ConflictPawn.Shields)
	assert.Equal(t, 2+4+8, len(game.discardedCards))
	assert.Equal(t, 2+2, game.WondersState.CountBuiltByPlayer(0))
	assert.Equal(t, 0+1+2, game.WondersState.CountBuiltByPlayer(1))
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
		SumVP:      61,
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
		SumVP:      40,
	}, game.vps[1])
}

func TestZeroGame_MilitarySupremacy(t *testing.T) { //nolint: funlen
	game, err := NewGame(WithSeed(0))
	require.NoError(t, err)

	err = game.SelectWonders(
		[...]WonderID{wonderID("Temple of Artemis"), wonderID("The Great Library"), wonderID("The Hanging Gardens"), wonderID("The Sphinx")},
		//
		[...]WonderID{wonderID("The Appian Way"), wonderID("The Statue of Zeus"), wonderID("The Great Lighthouse"), wonderID("The Mausoleum")},
	)
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Quarry"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Logging camp"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Guard tower"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Pharmacist"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructWonder(cardID("Apothecary"), wonderID("Temple of Artemis"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Clay pool"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Press"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructWonder(cardID("Baths"), wonderID("The Sphinx"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Garrison"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Clay pit"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Lumber yard"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Workshop"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Stone pit"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Glassworks"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Palisade"))
	require.NoError(t, err)
	_, err = game.DiscardCard(cardID("Theater"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Stable"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Scriptorium"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Wood reserve"))
	require.NoError(t, err)
	_, err = game.DiscardCard(cardID("Altar"))
	require.NoError(t, err)

	// === End Age I ===
	assert.Equal(t, [numPlayers]Shields{4, 0}, game.Military().ConflictPawn.Shields)
	assert.Equal(t, 2, len(game.discardedCards))
	assert.Equal(t, 2, game.WondersState.CountBuiltByPlayer(0))
	assert.Equal(t, 0, game.WondersState.CountBuiltByPlayer(1))
	assert.Equal(t, 9, countBuiltCards(game, 0))
	assert.Equal(t, 7, countBuiltCards(game, 1))
	assert.Equal(t, VP(6), countVPs(game, 0))
	assert.Equal(t, VP(1), countVPs(game, 1))
	assert.Equal(t, Coins(2), game.Player(0).Coins)
	assert.Equal(t, Coins(2), game.Player(1).Coins)
	assert.Equal(t, uint8(1), game.currentAge)
	assert.Equal(t, game.GetState(), StateChooseFirstPlayer)

	// === Age II ===
	err = game.ChooseFirstPlayer(1)
	require.NoError(t, err)
	assert.Equal(t, game.GetState(), StateGameTurn)

	// 1
	_, err = game.DiscardCard(cardID("Brick yard"))
	require.NoError(t, err)

	// 0
	_, err = game.DiscardCard(cardID("Horse breeders"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Library"))
	require.NoError(t, err)
	err = game.ChoosePToken(pTokenID("Law"))
	require.NoError(t, err)

	// 0
	_, err = game.DiscardCard(cardID("Tribunal"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Forum"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Barracks"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("School"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Drying room"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Laboratory"))
	require.NoError(t, err)
	err = game.ChoosePToken(pTokenID("Philosophy"))
	require.NoError(t, err)

	// 0
	_, err = game.DiscardCard(cardID("Postrum"))
	require.NoError(t, err)
	_, err = game.DiscardCard(cardID("Aqueduct"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Brewery"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Dispensary"))
	require.NoError(t, err)
	err = game.ChoosePToken(pTokenID("Architecture"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Archery range"))
	require.NoError(t, err)
	_, err = game.ConstructWonder(cardID("Shelf quarry"), wonderID("The Great Lighthouse"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Caravansery"))
	require.NoError(t, err)
	_, err = game.DiscardCard(cardID("Sawmill"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Walls"))
	require.NoError(t, err)

	// === End Age II ===
	assert.Equal(t, game.GetState(), StateVictory)
	assert.Equal(t, Winner1Player, game.winner)
	assert.Equal(t, MilitarySupremacy, game.victoryType)

	assert.Equal(t, [numPlayers]Shields{4 + 5, 0}, game.Military().ConflictPawn.Shields)
	assert.Equal(t, 2+6, len(game.discardedCards))
	assert.Equal(t, 2, game.WondersState.CountBuiltByPlayer(0))
	assert.Equal(t, 0+1, game.WondersState.CountBuiltByPlayer(1))
	assert.Equal(t, 9+6, countBuiltCards(game, 0))
	assert.Equal(t, 7+5, countBuiltCards(game, 1))

	assert.Equal(t, VP(6+14), countVPs(game, 0))
	assert.Equal(t, VP(1+18), countVPs(game, 1))
	assert.Equal(t, Coins(12), game.Player(0).Coins)
	assert.Equal(t, Coins(3), game.Player(1).Coins)

	assert.Equal(t, [numVPTypes]VP{
		BlueVP:     0,
		GreenVP:    0,
		YellowVP:   0,
		PurpleVP:   0,
		WonderVP:   6,
		PTokenVP:   0,
		CoinsVP:    4,
		MilitaryVP: 10,
		SumVP:      20,
	}, game.vps[0])
	assert.Equal(t, [numVPTypes]VP{
		BlueVP:     0,
		GreenVP:    7,
		YellowVP:   0,
		PurpleVP:   0,
		WonderVP:   4,
		PTokenVP:   7,
		CoinsVP:    1,
		MilitaryVP: 0,
		SumVP:      19,
	}, game.vps[1])
}

func TestZeroGame_ScientificSupremacy(t *testing.T) { //nolint: funlen
	game, err := NewGame(WithSeed(0))
	require.NoError(t, err)

	err = game.SelectWonders(
		[...]WonderID{wonderID("Temple of Artemis"), wonderID("The Great Library"), wonderID("The Hanging Gardens"), wonderID("The Sphinx")},
		//
		[...]WonderID{wonderID("The Appian Way"), wonderID("The Statue of Zeus"), wonderID("The Great Lighthouse"), wonderID("The Mausoleum")},
	)
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Quarry"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Logging camp"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Guard tower"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Pharmacist"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructWonder(cardID("Baths"), wonderID("Temple of Artemis"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Garrison"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Apothecary"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Clay pool"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Press"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Clay pit"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Workshop"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Lumber yard"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Glassworks"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Stable"))
	require.NoError(t, err)
	_, err = game.DiscardCard(cardID("Theater"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Stone pit"))
	require.NoError(t, err)
	_, err = game.DiscardCard(cardID("Palisade"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Altar"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Scriptorium"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Wood reserve"))
	require.NoError(t, err)

	// === End Age I ===
	assert.Equal(t, [numPlayers]Shields{3, 0}, game.Military().ConflictPawn.Shields)
	assert.Equal(t, 2, len(game.discardedCards))
	assert.Equal(t, 1, game.WondersState.CountBuiltByPlayer(0))
	assert.Equal(t, 0, game.WondersState.CountBuiltByPlayer(1))
	assert.Equal(t, 10, countBuiltCards(game, 0))
	assert.Equal(t, 7, countBuiltCards(game, 1))
	assert.Equal(t, VP(3), countVPs(game, 0))
	assert.Equal(t, VP(2), countVPs(game, 1))
	assert.Equal(t, Coins(5), game.Player(0).Coins)
	assert.Equal(t, Coins(2), game.Player(1).Coins)
	assert.Equal(t, uint8(1), game.currentAge)
	assert.Equal(t, game.GetState(), StateChooseFirstPlayer)

	// === Age II ===
	err = game.ChooseFirstPlayer(1)
	require.NoError(t, err)
	assert.Equal(t, game.GetState(), StateGameTurn)

	// 1
	_, err = game.DiscardCard(cardID("Brick yard"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Horse breeders"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Library"))
	require.NoError(t, err)
	err = game.ChoosePToken(pTokenID("Law"))
	require.NoError(t, err)

	// 0
	_, err = game.DiscardCard(cardID("Tribunal"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Drying room"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Barracks"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Laboratory"))
	require.NoError(t, err)
	err = game.ChoosePToken(pTokenID("Architecture"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Forum"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("School"))
	require.NoError(t, err)
	err = game.ChoosePToken(pTokenID("Strategy"))
	require.NoError(t, err)

	// 0
	_, err = game.DiscardCard(cardID("Postrum"))
	require.NoError(t, err)
	_, err = game.DiscardCard(cardID("Aqueduct"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Brewery"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Caravansery"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Parade ground"))
	require.NoError(t, err)
	_, err = game.ConstructWonder(cardID("Shelf quarry"), wonderID("The Great Lighthouse"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Sawmill"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("Walls"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Dispensary"))
	require.NoError(t, err)
	_, err = game.DiscardCard(cardID("Archery range"))
	require.NoError(t, err)

	// 0
	_, err = game.ConstructBuilding(cardID("Customs house"))
	require.NoError(t, err)

	// === End Age II ===
	assert.Equal(t, [numPlayers]Shields{3 + 4, 0 + 3}, game.Military().ConflictPawn.Shields)
	assert.Equal(t, 2+5, len(game.discardedCards))
	assert.Equal(t, 1+0, game.WondersState.CountBuiltByPlayer(0))
	assert.Equal(t, 0+1, game.WondersState.CountBuiltByPlayer(1))
	assert.Equal(t, 10+8, countBuiltCards(game, 0))
	assert.Equal(t, 7+6, countBuiltCards(game, 1))
	assert.Equal(t, VP(3+2), countVPs(game, 0))
	assert.Equal(t, VP(2+8), countVPs(game, 1))
	assert.Equal(t, Coins(9), game.Player(0).Coins)
	assert.Equal(t, Coins(3), game.Player(1).Coins)
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
	_, err = game.ConstructBuilding(cardID("Port"))
	require.NoError(t, err)
	_, err = game.ConstructBuilding(cardID("University"))
	require.NoError(t, err)

	// === End Age III ===
	assert.Equal(t, game.GetState(), StateVictory)
	assert.Equal(t, Winner2Player, game.winner)
	assert.Equal(t, ScientificSupremacy, game.victoryType)

	assert.Equal(t, [numPlayers]Shields{3 + 4, 0 + 3}, game.Military().ConflictPawn.Shields)
	assert.Equal(t, 2+5+1, len(game.discardedCards))
	assert.Equal(t, 1+0, game.WondersState.CountBuiltByPlayer(0))
	assert.Equal(t, 0+1, game.WondersState.CountBuiltByPlayer(1))
	assert.Equal(t, 10+8+1, countBuiltCards(game, 0))
	assert.Equal(t, 7+6+1, countBuiltCards(game, 1))

	assert.Equal(t, Coins(9+11), game.Player(0).Coins)
	assert.Equal(t, Coins(3+3), game.Player(1).Coins)

	assert.Equal(t, [numVPTypes]VP{
		BlueVP:     3,
		GreenVP:    2,
		YellowVP:   3,
		PurpleVP:   0,
		WonderVP:   0,
		PTokenVP:   0,
		CoinsVP:    6,
		MilitaryVP: 5,
		SumVP:      19,
	}, game.vps[0])
	assert.Equal(t, [numVPTypes]VP{
		BlueVP:     0,
		GreenVP:    8,
		YellowVP:   0,
		PurpleVP:   0,
		WonderVP:   4,
		PTokenVP:   0,
		CoinsVP:    2,
		MilitaryVP: 0,
		SumVP:      14,
	}, game.vps[1])
}

func countBuiltCards(g *Game, i PlayerIndex) (count int) {
	for _, cc := range g.builtCards[i] {
		count += len(cc)
	}
	return
}

func countVPs(g *Game, i PlayerIndex) (count VP) {
	for vi, vps := range g.vps[i] {
		if VPType(vi) == SumVP {
			continue
		}
		count += vps
	}
	return
}
