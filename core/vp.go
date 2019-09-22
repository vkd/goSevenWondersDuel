package core

// VP - victory points
type VP uint8

// FinalVP - extra VP at the end of a game
func (vp VP) FinalVP(g *Game, i PlayerIndex) VP {
	return vp
}
