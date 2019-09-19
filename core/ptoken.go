package core

// PToken - progress token
type PToken struct {
	Name    PTokenName
	Effects []interface{}
}

// PTokenName - name of a progress token
type PTokenName string
