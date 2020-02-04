package game

func (r Resource) effect(g *Game, i PlayerIndex) {
	g.player(i).Resources.Change(r, 1)
}

// Presets of Resources
var (
	EmptyResources = Resources{}
	Empty          = Resources{}
)

// Resources - stack of resources, total count of every kind
type Resources [numResources]int

// Sum of two stack resources
func Sum(r1, r2 Resources) Resources {
	for i := range r1 {
		r1[i] += r2[i]
	}
	return r1
}

// Add another stack of resources
func (r Resources) Add(rs Resources) Resources {
	return Sum(r, rs)
}

// Change ...
func (r Resources) Change(rt Resource, delta int) Resources {
	if delta < 0 && -1*delta > r[rt] {
		r[rt] = 0
	} else {
		r[rt] += delta
	}
	return r
}

// Sub ...
func (r Resources) Sub(rs Resources) Resources {
	for i := range r {
		if rs[i] > r[i] {
			r[i] = 0
		} else {
			r[i] -= rs[i]
		}
	}
	return r
}

// TakeOne ...
func (r Resources) TakeOne(res Resource) Resources {
	return r.Change(res, -1)
}

// IsPositive ...
func (r Resources) IsPositive(res Resource) bool {
	return r[res] > 0
}

// IsZero ...
func (r Resources) IsZero(res Resource) bool {
	return r[res] == 0
}

// // Inc ...
// func (r *Resources) Inc(res Resource) {
// 	r[res]++
// }

func (r Resources) effect(g *Game, i PlayerIndex) {
	g.log.Infof("Add resources: %v", r)
	g.player(i).Resources.Add(r)
}

// Res ...
func Res(rt Resource, val int) Resources {
	return EmptyResources.Change(rt, val)
}

// MaybeRes ...
type MaybeRes struct {
	R      Resource
	Exists bool
}

// Set ...
func (m *MaybeRes) Set(r Resource) {
	m.R = r
	m.Exists = true
}
