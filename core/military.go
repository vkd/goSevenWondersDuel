package core

// ConflictPawn - Conflict Pawn.
// The Conflict pawn indicates on the board one city’s military advantage over the other.
type ConflictPawn struct {
	Shields [numPlayers]Shields
}

// Position of the pawn on the board.
//
// [-9 . . -6 . . -3 . . 0 . . 3 . . 6 . . 9]
func (c ConflictPawn) Position() int {
	pos := int(c.Shields[0]) - int(c.Shields[1])
	if pos > 9 {
		pos = 9
	} else if pos < -9 {
		pos = -9
	}
	return pos
}

func (c *ConflictPawn) add(i PlayerIndex, s Shields) {
	c.Shields[i] += s
}

// --- TODO ---

// Military board
type Military struct {
	ConflictPawn ConflictPawn

	// Military Tokens.
	// "true" means "is taken".
	Tokens2 [numPlayers]bool
	Tokens5 [numPlayers]bool
}

func (m *Military) addShields(g *Game, i PlayerIndex, s Shields) {
	m.ConflictPawn.add(i, s)

	var pos = m.ConflictPawn.Position()
	var lagging PlayerIndex
	if pos > 0 {
		lagging = 1
	} else if pos < 0 {
		lagging = 0
		pos = -pos
	}

	if pos >= 3 && !m.Tokens2[lagging] {
		g.player(lagging).Coins.sub(2)
		m.Tokens2[lagging] = true
	}
	if pos >= 6 && !m.Tokens5[lagging] {
		g.player(lagging).Coins.sub(5)
		m.Tokens5[lagging] = true
	}
	if pos >= 9 {
		g.victory(lagging.Next().winner(), MilitarySupremacy)
	}
}

func (m Military) VP(i PlayerIndex) VP {
	var dt = m.ConflictPawn.Position()

	var leader PlayerIndex
	if dt > 0 {
		leader = 0
	} else {
		leader = 1
		dt = -dt
	}
	if leader != i {
		return 0
	}

	switch {
	case dt == 0:
		return 0
	case dt <= 2:
		return 2
	case dt <= 5:
		return 5
	default:
		return 10
	}
}

// Shields - military power
type Shields uint8

var _ Effect = Shields(0)

func (s Shields) applyEffect(g *Game, i PlayerIndex) {
	g.military.addShields(g, i, s)
}

func (s Shields) Sub(s2 Shields) Shields {
	if s <= s2 {
		return 0
	}
	return s - s2
}
