package game

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
