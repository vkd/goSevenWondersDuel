package core

// Resource - type of goods
type Resource uint8

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

// Resources - stack of resources
type Resources [numResources]uint
