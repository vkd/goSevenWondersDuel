package core

import "math/rand"

// TakeNfromM takes n numbers from shuffled sequence [0, 1, ..., m].
func TakeNfromM(n, m int, rnd *rand.Rand) []int {
	var out = make([]int, m)
	for i := range out {
		out[i] = i
	}
	rnd.Shuffle(m, func(i, j int) { out[i], out[j] = out[j], out[i] })
	return out[:n]
}
