package core

// Resource - type of goods
type Resource uint8

var _ Effect = Resource(0)

// Different kind of resources
const (
	Wood Resource = iota
	Stone
	Clay
	Papyrus
	Glass
	numResources int = iota
)

// Sets of resources
var (
	rawMaterials      = []Resource{Wood, Stone, Clay}
	manufacturedGoods = []Resource{Papyrus, Glass}

	allResources = []Resource{Wood, Stone, Clay, Papyrus, Glass}
	// compile time checkers - len(allResources) == numResources
	_ = [1]struct{}{}[len(allResources)-numResources]
)

var (
	nameOfResource = map[Resource]string{
		Wood:    "Wood",
		Stone:   "Stone",
		Clay:    "Clay",
		Papyrus: "Papyrus",
		Glass:   "Glass",
	}
	_ = [1]struct{}{}[len(nameOfResource)-numResources]
)

// String representation of resource
func (r Resource) String() string {
	return nameOfResource[r]
}

func (r Resource) applyEffect(g *Game, i PlayerIndex) {
	g.player(i).Resources.addOne(r)
}

func (r Resource) applyPrice(p *Price) {
	p.Resources.addOne(r)
}

// Resources - stack of resources
type Resources [numResources]uint

// NewRes - construct a new resources stack
func NewRes(rs ...Resource) Resources {
	var out Resources
	for _, r := range rs {
		out[r]++
	}
	return out
}

func (rs Resources) ReduceOne(r Resource) Resources {
	rs.reduceOne(r)
	return rs
}

// Reduce by player's resources
func (rs Resources) Reduce(byRs Resources) Resources {
	for i := range rs {
		if rs[i] < byRs[i] {
			rs[i] = 0
		} else {
			rs[i] -= byRs[i]
		}
	}
	return rs
}

// Add resources
func (rs Resources) Add(rs2 Resources) Resources {
	for i := range rs {
		rs[i] += rs2[i]
	}
	return rs
}

func (rs *Resources) add(rs2 Resources) {
	for i := range rs {
		rs[i] += rs2[i]
	}
}

func (rs *Resources) addOne(r Resource) {
	rs[r]++
}

func (rs *Resources) reduceOne(r Resource) {
	if rs[r] > 0 {
		rs[r]--
	}
}

func (rs Resources) applyPrice(p *Price) {
	p.Resources.add(rs)
}
