package core

// Chain ...
type Chain uint8

var _ Effect = Chain(0)

// Chains
const (
	Horseshoe Chain = iota
	Sword
	Wall
	Target
	Helm
	Book
	Gear
	Harp
	Lamp
	Mask
	Column
	Moon
	Sun
	Water
	Pantheon
	Bottle
	Barrel
	numOfChains = iota
)

// String representation of chain symbol
func (c Chain) String() string {
	return nameOfChain[c]
}

// Apply effect of chain symbol
func (c Chain) Apply(g *Game, i PlayerIndex) {
	g.player(i).Chains[c] = true
}

// Chains - set of chains
//
// TODO - use bits
type Chains [numOfChains]bool

// NewChains with preinstalled values
func NewChains(cc ...Chain) Chains {
	var out Chains
	for _, c := range cc {
		out[c] = true
	}
	return out
}

// Strings representation of chains
func (cs Chains) Strings() []string {
	var out []string
	for i, ok := range cs {
		if ok {
			out = append(out, Chain(i).String())
		}
	}

	return out
}

// Contain a chain
func (cs Chains) Contain(c Chain) bool {
	return cs[c]
}

var (
	nameOfChain = map[Chain]string{
		Horseshoe: "Horseshoe",
		Sword:     "Sword",
		Wall:      "Wall",
		Target:    "Target",
		Helm:      "Helm",
		Book:      "Book",
		Gear:      "Gear",
		Harp:      "Harp",
		Lamp:      "Lamp",
		Mask:      "Mask",
		Column:    "Column",
		Moon:      "Moon",
		Sun:       "Sun",
		Water:     "Water",
		Pantheon:  "Pantheon",
		Bottle:    "Bottle",
		Barrel:    "Barrel",
	}
	_ = [1]struct{}{}[numOfChains-len(nameOfChain)]
)

// MaybeChain - optional field of a chain
type MaybeChain struct {
	OK bool
	Chain
}

// Set a value
func (m *MaybeChain) Set(c Chain) {
	m.OK = true
	m.Chain = c
}
