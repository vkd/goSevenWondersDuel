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

// Resources - stack of resources
type Resources [numResources]uint
