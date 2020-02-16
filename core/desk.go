package core

import (
	"fmt"
)

// CardState on the table
type CardState struct {
	// FaceUp - card is displayed face up
	FaceUp bool
	// Exists - on this position card is exists
	Exists bool
	// Covered - this card could not be built
	Covered bool
	// ID - actual card ID
	ID CardID
}

// TestBuilt - card could be built
func (c CardState) TestBuilt() bool {
	return c.Exists && !c.Covered
}

// CardsState on the table
type CardsState [SizeAge]CardState

func (s CardsState) testBuilt(i cardIndex) bool {
	return s[i].TestBuilt()
}

func (s *CardsState) set(i cardIndex, id CardID, covered bool) {
	s[i].FaceUp = true
	s[i].Exists = true
	s[i].Covered = covered
	s[i].ID = id
}

func (s *CardsState) take(i cardIndex) {
	s[i].Exists = false
}

func (s *CardsState) open(i cardIndex, cards []CardID) {
	s[i].FaceUp = true
	s[i].ID = cards[i]
}

func (s *CardsState) free(i cardIndex) {
	s[i].Covered = false
}

func (s *CardsState) hide(i cardIndex) {
	if s[i].Covered {
		s[i].ID = 0
		s[i].FaceUp = false
	}
}

func (ss CardsState) anyExists() bool {
	for _, s := range ss {
		if s.Exists {
			return true
		}
	}
	return false
}

func (s CardsState) isAnyExists(idxs []cardIndex) bool {
	for _, idx := range idxs {
		if s[idx].Exists {
			return true
		}
	}
	return false
}

var (
	//      0  1
	//     2  3  4
	//    5  6  7  8
	//   9 10 11 12 13
	// 14 15 16 17 18 19
	ageICoveredBy = map[cardIndex][]cardIndex{
		0: {2, 3},
		1: {3, 4},

		2: {5, 6},
		3: {6, 7},
		4: {7, 8},

		5: {9, 10},
		6: {10, 11},
		7: {11, 12},
		8: {12, 13},

		9:  {14, 15},
		10: {15, 16},
		11: {16, 17},
		12: {17, 18},
		13: {18, 19},
	}
	ageIHiddenCards = []cardIndex{2, 3, 4, 9, 10, 11, 12, 13}

	// 0  1  2  3  4  5
	//  6  7  8  9  10
	//   11 12 13 14
	//    15 16 17
	//     18 19
	ageIICoveredBy = map[cardIndex][]cardIndex{
		0: {6},
		1: {6, 7},
		2: {7, 8},
		3: {8, 9},
		4: {9, 10},
		5: {10},

		6:  {11},
		7:  {11, 12},
		8:  {12, 13},
		9:  {13, 14},
		10: {14},

		11: {15},
		12: {15, 16},
		13: {16, 17},
		14: {17},

		15: {18},
		16: {18, 19},
		17: {19},
	}
	ageIIHiddenCards = []cardIndex{6, 7, 8, 9, 10, 15, 16, 17}

	//     0   1
	//   2   3   4
	// 5   6   7   8
	//   9       10
	// 11  12  13  14
	//   15  16  17
	//     18  19
	ageIIICoveredBy = map[cardIndex][]cardIndex{
		0: {2, 3},
		1: {3, 4},

		2: {5, 6},
		3: {6, 7},
		4: {7, 8},

		5: {9},
		6: {9},
		7: {10},
		8: {10},

		9:  {11, 12},
		10: {13, 14},

		11: {15},
		12: {15, 16},
		13: {16, 17},
		14: {17},

		15: {18},
		16: {18, 19},
		17: {19},
	}
	ageIIIHiddenCards = []cardIndex{2, 3, 4, 9, 10, 15, 16, 17}
)

var (
	structureAgeI   = newAgeStructure(ageICoveredBy, ageIHiddenCards)
	structureAgeII  = newAgeStructure(ageIICoveredBy, ageIIHiddenCards)
	structureAgeIII = newAgeStructure(ageIIICoveredBy, ageIIIHiddenCards)
)

type cardIndex uint8

type cardRelations map[cardIndex][]cardIndex

type ageStructure struct {
	coveredBy   cardRelations
	covers      cardRelations
	hiddenCards []cardIndex
}

func newAgeStructure(coveredBy cardRelations, hiddenCards []cardIndex) *ageStructure {
	age := ageStructure{
		coveredBy:   coveredBy,
		covers:      makeRevertCovers(coveredBy),
		hiddenCards: hiddenCards,
	}
	return &age
}

func makeRevertCovers(coveredBy map[cardIndex][]cardIndex) map[cardIndex][]cardIndex {
	out := make(map[cardIndex][]cardIndex)
	for k, v := range coveredBy {
		for _, idx := range v {
			out[idx] = append(out[idx], k)
		}
	}
	return out
}

type ageDesk struct {
	structure *ageStructure
	cards     []CardID
	state     CardsState
}

func newAgeDesk(structure *ageStructure, cards []CardID) (desk ageDesk, _ error) {
	desk.structure = structure
	desk.cards = cards

	for i, id := range cards {
		coveredBy, ok := structure.coveredBy[cardIndex(i)]
		covered := ok && len(coveredBy) > 0
		desk.state.set(cardIndex(i), id, covered)
	}

	// hide cards
	for _, i := range structure.hiddenCards {
		desk.state.hide(i)
	}
	return desk, nil
}

func (d *ageDesk) Build(id CardID) error {
	idx, ok := indexOfCards(id, d.cards)
	if !ok || !d.state.testBuilt(idx) {
		return fmt.Errorf("wrong card id = %d", id)
	}

	d.state.take(idx)

	for _, cover := range d.structure.covers[idx] {
		isCovered := d.state.isAnyExists(d.structure.coveredBy[cover])
		if !isCovered {
			d.state.free(cover)
			d.state.open(cover, d.cards)
		}
	}
	return nil
}

func (d *ageDesk) testBuild(id CardID) bool {
	idx, ok := indexOfCards(id, d.cards)
	return ok && d.state.testBuilt(idx)
}

func indexOfCards(id CardID, cards []CardID) (cardIndex, bool) {
	for i, c := range cards {
		if c == id {
			return cardIndex(i), true
		}
	}
	return 0, false
}
