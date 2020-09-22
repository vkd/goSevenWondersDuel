package core

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAgeStructure(t *testing.T) {
	cardSource := &singleCardSource{
		cards: []CardID{0, 1, 5, 6, 7, 8, 14, 15, 16, 17, 18, 19},
	}

	//      0  1
	//     2  3  4
	//    5  6  7  8
	//   9 10 11 12 13
	// 14 15 16 17 18 19
	as, err := NewAgeStructure(AgeI, cardSource)
	require.NoError(t, err)

	//      0  1
	//     *  *  *
	//    5  6  7  8
	//   * ** ** ** **
	// 14 15 16 17 __ 19
	as, err = as.Take(18)
	require.NoError(t, err)

	//      0  1
	//     *  *  *
	//    5  6  7  8
	//   * ** ** ** 13
	// 14 15 16 17 __ __
	cardSource.cards = append(cardSource.cards, 13)
	as, err = as.Take(19)
	require.NoError(t, err)

	//      0  1
	//     *  *  *
	//    5  6  7  8
	//   * ** ** ** __
	// 14 15 16 17 __ __
	as, err = as.Take(13)
	require.NoError(t, err)

	//      0  1
	//     *  *  *
	//    5  6  7  8
	//   * ** ** 12 __
	// 14 15 16 __ __ __
	cardSource.cards = append(cardSource.cards, 12)
	as, err = as.Take(17)
	require.NoError(t, err)

	//      0  1
	//     *  *  *
	//    5  6  7  8
	//   * ** ** __ __
	// 14 15 16 __ __ __
	as, err = as.Take(12)
	require.NoError(t, err)

	//      0  1
	//     *  *  *
	//    5  6  7  _
	//   * ** ** __ __
	// 14 15 16 __ __ __
	as, err = as.Take(8)
	require.NoError(t, err)
}

var errTest = errors.New("test error")

func TestAgeStructure_errors(t *testing.T) {
	cardSource := &singleCardSource{err: errTest}

	//      0  1
	//     2  3  4
	//    5  6  7  8
	//   9 10 11 12 13
	// 14 15 16 17 18 19
	_, err := NewAgeStructure(AgeI, cardSource)
	require.Error(t, err)

	cardSource = &singleCardSource{
		cards: []CardID{0, 1, 5, 6, 7, 8, 14, 15, 16, 17, 18, 19},
	}
	//      0  1
	//     2  3  4
	//    5  6  7  8
	//   9 10 11 12 13
	// 14 15 16 17 18 19
	as, err := NewAgeStructure(AgeI, cardSource)
	require.NoError(t, err)

	// id is out of range
	_, err = as.Take(20)
	assert.Error(t, err)

	// covered card
	_, err = as.Take(13)
	assert.Error(t, err)

	//      0  1
	//     *  *  *
	//    5  6  7  8
	//   * ** ** ** **
	// 14 15 16 17 __ 19
	as, err = as.Take(18)
	require.NoError(t, err)

	// card is not exist
	_, err = as.Take(18)
	assert.Error(t, err)

	// cannot get next card
	cardSource.err = errTest
	_, err = as.Take(19)
	assert.Error(t, err)
}

type singleCardSource struct {
	cards []CardID
	next  int
	err   error
}

func (c *singleCardSource) Next() (CardID, error) {
	if c.err != nil {
		return 0, c.err
	}

	cid := c.cards[c.next]
	c.next++
	return cid, nil
}
