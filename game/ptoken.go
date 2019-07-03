package game

const (
	numPTokens = 10
)

// PToken - progress token
type PToken struct {
	Name    PTokenName
	Effects []Effect
}

// PTokenName - name of a progress token
type PTokenName string

func (n PTokenName) find() *PToken {
	return mapPTokens[n]
}

// func shuffledPTokenNames() [numPTokens]PTokenName {
// 	var res [numPTokens]PTokenName
// 	for i := range res {
// 		res[i] = listPTokens[i].Name
// 	}
// 	rnd.Shuffle(numPTokens, func(i, j int) { res[i], res[j] = res[j], res[i] })
// 	return res
// }

// PTokenNames - list of progress tokens
type PTokenNames []PTokenName

// NewAllPTokenNames - new slice of progress tokens
func NewAllPTokenNames() PTokenNames {
	var out = make(PTokenNames, len(listPTokens))
	for i := range listPTokens {
		out[i] = listPTokens[i].Name
	}
	return out
}

// Shuffle list of progress tokens
func (ns PTokenNames) Shuffle() PTokenNames {
	rnd.Shuffle(len(ns), func(i, j int) { ns[i], ns[j] = ns[j], ns[i] })
	return ns
}

var (
	listPTokens = []PToken{
		newPToken("Agriculture", Money(6), VP(4)),
		newPToken("Architecture"),
		newPToken("Economy"),
		newPToken("Law", Scales),
		newPToken("Masonry"),
		newPToken("Mathematics"),
		newPToken("Philosophy", VP(7)),
		newPToken("Strategy"),
		newPToken("Theology"),
		newPToken("Urbanism", Money(6)),
	}
	_ = [1]struct{}{}[len(listPTokens)-numPTokens]

	mapPTokens = makeMapPTokensByName()
	_          = [1]struct{}{}[len(mapPTokens)-numPTokens]
)

func newPToken(name PTokenName, ee ...Effect) PToken {
	return PToken{
		Name:    name,
		Effects: ee,
	}
}

func makeMapPTokensByName() map[PTokenName]*PToken {
	m := map[PTokenName]*PToken{}
	for i, t := range listPTokens {
		m[t.Name] = &listPTokens[i]
	}
	return m
}
