package core

// ScientificSymbol - Each time you gather a pair of identical scientific
// symbols, you may immediately choose one of the Progress tokens on the game
// board. That token will be kept in your city until the end of the game.
type ScientificSymbol uint8

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
