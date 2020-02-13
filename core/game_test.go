package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestZeroGame(t *testing.T) {
	game, err := NewGame(WithSeed(0))
	require.NoError(t, err)
	assert.Equal(t, game.GetState(), StateNone)

	wonders, ptokens, ok := game.Init()
	assert.True(t, ok)
	// [0:The Appian Way 1:The Statue of Zeus 2:The Great Library 3:Temple of Artemis
	// 4:The Hanging Gardens 5:The Great Lighthouse 6:The Mausoleum 7:The Sphinx]
	assert.Len(t, wonders, initialWonders)
	assert.Len(t, ptokens, initialPTokens)

	err = game.SelectWonders(
		// Temple of Artemis, The Great Library, The Hanging Gardens, The Sphinx
		[...]WonderID{wonders[3], wonders[2], wonders[4], wonders[7]},
		// The Appian Way, The Statue of Zeus, The Great Lighthouse, The Mausoleum
		[...]WonderID{wonders[0], wonders[1], wonders[5], wonders[6]},
	)
	assert.NoError(t, err)
	assert.Equal(t, game.GetState(), StateGameTurn)
}
