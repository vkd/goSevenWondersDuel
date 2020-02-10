package core

// Military board
type Military struct {
	// Conflict pawn
	Shields [numPlayers]Shields

	// Military tokens
	Tokens2 [numPlayers]bool
	Tokens5 [numPlayers]bool
}

func (m *Military) addShields(g *Game, i PlayerIndex, s Shields) {
	m.Shields[i] += s

	var diff = m.diffFor(i)

	if diff >= 3 && m.Tokens2[i.Next()] == false {
		g.player(i.Next()).Coins.sub(2)
		m.Tokens2[i.Next()] = true
	}
	if diff >= 6 && m.Tokens5[i.Next()] == false {
		g.player(i.Next()).Coins.sub(5)
		m.Tokens5[i.Next()] = true
	}
	if diff >= 10 {
		g.Victory()
	}
}

func (m *Military) diffFor(i PlayerIndex) Shields {
	if m.Shields[i] <= m.Shields[i.Next()] {
		return 0
	}
	return m.Shields[i] - m.Shields[i.Next()]
}

// Shields - military power
type Shields uint8

var _ Effect = Shields(0)

func (s Shields) applyEffect(g *Game, i PlayerIndex) {
	g.military.addShields(g, i, s)
}
