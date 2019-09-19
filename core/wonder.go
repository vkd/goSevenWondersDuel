package core

// Wonder from the Age of Antiquity
type Wonder struct {
	Name    WonderName
	Cost    Cost
	Effects []interface{}
}

// WonderName - name of a wonder
type WonderName string
