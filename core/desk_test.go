package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_desk_ageI(t *testing.T) {
	var d deskAgeI
	for i := range d {
		d[i] = CardName(fmt.Sprintf("card-%d", i))
	}

	cards := d.Cards()
	allNotEquals(t, cards[:2], hiddenCardName)
	allEquals(t, cards[2:5], hiddenCardName)
	allNotEquals(t, cards[5:9], hiddenCardName)
	allEquals(t, cards[9:14], hiddenCardName)
	allNotEquals(t, cards[14:], hiddenCardName)

	for i, name := range cards {
		assert.Equal(t, i >= 14, d.Check(name), "wrong value for %d index", i)
	}

	ok := d.Take(cards[16])
	assert.True(t, ok)
	cards = d.Cards()

	for i, name := range cards {
		if name == noCardName {
			continue
		}
		assert.Equal(t, i >= 14, d.Check(name), "wrong value for %d index", i)
	}

	assert.Equal(t, hiddenCardName, cards[11])

	ok = d.Take(cards[17])
	assert.True(t, ok)
	cards = d.Cards()

	assert.NotEqual(t, hiddenCardName, cards[11])

	for i, name := range cards {
		if name == noCardName {
			continue
		}
		assert.Equal(t, i == 11 || i >= 14, d.Check(name), "wrong value for %d index", i)
	}
}

func allEquals(t *testing.T, cards []CardName, expected CardName) {
	for _, name := range cards {
		assert.Equal(t, expected, name)
	}
}

func allNotEquals(t *testing.T, cards []CardName, expected CardName) {
	for _, name := range cards {
		assert.NotEqual(t, expected, name)
	}
}
