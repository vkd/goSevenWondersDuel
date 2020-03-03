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

// Leader of the board by military score.
//
// Returns leader's PlayerIndex and leader's score.
func (c ConflictPawn) Leader() (PlayerIndex, uint) {
	pos := c.Position()
	if pos > 0 {
		return 0, uint(pos)
	} else if pos < 0 {
		return 1, uint(-pos)
	}
	return 0, 0
}

func (c *ConflictPawn) add(i PlayerIndex, s Shields) {
	c.Shields[i] += s
}

// Board - The board represents the military rivalry between the two cities.
// It is divided into zones (9) and spaces (19). The last space on each end
// represents the player’s capital.
type Board struct {
	ConflictPawn ConflictPawn

	// The Military tokens represent the benefits a city earns when it
	// manages to gain the upper hand, militarily, over its opponent.
	// "true" means "is taken".
	MilitaryTokens2 [numPlayers]bool
	MilitaryTokens5 [numPlayers]bool
}

func (m *Board) addShields(g *Game, i PlayerIndex, s Shields) {
	m.ConflictPawn.add(i, s)

	var leader, val = m.ConflictPawn.Leader()
	var lagging = leader.Next()

	if val >= 3 && !m.MilitaryTokens2[lagging] {
		g.player(lagging).Coins.sub(2)
		m.MilitaryTokens2[lagging] = true
	}
	if val >= 6 && !m.MilitaryTokens5[lagging] {
		g.player(lagging).Coins.sub(5)
		m.MilitaryTokens5[lagging] = true
	}
	if val >= 9 {
		g.victory(lagging.Next().winner(), MilitarySupremacy)
	}
}

func (m Board) VP(i PlayerIndex) VP {
	var leader, val = m.ConflictPawn.Leader()

	if leader != i {
		return 0
	}

	switch {
	case val == 0:
		return 0
	case val <= 2:
		return 2
	case val <= 5:
		return 5
	default:
		return 10
	}
}

// --- TODO ---

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
