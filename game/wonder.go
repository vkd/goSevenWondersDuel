package game

import "fmt"

func (n WonderName) find() *Wonder {
	w, ok := mapWonders[n]
	if !ok {
		panic(fmt.Sprintf("cannot find %q wonder", n))
	}
	return w
}

// WonderNames - list of name of wonders
type WonderNames []WonderName

// NewAllWonderNames - new slice of wonder name from all wonders
func NewAllWonderNames() WonderNames {
	var out = make(WonderNames, len(listWonders))
	for i := range listWonders {
		out[i] = listWonders[i].Name
	}
	return out
}

// Exists name in slice of wonder names
func (w WonderNames) Exists(n WonderName) bool {
	for _, wn := range w {
		if wn == n {
			return true
		}
	}
	return false
}

// Shuffle slice of wonder names
func (w WonderNames) Shuffle() WonderNames {
	rnd.Shuffle(len(w), func(i, j int) { w[i], w[j] = w[j], w[i] })
	return w
}

// Append wonder name
func (w *WonderNames) Append(nn ...WonderName) {
	*w = append(*w, nn...)
}

// Index of searched name
func (w WonderNames) Index(n WonderName) int {
	for i := range w {
		if w[i] == n {
			return i
		}
	}
	return -1
}

// Count - return count of this name
func (w WonderNames) Count(n WonderName) int {
	var count int
	for _, name := range w {
		if name == n {
			count++
		}
	}
	return count
}
