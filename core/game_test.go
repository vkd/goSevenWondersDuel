package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZeroGame(t *testing.T) {
	game := NewGame(WithSeed(0))
	assert.Equal(t, game.GetState(), StateInit)

	wonders, ptokens, ok := game.Init()
	assert.True(t, ok)
	// [0:The Appian Way 1:The Statue of Zeus 2:The Great Library 3:Temple of Artemis
	// 4:The Hanging Gardens 5:The Great Lighthouse 6:The Mausoleum 7:The Sphinx]
	assert.Len(t, wonders, initialWonders)
	assert.Len(t, ptokens, initialPTokens)

	err := game.SelectWonders(
		// Temple of Artemis, The Great Library, The Hanging Gardens, The Sphinx
		[...]WonderName{wonders[3], wonders[2], wonders[4], wonders[7]},
		// The Appian Way, The Statue of Zeus, The Great Lighthouse, The Mausoleum
		[...]WonderName{wonders[0], wonders[1], wonders[5], wonders[6]},
	)
	assert.NoError(t, err)
	assert.Equal(t, game.GetState(), StateAgeI)
}
