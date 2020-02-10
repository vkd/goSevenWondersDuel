package core

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

func (c Chain) applyEffect(g *Game, i PlayerIndex) {
	g.player(i).Chains.Set(c)
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

func (cs *Chains) Set(c Chain) {
	cs[c] = true
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
