package core

// ScientificSymbol - Each time you gather a pair of identical scientific
// symbols, you may immediately choose one of the Progress tokens on the game
// board. That token will be kept in your city until the end of the game.
type ScientificSymbol uint8

var _ Effect = ScientificSymbol(0)

//ScientificSymbols
const (
	Wheel ScientificSymbol = iota
	Mortar
	Clock
	Tool
	Pen
	Astronomy
	Scales
	numOfScientificSymbols = iota
)

var (
	nameOfScientificSymbol = map[ScientificSymbol]string{
		Wheel:     "Wheel",
		Mortar:    "Mortar",
		Clock:     "Clock",
		Tool:      "Tool",
		Pen:       "Pen",
		Astronomy: "Astronomy",
		Scales:    "Scales",
	}
	_ = [1]struct{}{}[len(nameOfScientificSymbol)-numOfScientificSymbols]
)

// String - return name of the scientific symbol
func (s ScientificSymbol) String() string {
	return nameOfScientificSymbol[s]
}

func (s ScientificSymbol) applyEffect(g *Game, i PlayerIndex) {
	g.players[i].ScientificSymbols.set(s)
	if g.players[i].ScientificSymbols.uniqs() >= 6 {
		g.victory(i.winner(), WinScience)
		return
	}
	if g.players[i].ScientificSymbols[s]%2 == 0 {
		g.GettingPToken(i)
	}
}

// IsFree contruction card if this symbol is presented
// func (s ScientificSymbol) IsFree(p Player) bool {
// 	return p.ScientificSymbols[s]
// }

// func (s ScientificSymbol) Price(g *Game, i PlayerIndex) (Coins, bool) {
// 	if g.player(i).ScientificSymbols[s] {
// 		return 0, true
// 	}
// 	return 0, false
// }

// ScientificSymbols - set of scientific symbols
type ScientificSymbols [numOfScientificSymbols]uint8

func (ss *ScientificSymbols) set(s ScientificSymbol) {
	ss[s]++
}

func (ss ScientificSymbols) uniqs() (count int) {
	for _, s := range ss {
		if s > 0 {
			count++
		}
	}
	return
}
