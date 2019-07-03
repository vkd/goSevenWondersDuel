package game

// Wonder from the Age of Antiquity
type Wonder struct {
	Name    WonderName
	Cost    CostOfCard
	Effects []Effect
}

// WonderName - name of a wonder
type WonderName string

func (n WonderName) find() *Wonder {
	return mapWonders[n]
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

const (
	numWonders = 12
)

var (
	listWonders = []Wonder{
		newWonder("Temple of Artemis", Money(12), RepeatTurn(), Cost(Wood, Stone, Glass, Papyrus)),
		newWonder("The Great Lighthouse", OneOfAnyMarket(rawMaterials), VP(4), Cost(Wood, Stone, Papyrus, Papyrus)),
		newWonder("The Colossus", Shields(2), VP(3), Cost(Clay, Clay, Clay, Glass)),
		newWonder("The Pyramids", VP(9), Cost(Stone, Stone, Stone, Papyrus)),
		newWonder("The Mausoleum", VP(2), Cost(Clay, Clay, Glass, Glass, Papyrus)),
		newWonder("The Statue of Zeus", Shields(1), VP(3)),
		newWonder("The Appian Way", Money(3), Opponent(DiscardMoney(3)), RepeatTurn(), VP(3), Cost(Stone, Stone, Clay, Clay, Papyrus)),
		newWonder("Circus Maximus", Shields(1), VP(3), Cost(Stone, Stone, Wood, Glass)),
		newWonder("The Great Library", VP(4), Cost(Wood, Wood, Wood, Glass, Papyrus)),
		newWonder("Piraeus", OneOfAnyMarket(manufacturedGoods), RepeatTurn(), VP(2), Cost(Wood, Wood, Stone, Clay)),
		newWonder("The Hanging Gardens", Money(6), RepeatTurn(), VP(3), Cost(Wood, Wood, Glass, Papyrus)),
		newWonder("The Sphinx", RepeatTurn(), VP(6), Cost(Stone, Clay, Glass, Glass)),
	}
	_ = [1]struct{}{}[numWonders-len(listWonders)]

	mapWonders = makeMapWondersByName()
	_          = [1]struct{}{}[numWonders-len(mapWonders)]
)

func newWonder(name WonderName, args ...interface{}) (w Wonder) {
	w.Name = name
	for i := range args {
		switch a := args[i].(type) {
		case CostOfCard:
			w.Cost = a
		case Effect:
			w.Effects = append(w.Effects, a)
		default:
			panic("Not implemented")
		}
	}
	return w
}

func makeMapWondersByName() map[WonderName]*Wonder {
	m := map[WonderName]*Wonder{}
	for i, w := range listWonders {
		m[w.Name] = &listWonders[i]
	}
	return m
}
