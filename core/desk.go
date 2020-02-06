package core

// Desk of the age's structure
type Desk interface {
	Cards() []CardName
	Check(CardName) bool
	Take(CardName) bool
}

type cardIndex int

var ageICoveredBy = map[cardIndex][]cardIndex{
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

// 18: {12, 13},
// 19: {13},
var ageICovers = makeRevertCovers(ageICoveredBy)

func makeRevertCovers(coveredBy map[cardIndex][]cardIndex) map[cardIndex][]cardIndex {
	out := make(map[cardIndex][]cardIndex)
	for k, v := range coveredBy {
		for _, idx := range v {
			out[idx] = append(out[idx], k)
		}
	}
	return out
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

func (d *deskAgeI) Build(id CardID) bool {
	idx, ok := d.getIndex(id)
	if !ok {
		return false
	}

	d[idx].IsSkipped = true

	for _, cover := range ageICovers[idx] {
		if d.isFree(cover) {
			d.open(cover)
		}
	}

	return true
}

func (d *deskAgeI) getIndex(id CardID) (cardIndex, bool) {
	for i, c := range d {
		if c.ID == id {
			return cardIndex(i), true
		}
	}
	return 0, false
}

func (d *deskAgeI) open(idx cardIndex) {
	d[idx].IsAvailable = true
	d[idx].IsVisible = true
}

func (d *deskAgeI) take(idx cardIndex) {
	d[idx].IsSkipped = true
}

func (d *deskAgeI) isFree(idx cardIndex) bool {
	for _, cover := range ageICoveredBy[idx] {
		if !d[cover].IsSkipped {
			return false
		}
	}
	return true
}
