package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//   0
	// 1   2
	//   3
	testCoveredBy = cardRelations{
		0: {1, 2},
		1: {3},
		2: {3},
	}
	testHiddenCards = []cardIndex{1, 2}
	testStructure   = newAgeStructure(testCoveredBy, testHiddenCards)

	testCards = []CardID{10, 11, 12, 13}
)

func Test_ageDesk_Build(t *testing.T) {
	var desk, err = newAgeDesk(testStructure, testCards)
	require.NoError(t, err)

	// Precheck
	assert.False(t, desk.state[1].FaceUp)
	assert.Zero(t, desk.state[1].ID)
	assert.Error(t, desk.Build(CardID(12)))

	// Action
	assert.NoError(t, desk.Build(CardID(13)))

	// Postcheck
	assert.True(t, desk.state[2].FaceUp)
	assert.Equal(t, CardID(12), desk.state[2].ID)
	assert.NoError(t, desk.Build(CardID(12)))
}
