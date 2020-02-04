package core

// Desk of the age's structure
type Desk interface {
	Cards() []CardName
	Check(CardName) bool
	Take(CardName) bool
}

const (
	noCardName     CardName = "-"
	hiddenCardName CardName = "???"
)

type olddeskAgeI [20]CardName

func (a olddeskAgeI) Cards() []CardName {
	a2 := a
	out := a2[:]

	// 2 3 4
	hideCardIf(out, ageIRelations, 2)
	hideCardIf(out, ageIRelations, 3)
	hideCardIf(out, ageIRelations, 4)

	// 9 10 11 12 13
	hideCardIf(out, ageIRelations, 9)
	hideCardIf(out, ageIRelations, 10)
	hideCardIf(out, ageIRelations, 11)
	hideCardIf(out, ageIRelations, 12)
	hideCardIf(out, ageIRelations, 13)

	return out
}

func (a olddeskAgeI) Check(card CardName) bool {
	if card == hiddenCardName || card == noCardName {
		return false
	}
	desk := a[:]
	index := indexOf(desk, card)
	if index < 0 {
		return false
	}
	if index >= 14 && index <= 19 {
		return true
	}
	reqs, ok := ageIRelations[index]
	if !ok {
		return false
	}
	return isAvailable(desk, reqs.left, reqs.right)
}

func (a *olddeskAgeI) Take(card CardName) bool {
	if !a.Check(card) {
		return false
	}

	index := indexOf(a[:], card)
	a[index] = noCardName
	return true
}

func hideCardIf(d []CardName, m map[cardIndex]relationNodes, i cardIndex) {
	r := m[i]
	if d[r.left] != noCardName || d[r.right] != noCardName {
		d[i] = hiddenCardName
	}
}

func isAvailable(desk []CardName, req1, req2 cardIndex) bool {
	if req1 != -1 && desk[req1] != noCardName {
		return false
	}
	if req2 != -1 && desk[req2] != noCardName {
		return false
	}
	return true
}

func indexOf(d []CardName, card CardName) cardIndex {
	for i, n := range d {
		if n == card {
			return cardIndex(i)
		}
	}
	return -1
}

type cardIndex int

type relationNodes struct {
	left, right cardIndex
}

var ageIRelations = map[cardIndex]relationNodes{
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

// Desk of I age
//
//      0  1
//     2  3  4
//    5  6  7  8
//   9 10 11 12 13
// 14 15 16 17 18 19
type deskAgeI [20]DeskCard

func newDeskAgeI(cards []CardID) (desk deskAgeI) {
	if len(cards) < 20 {
		panic("wrong amount of cards")
	}
	for i := range desk {
		desk[i].IsVisible = i <= 1 || (i >= 5 && i <= 8) || i >= 14
		desk[i].IsAvailable = i >= 14
		desk[i].ID = cards[i]
	}
	return
}
