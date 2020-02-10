package core

// Military board
type Military struct {
	// Conflict pawn
	shields [numPlayers]Shields

	// Military tokens
	tokens2 [numPlayers]bool
	tokens5 [numPlayers]bool
}

func (m *Military) addShields(g *Game, i PlayerIndex, s Shields) {
	m.shields[i] += s

	var diff = m.diffFor(i)

	if diff >= 3 && m.tokens2[i] == false {
		g.player(i).Coins.sub(2)
		m.tokens2[i] = true
	}
	if diff >= 6 && m.tokens5[i] == false {
		g.player(i).Coins.sub(5)
		m.tokens5[i] = true
	}
	if diff >= 10 {
		g.Victory()
	}
}

func (m *Military) diffFor(i PlayerIndex) Shields {
	if m.shields[i] <= m.shields[i.Next()] {
		return 0
	}
	return m.shields[i] - m.shields[i.Next()]
}

// Shields - military power
type Shields uint8

var _ Effect = Shields(0)

func (s Shields) applyEffect(g *Game, i PlayerIndex) {
	g.military.addShields(g, i, s)
}
