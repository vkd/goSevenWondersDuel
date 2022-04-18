package core

import (
	"errors"
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

func (c CardState) String() string {
	return fmt.Sprintf("%v: faced=%t, exists=%t, covered=%t", c.ID, c.FaceUp, c.Exists, c.Covered)
}

// TestBuilt - card could be built
func (c CardState) TestBuilt() error {
	if !c.Exists {
		return fmt.Errorf("does not exist")
	}
	if c.Covered {
		return fmt.Errorf("is covered")
	}
	return nil
}

// CardsState on the table
type CardsState = AgeCards // [SizeAge]CardState

// func (s CardsState) testBuilt(i cardIndex) error {
// 	return s[i].TestBuilt()
// }

// func (s *CardsState) set(i cardIndex, id CardID, covered bool) {
// 	s[i].FaceUp = true
// 	s[i].Exists = true
// 	s[i].Covered = covered
// 	s[i].ID = id
// }

// func (s *CardsState) take(i cardIndex) {
// 	s[i].Exists = false
// }

// func (s *CardsState) open(i cardIndex, cards []CardID) {
// 	s[i].FaceUp = true
// 	s[i].ID = cards[i]
// }

// func (s *CardsState) free(i cardIndex) {
// 	s[i].Covered = false
// }

// func (s *CardsState) hide(i cardIndex) {
// 	if s[i].Covered {
// 		s[i].ID = 0
// 		s[i].FaceUp = false
// 	}
// }

// func (s CardsState) anyExists() bool {
// 	for _, cs := range s {
// 		if cs.Exists {
// 			return true
// 		}
// 	}
// 	return false
// }

// func (s CardsState) isAnyExists(idxs []cardIndex) bool {
// 	for _, idx := range idxs {
// 		if s[idx].Exists {
// 			return true
// 		}
// 	}
// 	return false
// }

var (
	structureAgeI   = newAgeStructure(ageICoveredBy, ageIHiddenCards)
	structureAgeII  = newAgeStructure(ageIICoveredBy, ageIIHiddenCards)
	structureAgeIII = newAgeStructure(ageIIICoveredBy, ageIIIHiddenCards)
)

type cardIndex uint8

type cardRelations = CoverageByAgeStructure

type ageStructure struct {
	coveredBy   cardRelations
	covers      CoveredsAgeStructure
	hiddenCards HiddenCardsAgeStructure
}

func newAgeStructure(coveredBy cardRelations, hiddenCards HiddenCardsAgeStructure) *ageStructure {
	age := ageStructure{
		coveredBy:   coveredBy,
		covers:      makeRevertCovers(coveredBy),
		hiddenCards: hiddenCards,
	}
	return &age
}

type shuffledCards struct {
	cards []CardID
	index int
}

func (s *shuffledCards) Next() (CardID, error) {
	if s.index < 0 || s.index >= len(s.cards) {
		return 0, errors.New("wrong index")
	}
	cid := s.cards[s.index]
	s.index++
	return cid, nil
}

type ageDesk struct {
	cards    []CardID
	age      Age
	shuffler CardShuffler

	ageStructure AgeStructure
}

func newAgeDesk(cards []CardID, age Age) (ageDesk, error) {
	var desk ageDesk
	desk.cards = cards
	desk.age = age
	desk.shuffler = &shuffledCards{cards: cards}

	var err error
	desk.ageStructure, err = NewAgeStructure(desk.shuffler, CoverageByForAge(age), HiddenCardsForAge(age))
	if err != nil {
		return desk, err
	}

	return desk, nil
}

func (d *ageDesk) Build(id CardID) error {
	idx, ok := indexOfCards(id, d.cards)
	if !ok {
		return fmt.Errorf("card not found (id = %d)", id)
	}

	nextAs, err := d.ageStructure.Take(idx, d.shuffler, CoverageByForAge(d.age), CoveredsForAge(d.age))
	if err != nil {
		return fmt.Errorf("cannot take card on structure: %w", err)
	}

	d.ageStructure = nextAs
	return nil
}

func (d *ageDesk) testBuild(id CardID) error {
	idx, ok := indexOfCards(id, d.cards)
	if !ok {
		return fmt.Errorf("card does not exist (id = %d)", id)
	}
	as := d.ageStructure
	_, err := as.Take(idx, d.shuffler, CoverageByForAge(d.age), CoveredsForAge(d.age))
	if err != nil {
		return fmt.Errorf("cannot take card [%d] on structure: %w", idx, err)
	}
	return nil
}

func indexOfCards(id CardID, cards []CardID) (cardIndex, bool) {
	for i, c := range cards {
		if c == id {
			return cardIndex(i), true
		}
	}
	return 0, false
}
