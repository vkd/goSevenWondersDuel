package core

import (
	"errors"
	"fmt"
)

func GetUncovered(a AgeStructure) []CardID {
	out := make([]CardID, 0, 6)
	for i, s := range a.Cards {
		if !s.Exists {
			continue
		}
		var covered bool
		for _, byID := range CoverageByForAge(a.age)[StructureID(i)] {
			if a.Cards[byID].Exists {
				covered = true
				break
			}
		}
		if !covered {
			out = append(out, s.CardID)
		}
	}
	return out
}

type AgeStructure struct {
	Cards [SizeAge]StructureCard

	age      Age
	shuffler CardShuffler
}

type CardShuffler interface {
	Next() (CardID, error)
}

func NewAgeStructure(age Age, shuffler CardShuffler) (AgeStructure, error) {
	a := AgeStructure{
		age:      age,
		shuffler: shuffler,
	}

	var isHidden [len(a.Cards)]bool
	for _, hid := range HiddenCardsForAge(age) {
		isHidden[hid] = true
	}

	for i := range a.Cards {
		a.Cards[i].placeHidden()
	}
	for i := range a.Cards {
		if isHidden[i] {
			continue
		}
		err := a.Cards[i].Open(a.shuffler)
		if err != nil {
			return a, fmt.Errorf("cannot open card: %w", err)
		}
	}

	for i := range a.Cards {
		if !isHidden[i] {
			continue
		}

		id := StructureID(i)
		err := a.IsCovered(id)
		if err == nil {
			return a, fmt.Errorf("uncorrect hiddens: %d card is not covered", id)
		}
	}
	return a, nil
}

func (a AgeStructure) Take(id StructureID) (AgeStructure, error) {
	if id < 0 || int(id) >= len(a.Cards) {
		return a, fmt.Errorf("id is out of range [0:%d)", len(a.Cards))
	}

	err := a.IsCovered(id)
	if err != nil {
		return a, fmt.Errorf("card is covered: %w", err)
	}

	err = a.Cards[id].Take()
	if err != nil {
		return a, fmt.Errorf("card is not buildable: %w", err)
	}

	for _, coveredID := range CoveredsForAge(a.age)[id] {
		if !a.Cards[coveredID].Exists {
			continue
		}

		err = a.IsCovered(coveredID)
		if err != nil {
			continue
		}
		err = a.Cards[coveredID].Open(a.shuffler)
		if err != nil {
			return a, fmt.Errorf("cannot open card ID %d: %w", coveredID, err)
		}
	}

	return a, nil
}

func (a AgeStructure) IsCovered(id StructureID) error {
	for _, byID := range CoverageByForAge(a.age)[id] {
		if a.Cards[byID].Exists {
			return fmt.Errorf("by card id = %d", byID)
		}
	}
	return nil
}

type StructureCard struct {
	Exists bool
	FaceUp bool
	CardID CardID
}

func (s *StructureCard) Take() error {
	if !s.Exists {
		return ErrStructureCardNotExists
	}

	s.Exists = false
	return nil
}

func (s *StructureCard) Open(shuffler CardShuffler) error {
	if !s.Exists {
		return ErrStructureCardNotExists
	}

	if s.FaceUp {
		return nil
	}

	newID, err := shuffler.Next()
	if err != nil {
		return fmt.Errorf("shuffler: %w", err)
	}

	s.FaceUp = true
	s.CardID = newID
	return nil
}

func (s *StructureCard) placeHidden() {
	s.Exists = true
}

var ErrStructureCardNotExists = errors.New("not exists")

type StructureID = cardIndex

type CoverageByAgeStructure map[StructureID][]StructureID

func CoverageByForAge(a Age) CoverageByAgeStructure {
	switch a {
	case AgeI:
		return ageICoveredBy
	case AgeII:
		return ageIICoveredBy
	case AgeIII:
		return ageIIICoveredBy
	}
	return nil
}

type CoveredsAgeStructure map[StructureID][]StructureID

func CoveredsForAge(a Age) CoveredsAgeStructure {
	switch a {
	case AgeI:
		return ageICovers
	case AgeII:
		return ageIICovers
	case AgeIII:
		return ageIIICovers
	}
	return nil
}

func makeRevertCovers(coveredBy CoverageByAgeStructure) CoveredsAgeStructure {
	out := make(CoveredsAgeStructure)
	for k, v := range coveredBy {
		for _, idx := range v {
			out[idx] = append(out[idx], k)
		}
	}
	return out
}

type HiddenCardsAgeStructure []StructureID

func HiddenCardsForAge(a Age) HiddenCardsAgeStructure {
	switch a {
	case AgeI:
		return ageIHiddenCards
	case AgeII:
		return ageIIHiddenCards
	case AgeIII:
		return ageIIIHiddenCards
	}
	return nil
}

var (
	//      0  1
	//     2  3  4
	//    5  6  7  8
	//   9 10 11 12 13
	// 14 15 16 17 18 19
	ageICoveredBy = CoverageByAgeStructure{
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
	ageICovers      = makeRevertCovers(ageICoveredBy)
	ageIHiddenCards = HiddenCardsAgeStructure{2, 3, 4, 9, 10, 11, 12, 13}

	// 0  1  2  3  4  5
	//  6  7  8  9  10
	//   11 12 13 14
	//    15 16 17
	//     18 19
	ageIICoveredBy = CoverageByAgeStructure{
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
	ageIICovers      = makeRevertCovers(ageIICoveredBy)
	ageIIHiddenCards = HiddenCardsAgeStructure{6, 7, 8, 9, 10, 15, 16, 17}

	//     0   1
	//   2   3   4
	// 5   6   7   8
	//   9       10
	// 11  12  13  14
	//   15  16  17
	//     18  19
	ageIIICoveredBy = CoverageByAgeStructure{
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
	ageIIICovers      = makeRevertCovers(ageIIICoveredBy)
	ageIIIHiddenCards = HiddenCardsAgeStructure{2, 3, 4, 9, 10, 15, 16, 17}
)
