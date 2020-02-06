package core

import (
	"fmt"
)

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

type cardIndex int

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

type CardState struct {
	Hidden     bool
	Built      bool
	Accessible bool
	ID         CardID
}

type CardsState [SizeAge]CardState

func (s *CardsState) build(idx cardIndex, str *ageStructure, cards []CardID) error {
	if s[idx].Built {
		return fmt.Errorf("card (index=%d) is already built", idx)
	}
	if !s[idx].Accessible {
		return fmt.Errorf("card (index=%d) is not accessible", idx)
	}

	s[idx].Built = true
	s[idx].Accessible = false

	for _, cover := range str.covers[idx] {
		if s.isAllBuilt(str.coveredBy[cover]) {
			s.open(cover, cards)
		}
	}
	return nil
}

func (s *CardsState) open(idx cardIndex, cards []CardID) {
	s[idx].Accessible = true
	s[idx].Hidden = false
	s[idx].ID = cards[idx]
}

func (s CardsState) isAllBuilt(coveredBy []cardIndex) bool {
	for _, idx := range coveredBy {
		if !s[idx].Built {
			return false
		}
	}
	return true
}

type ageDesk struct {
	structure *ageStructure
	cards     []CardID
	state     CardsState
}

func newAgeDesk(structure *ageStructure, cards []CardID) (desk ageDesk, _ error) {
	if len(cards) < 20 {
		return desk, fmt.Errorf("wrong amount of cards (less than 20): %d", len(cards))
	}

	desk.structure = structure
	desk.cards = cards

	for i := range desk.state {
		_, covered := structure.coveredBy[cardIndex(i)]
		desk.state[i].Accessible = !covered
		desk.state[i].ID = cards[i]
	}

	// hide cards
	for _, i := range structure.hiddenCards {
		if !desk.state[i].Accessible {
			desk.state[i].ID = 0
			desk.state[i].Hidden = true
		}
	}
	return desk, nil
}

func (d *ageDesk) Build(id CardID) error {
	idx, ok := d.getIndex(id)
	if !ok || !d.state[idx].Accessible || d.state[idx].Built || d.state[idx].Hidden {
		return fmt.Errorf("wrong card id")
	}

	return d.state.build(idx, d.structure, d.cards)
}

func (d *ageDesk) getIndex(id CardID) (cardIndex, bool) {
	for i, c := range d.cards {
		if c == id {
			return cardIndex(i), true
		}
	}
	return 0, false
}
