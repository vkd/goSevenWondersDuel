package core

import "math/rand"

// PToken - progress token
type PToken struct {
	Name    PTokenName
	Effects []Effect
}

type PTokenID uint8

// PTokenName - name of a progress token
type PTokenName string

const (
	numPTokens = 10
)

var (
	listPTokens = []PToken{
		newPToken("Agriculture"),
		// newPToken("Agriculture", Coins(6), VP(4)),
		newPToken("Architecture"),
		newPToken("Economy"),
		newPToken("Law"),
		// newPToken("Law", Scales),
		newPToken("Masonry"),
		newPToken("Mathematics"),
		newPToken("Philosophy"),
		// newPToken("Philosophy", VP(7)),
		newPToken("Strategy"),
		newPToken("Theology"),
		newPToken("Urbanism"),
		// newPToken("Urbanism", Coins(6)),
	}
	_ = [1]struct{}{}[len(listPTokens)-numPTokens]

	mapPTokens = makeMapPTokensByName(listPTokens)
	_          = [1]struct{}{}[len(mapPTokens)-numPTokens]

	listPTokensIDs [numPTokens]PTokenID
)

func init() {
	for i := range listPTokensIDs {
		listPTokensIDs[i] = PTokenID(i)
	}
}

func newPToken(name PTokenName, ee ...interface{}) PToken {
	return PToken{
		Name: name,
		// Effects: ee,
	}
}

func makeMapPTokensByName(list []PToken) map[PTokenName]*PToken {
	m := map[PTokenName]*PToken{}
	for i, t := range list {
		m[t.Name] = &list[i]
	}
	return m
}

func shufflePTokens(rnd *rand.Rand) []PTokenID {
	var ptokens = listPTokensIDs
	rnd.Shuffle(len(ptokens), func(i, j int) {
		ptokens[i], ptokens[j] = ptokens[j], ptokens[i]
	})
	return ptokens[:]
}
