package gfx

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/vkd/goSevenWondersDuel/core"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

const (
	// cardWidth      float64 = 264
	cardWidth float64 = 132
	// cardHeight     float64 = 400
	cardHeight float64 = 200

	wonderWidth  float64 = 294
	wonderHeight float64 = 191

	cardTitleHeight float64 = 70
	deltaEpoh       float64 = 30

	windowWidth  float64 = 1200
	windowHeight float64 = 800

	progressWidth float64 = 159
)

var (
	currentWonder  int
	wonderToPlayer = [8]core.PlayerIndex{0, 1, 1, 0, 1, 0, 0, 1}
	wonderTaken    [8]bool
	wonderBuilt    [2][4]uint8
	ptokens        []core.PTokenID

	topCenter = pixel.V(windowWidth/2, windowHeight-100)

	ageIRects   = genCardRects(ageGrid.genAgeIVecs(topCenter.Sub(pixel.V(0, cardHeight))))
	ageIIRects  = genCardRects(ageGrid.genAgeIIVecs(topCenter.Sub(pixel.V(0, cardHeight))))
	ageIIIRects = genCardRects(ageGrid.genAgeIIIVecs(topCenter.Sub(pixel.V(0, cardHeight))))
)

var discardedRects []pixel.Rect

func init() {
	var dx float64 = 10
	var dy float64 = dx
	for j := 0; j < 3; j++ {
		var y float64 = (windowHeight+dy)/2 - float64(j)*(cardHeight+dy)
		for i := 0; i < 8; i++ {
			var x float64 = 50 + (cardWidth+dx)*float64(i)
			discardedRects = append(discardedRects, pixel.R(x, y, x+cardWidth, y+cardHeight))
		}
	}
}

func Run() error {
	pixelgl.Run(func() {
		err := run()
		if err != nil {
			panic(err)
		}
	})
	return nil
}

type BoardState uint8

const (
	NoneState BoardState = iota
	Table
	Wonders
	DiscardedCards
	PTokensState
	DiscardedPTokensState
	OpponentsDestroyState
	PlayerCardsState
	OpponentCardsState
)

func run() error { //nolint: gocognit, funlen, gocyclo
	g, err := core.NewGame(core.WithSeed(0))
	if err != nil {
		return err
	}
	var w0 = [4]core.WonderID{0, 8, 10, 11}
	var w1 = [4]core.WonderID{5, 6, 1, 4}
	err = g.SelectWonders(w0, w1)
	if err != nil {
		return err
	}

	wonders := g.GetAvailableWonders()
	ptokens = g.GetAvailablePTokens()
	currentWonder = 0

	var discardedCards []core.CardID

	cfg := pixelgl.WindowConfig{
		Title:  "7 Wonders",
		Bounds: pixel.R(0, 0, windowWidth, windowHeight),
		// VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		log.Fatalf("Error on create window: %v", err)
	}

	err = loadTextures()
	if err != nil {
		log.Fatalf("Error on load textures: %v", err)
	}

	txt := text.New(pixel.V(80, windowHeight-30), atlas)
	war := text.New(pixel.V(windowWidth/2, windowHeight-30), atlas)
	statsLPlayer := text.New(pixel.V(30, 100), atlas)
	statsLPlayer.Color = colornames.Yellow
	statsRPlayer := text.New(pixel.V(windowWidth-330, 100), atlas)
	statsRPlayer.Color = colornames.Yellow

	// * - topCenter
	// X - topCenter.Y -= cardHeight
	// O - cards grid
	//
	// y
	// ^
	// |  +---+ +---+ * +---+ +---+
	// |  |   | |   |   |   | |   |
	// |  |  ++-++  |   |  ++-++  |
	// |  O--+  O+--+ X O--+  O+--+
	// |     |   |         |   |
	// |     O---+         O---+
	// |
	// +------------------------------> x
	// topCenter.Y -= cardHeight
	var tableCards = TableCards{
		Cards: g.CardsState(),
		Rects: ageIRects,
	}

	var wonderRects []pixel.Rect
	{
		var dx float64 = 10
		var dy float64 = dx
		var x float64 = windowWidth/2 - wonderWidth - dx/2
		var y float64 = 100
		for i := 0; i < 2; i++ {
			wonderRects = append(wonderRects, pixel.R(x, y, x+wonderWidth, y+wonderHeight))
			wonderRects = append(wonderRects, pixel.R(x+wonderWidth+dx, y, x+wonderWidth+dx+wonderWidth, y+wonderHeight))
			y += wonderHeight + dy
		}
	}

	var ptokenCircles []pixel.Circle
	{
		var width float64 = progressWidth
		var radius = width / 2
		var dx float64 = 30
		var x = windowWidth/2 - 2*(width+dx)
		var y = windowHeight / 2
		for i := 0; i < 5; i++ {
			ptokenCircles = append(ptokenCircles, pixel.C(pixel.V(x, y), radius))
			x += width + dx
		}
	}

	var userWonders [2][]core.WonderID = [2][]core.WonderID{w0[:], w1[:]}
	var userWonderRects [2][4]pixel.Rect
	{
		for i := 0; i < 4; i++ {
			var y = float64(i) * wonderHeight
			userWonderRects[0][i] = pixel.R(
				0,
				y,
				wonderWidth,
				y+wonderHeight,
			)
			userWonderRects[1][i] = pixel.R(
				windowWidth-wonderWidth,
				y,
				windowWidth,
				y+wonderHeight,
			)
		}
	}

	var discardedPTokens []core.PTokenID
	var opponentsDestroyableBuildings []core.CardID
	var playerCards [2][core.CardColorSize][]core.CardID

	var fps = time.NewTicker(time.Second / 15)
	defer fps.Stop()

	var boardState BoardState = Table

	var verbose bool = true
	var currentAge = 1
	var wondersBuilt int = 0

	// for i, idx := range [8]int{3, 0, 1, 2, 5, 4, 7, 6} {
	// 	wonderTaken[i] = true
	// 	wonderChosen[i] = wonders[idx]
	// 	currentWonder++
	// }

	var nextTurn = func() {
		var isNextAge bool = true
		for _, cs := range tableCards.Cards {
			if cs.Exists {
				isNextAge = false
				break
			}
		}
		if isNextAge {
			currentAge++
			switch currentAge {
			case 2:
				tableCards.Rects = ageIIRects
			case 3:
				tableCards.Rects = ageIIIRects
			case 4:
				w, r, vps, _ := g.VictoryResult()
				log.Printf("w: %v, r: %v, vps: %v", w, r, vps)
			}
			tableCards.Cards = g.CardsState()
		}

		switch g.GetState() {
		case core.StateGameTurn:
			boardState = Table
		case core.StateChooseFirstPlayer:
			boardState = Table
		case core.StateBuildFreePToken:
			discardedPTokens, err = g.GetDiscardedPTokens()
			if err != nil {
				log.Printf("Error on next turn (build free ptoken): %v", err)
			}
			boardState = DiscardedPTokensState
		case core.StateDiscardOpponentBuild:
			opponentsDestroyableBuildings, err = g.GetDiscardedOpponentsBuildings()
			if err != nil {
				log.Printf("Error on next turn (discard opponents building): %v", err)
			}
			boardState = OpponentsDestroyState
		case core.StateBuildFreeDiscarded:
			boardState = DiscardedCards
		case core.StateChoosePToken:
			boardState = PTokensState
		case core.StateVictory:
			w, reason, vps, err := g.VictoryResult()
			if err != nil {
				log.Printf("Error on get victory result: %v", err)
			}
			log.Printf("w: %v, reason: %v, vps: %v", w, reason, vps)
			for i := 0; i < 2; i++ {
				log.Printf("Player: %d", i)
				for tp, vp := range vps[i] {
					log.Printf("%s: %d", core.VPType(tp).String(), vp)
				}
				log.Printf("---")
			}
		}
	}

	var (
		selectedCardIndex       int = -1
		selectedWonderIndex     int = -1
		selectedUserWonderIndex int = -1
		selectedConstructWonder int = -1
		selectedDiscardedIndex  int = -1
		selectedPTokenIndex     int = -1
		selectDiscardedPToken   int = -1
		selectOpponentsBuilding int = -1
	)

	for !win.Closed() {
		var pIndex = g.CurrentPlayerIndex()
		var gameState = g.GetState()

		win.Clear(colornames.Purple)

		// if win.Pressed(pixelgl.KeyUp) {
		// }
		// if win.Pressed(pixelgl.KeyDown) {
		// }
		// if win.Pressed(pixelgl.KeyLeft) {
		// }
		// if win.Pressed(pixelgl.KeyRight) {
		// }

		if win.JustPressed(pixelgl.KeyF12) {
			g, err = core.NewGame(core.WithSeed(0))
			if err != nil {
				return err
			}
			wonders = g.GetAvailableWonders()
			ptokens = g.GetAvailablePTokens()
			tableCards.Cards = g.CardsState()
			tableCards.Rects = ageIRects
			discardedCards = nil
			currentAge = 1

			currentWonder = 0
			userWonders = [2][]core.WonderID{}
			wonderTaken = [8]bool{}
			boardState = Wonders
			wonderBuilt = [2][4]uint8{}
			playerCards = [2][core.CardColorSize][]core.CardID{}
		}

		if win.JustPressed(pixelgl.KeyW) {
			boardState = Wonders
		}
		if win.JustPressed(pixelgl.KeyT) {
			boardState = Table
		}
		if win.JustPressed(pixelgl.KeyD) {
			boardState = DiscardedCards
		}
		if win.JustPressed(pixelgl.KeyV) {
			verbose = !verbose
		}
		if win.JustPressed(pixelgl.KeyP) {
			boardState = PTokensState
		}
		if win.JustPressed(pixelgl.KeyH) {
			boardState = PlayerCardsState
		}
		if win.JustPressed(pixelgl.KeyO) {
			boardState = OpponentCardsState
		}
		if win.JustPressed(pixelgl.Key1) {
			err = g.ChooseFirstPlayer(0)
			if err != nil {
				log.Printf("Error on choose first player: %v", err)
			}
		}
		if win.JustPressed(pixelgl.Key2) {
			err = g.ChooseFirstPlayer(1)
			if err != nil {
				log.Printf("Error on choose first player: %v", err)
			}
		}

		if win.JustPressed(pixelgl.KeyEscape) {
			win.SetClosed(true)
		}

		mouse := win.MousePosition()

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			if selectedCardIndex > -1 {
				if selectedConstructWonder > -1 {
					tableCards.Cards, err = g.ConstructWonder(tableCards.Cards[selectedCardIndex].ID, userWonders[g.CurrentPlayerIndex()][selectedConstructWonder])

					wonderBuilt[pIndex][selectedConstructWonder] = uint8(tableCards.Cards[selectedCardIndex].ID) + 1
					wondersBuilt++
				} else {
					cardID := tableCards.Cards[selectedCardIndex].ID
					tableCards.Cards, err = g.ConstructBuilding(cardID)
					if err == nil {
						playerCards[pIndex][cardID.Color()] = append(playerCards[pIndex][cardID.Color()], cardID)
					}
				}
				if err != nil {
					log.Printf("Error on build: %v", err)
				} else {
					selectedConstructWonder = -1
				}
			}
			if selectedWonderIndex > -1 {
				if !wonderTaken[selectedWonderIndex] {
					wonderTaken[selectedWonderIndex] = true
					// wonderChosen[currentWonder] = wonders[selectedWonderIndex]
					userWonders[wonderToPlayer[currentWonder]] = append(userWonders[wonderToPlayer[currentWonder]], wonders[selectedWonderIndex])

					currentWonder++
					if currentWonder >= 8 {
						{
							var fst, snd [4]core.WonderID
							for i, id := range userWonders[0] {
								fst[i] = id
							}
							for i, id := range userWonders[1] {
								snd[i] = id
							}

							err = g.SelectWonders(fst, snd)
							if err != nil {
								return err
							}
						}
						boardState = Table
					}
				}
			}
			if selectedDiscardedIndex > -1 {
				cardID := discardedCards[selectedDiscardedIndex]
				err = g.ConstructDiscardedCard(cardID)
				if err != nil {
					log.Printf("Error on build discarded: %v", err)
				} else {
					playerCards[pIndex][cardID.Color()] = append(playerCards[pIndex][cardID.Color()], cardID)
				}
			}
			if selectedUserWonderIndex > -1 {
				if selectedConstructWonder == selectedUserWonderIndex {
					selectedConstructWonder = -1
				} else {
					selectedConstructWonder = selectedUserWonderIndex % 4
				}
			}
			if selectedPTokenIndex > -1 {
				err = g.ChoosePToken(ptokens[selectedPTokenIndex])
				if err != nil {
					log.Printf("Error on choose PToken: %v", err)
				}
				ptokens = g.GetAvailablePTokens()
			}
			if selectDiscardedPToken > -1 {
				err = g.PlayDiscardedPToken(discardedPTokens[selectDiscardedPToken])
				if err != nil {
					log.Printf("Error on play discarded PTokens: %v", err)
				} else {
					discardedPTokens = nil
					boardState = Table
				}
			}
			if selectOpponentsBuilding > -1 {
				cardID := opponentsDestroyableBuildings[selectOpponentsBuilding]
				err = g.DiscardOpponentBuild(cardID)
				if err != nil {
					log.Printf("Error on discard opponent building: %v", err)
				} else {
					var newList []core.CardID
					for _, cid := range playerCards[(pIndex+1)%2][cardID.Color()] {
						if cid != cardID {
							newList = append(newList, cid)
						}
					}
					playerCards[(pIndex+1)%2][cardID.Color()] = newList
					opponentsDestroyableBuildings = nil
					boardState = Table
				}
			}
			nextTurn()
		} else if win.JustPressed(pixelgl.MouseButtonRight) {
			if selectedCardIndex > -1 {
				var err error
				var id = tableCards.Cards[selectedCardIndex].ID
				tableCards.Cards, err = g.DiscardCard(id)
				if err != nil {
					log.Printf("Card %d is not discarded: %v", id, err)
				} else {
					discardedCards = append(discardedCards, id)
				}
			}
			nextTurn()
		}

		// showCard = win.Pressed(pixelgl.KeyLeftShift) || win.Pressed(pixelgl.MouseButtonLeft)

		// drawCard(0, pixel.V(510, 10), win)
		// drawFirstEpoh(win, pixel.V(windowWidth/2, windowHeight-100))

		m := g.Military()
		var warPoint = m.ConflictPawn.Position()

		war.Clear()
		war.Color = colornames.Red
		warText := []string{
			"9 . . 6 . . 3 . . 0 . . 3 . . 6 . . 9\n",
			strings.Repeat(" ", warPoint*2+18) + "*\n", // 18 = 0
			"X    -5    -2           -2    -5    X\n",
			"  10     5     2     2     5     10  \n",
		}
		align := war.BoundsOf(warText[0]).W() / 2
		for _, text := range warText {
			war.Dot.X -= align
			war.WriteString(text) // nolint
		}
		war.Draw(win, pixel.IM)

		selectedCardIndex = -1
		selectedWonderIndex = -1
		selectedUserWonderIndex = -1
		selectedDiscardedIndex = -1
		selectedPTokenIndex = -1
		selectDiscardedPToken = -1
		selectOpponentsBuilding = -1

		switch boardState {
		case Table:
			for i, c := range tableCards.Cards {
				if !c.Exists {
					continue
				}
				if !c.Covered && tableCards.Rects[i].Contains(mouse) {
					selectedCardIndex = i
				}
			}
			for i, c := range tableCards.Cards {
				if !c.Exists {
					continue
				}
				cost := g.CardCostCoins(tableCards.Cards[i].ID)
				drawCard(c.ID, c.FaceUp, tableCards.Rects[i], win, cost)
				if i == selectedCardIndex {
					drawSelectedBorder(tableCards.Rects[i], win)
				}
				if !c.Covered {
					var color = colornames.Red
					if g.CurrentPlayer().Coins >= cost {
						color = colornames.Green
					}
					drawBorder(tableCards.Rects[i], win, color, 2)
				}
				idx := text.New(tableCards.Rects[i].Max, atlas)
				idx.Color = colornames.Lightgreen
				idx.Dot.X -= idx.BoundsOf(strconv.Itoa(i)).W()
				idx.Dot.Y -= 10
				idx.WriteString(strconv.Itoa(i)) // nolint
				idx.Draw(win, pixel.IM)
			}

		case Wonders:
			if gameState == core.StateSelectWonders {
				// --- draft
				var second int = 0
				if currentWonder >= 4 {
					second = 4
				}
				for i, r := range wonderRects {
					i += second
					if wonderTaken[i] {
						continue
					}
					drawWonder(win, true, wonders[i], r, 0, 0)
					if r.Contains(mouse) {
						selectedWonderIndex = i
					}
				}
				if selectedWonderIndex >= 0 {
					drawSelectedBorder(wonderRects[selectedWonderIndex%4], win)
				}
			}
			for pi, ws := range userWonders {
				for wi, id := range ws {
					r := userWonderRects[pi][wi]

					cs, trade := g.WonderCost(id)
					cost := cs + trade
					if wonderBuilt[pi][wi] > 0 {
						drawWonder(win, true, id, r, cost, wonderBuilt[pi][wi])
					} else {
						drawWonder(win, wondersBuilt < 7, id, r, cost, wonderBuilt[pi][wi])
					}

					if core.PlayerIndex(pi) == pIndex && wonderBuilt[pi][wi] == 0 && r.Contains(mouse) {
						selectedUserWonderIndex = wi + pi*4
					}
				}
			}
			if selectedConstructWonder > -1 {
				r := userWonderRects[pIndex][selectedConstructWonder]
				drawBorder(r, win, colornames.Aqua, 4)
			}
			if selectedUserWonderIndex >= 0 {
				r := userWonderRects[selectedUserWonderIndex/4][selectedUserWonderIndex%4]
				drawSelectedBorder(r, win)
			}
			for wi, id := range userWonders[pIndex] {
				if wonderBuilt[pIndex][wi] != 0 {
					continue
				}
				r := userWonderRects[pIndex][wi]

				cs, trade := g.WonderCost(id)
				cost := cs + trade

				var color = colornames.Red
				if g.CurrentPlayer().Coins >= cost {
					color = colornames.Green
				}
				drawBorder(r, win, color, 2)
			}
		case DiscardedCards:
			for i, did := range discardedCards {
				drawCard(did, true, discardedRects[i], win, 0)
				if discardedRects[i].Contains(mouse) {
					selectedDiscardedIndex = i
				}
			}
			if selectedDiscardedIndex >= 0 {
				drawSelectedBorder(discardedRects[selectedDiscardedIndex], win)
			}
		case PTokensState:
			for i, pt := range ptokens {
				drawPToken(win, pt, ptokenCircles[i])
				if ptokenCircles[i].Contains(mouse) {
					selectedPTokenIndex = i
				}
			}
			if selectedPTokenIndex > -1 {
				drawCircle(win, ptokenCircles[selectedPTokenIndex], borderColor, 4)
			}
		case DiscardedPTokensState:
			for i, pi := range discardedPTokens {
				drawPToken(win, pi, ptokenCircles[i])
				if ptokenCircles[i].Contains(mouse) {
					selectDiscardedPToken = i
				}
			}
			if selectDiscardedPToken > -1 {
				drawCircle(win, ptokenCircles[selectDiscardedPToken], borderColor, 4)
			}
		case OpponentsDestroyState:
			for i, cid := range opponentsDestroyableBuildings {
				drawCard(cid, true, discardedRects[i], win, 0)
				if discardedRects[i].Contains(mouse) {
					selectOpponentsBuilding = i
				}
			}
			if selectOpponentsBuilding > -1 {
				drawBorder(discardedRects[selectOpponentsBuilding], win, borderColor, 4)
			}
		case PlayerCardsState:
			drawHand(win, playerCards[pIndex])
		case OpponentCardsState:
			drawHand(win, playerCards[pIndex.Next()])
		}

		var currectPlayer = g.CurrentPlayerIndex()

		// debug info
		if verbose {
			txt.Clear()
			txt.Color = colornames.Orange
			fmt.Fprintf(txt,
				"Active player: %d\n"+
					"State: %s\n"+
					"Age: %d\n",
				currectPlayer,
				gameState.String(),
				currentAge,
			)
			txt.Draw(win, pixel.IM)

			statsLPlayer.Clear()
			fmt.Fprint(statsLPlayer, debugPlayerInfo(g.Player(0), currectPlayer == 0))
			statsLPlayer.Draw(win, pixel.IM)

			statsRPlayer.Clear()
			fmt.Fprint(statsRPlayer, debugPlayerInfo(g.Player(1), currectPlayer == 1))
			statsRPlayer.Draw(win, pixel.IM)
		}

		win.Update()

		<-fps.C
	}
	return nil
}

func drawHand(t pixel.Target, hand [core.CardColorSize][]core.CardID) {
	for color, stack := range hand {
		for i, cid := range stack {
			drawCard(cid, true, discardedRects[color].Moved(pixel.V(0, -50*float64(i))), t, 0)
		}
	}
}

func debugPlayerInfo(p core.Player, isActive bool) string {
	var active string
	if isActive {
		active = "*"
	}
	return fmt.Sprintf(
		"IsActive: %s\n"+
			"Money: %d\n"+
			"         : [W S C P G]\n"+
			"Resources: %v\n"+
			"Chains: %v\n"+
			"       : [W M C T P A S]\n"+
			"Science: %v",
		active,
		p.Coins,
		p.Resources,
		p.Chains.Strings(),
		p.ScientificSymbols,
	)
}

type TableCards struct {
	Rects []pixel.Rect
	Cards core.CardsState
}

var (
	atlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)

	cardIM = pixel.IM.Scaled(pixel.ZV, cardWidth/textureCardWidth) //.Moved(pixel.V(cardWidth/2, cardHeight/2))

	borderColor = colornames.Yellow
)

func drawCard(id core.CardID, faceUp bool, r pixel.Rect, win pixel.Target, cost core.Coins) {
	im := cardIM //.Moved(pixel.V(cardWidth/2, cardHeight/2))
	var t = cardsTxBack[0]
	if faceUp {
		t = cardsTx[id]
	}
	t.Draw(win, im.Moved(r.Center()))

	txt := text.New(pixel.V(r.Min.X, r.Max.Y-10), atlas)
	switch id.Color() {
	case core.Yellow:
		txt.Color = colornames.Black
	default:
		txt.Color = colornames.White
	}
	if faceUp {
		fmt.Fprintf(txt, "Index: %d\nCost: %d", id, cost)
	}
	txt.Draw(win, pixel.IM)
}

func drawWonder(win pixel.Target, faceUp bool, id core.WonderID, rect pixel.Rect, cost core.Coins, builtCardID uint8) {
	var tx = wondersTxBack
	if faceUp {
		tx = wondersTx[id]
	}
	tx.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.5).Moved(pixel.V(wonderWidth/2, wonderHeight/2)).Moved(rect.Min))

	if builtCardID > 0 {
		var t = cardsTx[builtCardID-1]
		t.Draw(win, cardIM.Rotated(pixel.ZV, -math.Pi/2).Moved(rect.Center()))
		return
	}

	if faceUp {
		txt := text.New(pixel.V(rect.Min.X, rect.Max.Y-10), atlas)
		txt.Color = colornames.White
		fmt.Fprintf(txt, "Index: %d\nCost: %d", id, cost)
		txt.Draw(win, pixel.IM)
	}
}

func drawPToken(win pixel.Target, id core.PTokenID, c pixel.Circle) {
	progressTx[id].Draw(win, pixel.IM.Moved(c.Center))
}

func drawSelectedBorder(c pixel.Rect, win pixel.Target) {
	drawBorder(c, win, colornames.Yellow, 4)
}

func drawBorder(c pixel.Rect, win pixel.Target, color color.RGBA, thickness float64) {
	imd := imdraw.New(nil)
	imd.Color = color
	imd.EndShape = imdraw.RoundEndShape
	imd.Push(c.Min, c.Max)
	imd.Rectangle(thickness)
	imd.Draw(win)
}

func drawCircle(t pixel.Target, c pixel.Circle, color color.RGBA, thickness float64) {
	circle := imdraw.New(nil)
	circle.Color = color
	circle.Push(c.Center)
	circle.Circle(c.Radius, thickness)
	circle.Draw(t)
}
