package game

func (s Shields) effect(g *Game, i PlayerIndex) {
	g.war.Shields[i] += s
	dt, dtI := g.war.dt()
	if dt >= 3 && g.war.firstGte3[dtI] == false {
		g.war.firstGte3[dtI] = true
		DiscardMoney(2).effect(g, i.Next())
	}
	if dt >= 6 && g.war.firstGte6[dtI] == false {
		g.war.firstGte6[dtI] = true
		DiscardMoney(5).effect(g, i.Next())
	}
	if dt >= 9 {
		g.victory(dtI, MilitarySupremacy)
	}
}

// Military ...
type Military struct {
	Shields   [numPlayers]Shields // Conflict pawn
	firstGte3 [numPlayers]bool
	firstGte6 [numPlayers]bool
}

func (w *Military) dt() (Shields, PlayerIndex) {
	res := w.Shields[0] - w.Shields[1]
	if res < 0 {
		return -res, 1
	}
	return res, 0
}

func (w *Military) finalVP() (VP, PlayerIndex) {
	var res VP
	dt, dtI := w.dt()
	if dt >= 1 {
		res = 2
	}
	if dt >= 3 {
		res = 5
	}
	if dt >= 6 {
		res = 10
	}
	return res, dtI
}
