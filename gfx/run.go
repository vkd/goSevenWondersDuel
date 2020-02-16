package gfx

import (
	"fmt"
	"image/color"
	"log"
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
	leftPaddingCards   float64 = 89
	bottomPaddingCards float64 = 74

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

	progressWidth  float64 = 159
	progressHeight float64 = 159
)

var (
	currentWonder  int
	wonderToPlayer = [8]core.PlayerIndex{0, 1, 1, 0, 1, 0, 0, 1}
	wonderTaken    [8]bool

	wonderChosen [8]core.WonderID
)

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
	Table BoardState = iota
	Wonders
	DiscardedCards
)

func run() error {
	g, err := core.NewGame(core.WithSeed(0))
	if err != nil {
		return err
	}
	var gg = g

	wonders, _, ok := gg.Init()
	if !ok {
		return fmt.Errorf("cannot init game")
	}
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

	txt := text.New(pixel.V(30, windowHeight-30), atlas)
	war := text.New(pixel.V(windowWidth/2, windowHeight-30), atlas)
	statsLPlayer := text.New(pixel.V(30, 100), atlas)
	statsLPlayer.Color = colornames.Yellow
	statsRPlayer := text.New(pixel.V(windowWidth-330, 100), atlas)
	statsRPlayer.Color = colornames.Yellow

	var left float64 = leftPaddingCards
	var bottom float64 = bottomPaddingCards

	var topCenter = pixel.V(windowWidth/2, windowHeight-100)

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
		Cards: gg.CardsState(),
		Rects: genCardRects(ageGrid.genAgeIVecs(topCenter.Sub(pixel.V(0, cardHeight)))),
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

	var discardedRects []pixel.Rect
	{
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

	var fps = time.Tick(time.Second / 15)

	var boardState BoardState

	for i, idx := range [8]int{3, 0, 1, 2, 5, 4, 7, 6} {
		wonderTaken[i] = true
		wonderChosen[i] = wonders[idx]
		currentWonder++
	}

	{
		var fst, snd [4]core.WonderID
		var i, j int
		for idx, wc := range wonderChosen {
			switch wonderToPlayer[idx] {
			case 0:
				fst[i] = wc
				i++
			case 1:
				snd[j] = wc
				j++
			}
		}

		err = g.SelectWonders(fst, snd)
		if err != nil {
			return err
		}
	}

	// var minX, minY, minW, minH float64
	// minW, minH = wonderWidth*2, wonderHeight*2
	// minX, minY = wonderLefts[1], wonderBottoms[1]

	for !win.Closed() {
		win.Clear(colornames.Purple)

		if win.Pressed(pixelgl.KeyUp) {
			// if win.Pressed(pixelgl.KeyLeftShift) {
			// 	// minH += 1
			// 	minY += minH
			// } else {
			// 	minY += 1
			// }
		}
		if win.Pressed(pixelgl.KeyDown) {
			// if win.Pressed(pixelgl.KeyLeftShift) {
			// 	// minH -= 1
			// 	minY -= minH
			// } else {
			// 	minY -= 1
			// }
		}
		if win.Pressed(pixelgl.KeyLeft) {
			// if win.Pressed(pixelgl.KeyLeftShift) {
			// 	// minW -= 1
			// 	minX -= minW
			// } else {
			// 	minX -= 1
			// }
		}
		if win.Pressed(pixelgl.KeyRight) {
			// if win.Pressed(pixelgl.KeyLeftShift) {
			// 	// minW += 1
			// 	minX += minW
			// } else {
			// 	minX += 1
			// }
		}

		if win.JustPressed(pixelgl.KeyR) {
			g, err = core.NewGame(core.WithSeed(0))
			if err != nil {
				return err
			}
			gg = g
			wonders, _, ok = gg.Init()
			if !ok {
				return fmt.Errorf("cannot init game")
			}
			tableCards.Cards = gg.CardsState()
			discardedCards = nil

			currentWonder = 0
			for i, idx := range [8]int{3, 0, 1, 2, 5, 4, 7, 6} {
				wonderTaken[i] = true
				wonderChosen[i] = wonders[idx]
				currentWonder++
			}

			{
				var fst, snd [4]core.WonderID
				var i, j int
				for idx, wc := range wonderChosen {
					switch wonderToPlayer[idx] {
					case 0:
						fst[i] = wc
						i++
					case 1:
						snd[j] = wc
						j++
					}
				}

				err = g.SelectWonders(fst, snd)
				if err != nil {
					return err
				}
			}
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

		if win.JustPressed(pixelgl.KeyQ) || win.JustPressed(pixelgl.KeyEscape) {
			win.SetClosed(true)
		}

		mouse := win.MousePosition()

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			if selectedCardIndex > -1 {
				tableCards.Cards, err = gg.ConstructBuilding(tableCards.Cards[selectedCardIndex].ID)
				if err != nil {
					log.Printf("Error on build: %v", err)
				}
			}
			if selectedWonderIndex > -1 {
				if !wonderTaken[selectedWonderIndex] {
					wonderTaken[selectedWonderIndex] = true
					wonderChosen[currentWonder] = wonders[selectedWonderIndex]
					currentWonder++
				}
			}
			if selectedDiscardedIndex > -1 {
				err = gg.ConstructDiscardedCard(discardedCards[selectedDiscardedIndex])
				if err != nil {
					log.Printf("Error on build discarded: %v", err)
				}
			}
		} else if win.JustPressed(pixelgl.MouseButtonRight) {
			if selectedCardIndex > -1 {
				var err error
				var id = tableCards.Cards[selectedCardIndex].ID
				tableCards.Cards, err = gg.DiscardCard(id)
				if err != nil {
					log.Printf("Card %d is not discarded: %v", id, err)
				} else {
					discardedCards = append(discardedCards, id)
				}
			}
		}

		// showCard = win.Pressed(pixelgl.KeyLeftShift) || win.Pressed(pixelgl.MouseButtonLeft)

		// drawCard(0, pixel.V(510, 10), win)
		// drawFirstEpoh(win, pixel.V(windowWidth/2, windowHeight-100))

		m := gg.Military()
		warPoint := int(m.Shields[0]) - int(m.Shields[1])

		war.Clear()
		war.Color = colornames.Red
		warText := []string{
			"9 . . 6 . . 3 . . 0 . . 3 . . 6 . . 9\n",
			strings.Repeat(" ", int(warPoint)*2+18) + "*\n", // 18 = 0
			"X    -5    -2           -2    -5    X\n",
			"  10     5     2     2     5     10  \n",
		}
		align := war.BoundsOf(warText[0]).W() / 2
		for _, text := range warText {
			war.Dot.X -= align
			war.WriteString(text)
		}
		war.Draw(win, pixel.IM)

		selectedCardIndex = -1
		selectedWonderIndex = -1
		selectedDiscardedIndex = -1

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
				drawCard(c.ID, c.FaceUp, tableCards.Rects[i], win)
				if i == selectedCardIndex {
					drawSelectedBorder(tableCards.Rects[i], win)
				}
				if !c.Covered {
					var color = colornames.Red
					if gg.CurrentPlayer().Coins >= gg.CardCostCoins(tableCards.Cards[i].ID) {
						color = colornames.Green
					}
					drawBorder(tableCards.Rects[i], win, color, 2)
				}
				idx := text.New(tableCards.Rects[i].Max, atlas)
				idx.Color = colornames.Lightgreen
				idx.Dot.X -= idx.BoundsOf(strconv.Itoa(i)).W()
				idx.Dot.Y -= 10
				idx.WriteString(strconv.Itoa(i))
				idx.Draw(win, pixel.IM)
			}

		case Wonders:
			var wonder0Y float64 = 0
			var wonder1Y float64 = 0
			for i := 0; i < currentWonder; i++ {
				switch wonderToPlayer[i] {
				case 0:
					drawWonder(win, wonderChosen[i], pixel.R(0, wonder0Y, wonderWidth, wonder0Y+wonderHeight))
					wonder0Y += wonderHeight
				case 1:
					drawWonder(win, wonderChosen[i], pixel.R(windowWidth-wonderWidth, wonder1Y, windowWidth, wonder1Y+wonderHeight))
					wonder1Y += wonderHeight
				}
			}

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
				drawWonder(win, wonders[i], r)
				if r.Contains(mouse) {
					selectedWonderIndex = i
				}
			}
			if selectedWonderIndex >= 0 {
				drawSelectedBorder(wonderRects[selectedWonderIndex%4], win)
			}
		case DiscardedCards:
			for i, did := range discardedCards {
				drawCard(did, true, discardedRects[i], win)
				if discardedRects[i].Contains(mouse) {
					selectedDiscardedIndex = i
				}
			}
			if selectedDiscardedIndex >= 0 {
				drawSelectedBorder(discardedRects[selectedDiscardedIndex], win)
			}
		}

		var currectPlayer = g.CurrentPlayerIndex()

		// debug info
		txt.Clear()
		txt.Color = colornames.Orange
		fmt.Fprintf(txt,
			"Left: %d\nBottom: %d\nTitle: %d\nDelta: %d\n\nMouse (%d;%d)\nActive player: %d\nState: %s",
			int(left),
			int(bottom),
			int(cardTitleHeight),
			int(deltaEpoh),
			int(win.MousePosition().X),
			int(win.MousePosition().Y),
			currectPlayer,
			g.GetState().String(),
		)
		txt.Draw(win, pixel.IM)

		statsLPlayer.Clear()
		fmt.Fprintf(statsLPlayer, debugPlayerInfo(gg.Player(0), currectPlayer == 0))
		statsLPlayer.Draw(win, pixel.IM)

		statsRPlayer.Clear()
		fmt.Fprintf(statsRPlayer, debugPlayerInfo(gg.Player(1), currectPlayer == 1))
		statsRPlayer.Draw(win, pixel.IM)

		win.Update()

		<-fps
	}
	return nil
}

func debugPlayerInfo(p core.Player, isActive bool) string {
	var active string
	if isActive {
		active = "*"
	}
	return fmt.Sprintf("IsActive: %s\nMoney: %d\n         : [W S C P G]\nResources: %v\nChains: %v\n       : [W M C T P A S]\nScience: %v", active, p.Coins, p.Resources, p.Chains.Strings(), p.ScientificSymbols)
}

var (
	selectedCardIndex      int = -1
	selectedWonderIndex    int = -1
	selectedDiscardedIndex int = -1
)

type TableCards struct {
	Rects []pixel.Rect
	Cards core.CardsState
}

var (
	showCard bool

	atlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)
)

func drawCard(id core.CardID, faceUp bool, r pixel.Rect, win pixel.Target) {
	im := cardIM
	if !showCard {
		im = im.Scaled(pixel.ZV, 0.5).Moved(pixel.V(cardWidth/2, cardHeight/2))
	} else {
		im = im.Moved(pixel.V(cardWidth, cardHeight))
	}
	if faceUp {
		cardsTx[id].Draw(win, im.Moved(r.Min))
	} else {
		cardsTxBack[0].Draw(win, im.Moved(r.Min))
	}
	// if r.Contains(mouse) && c.IsOnTop() {
	// 	drawSelectedBorder(c, win)
	// }

	txt := text.New(pixel.V(r.Min.X, r.Max.Y-10), atlas)
	txt.Color = colornames.Lightblue
	if faceUp {
		fmt.Fprintf(txt, "Index: %d", id)
	}
	txt.Draw(win, pixel.IM)
}

func drawWonder(win pixel.Target, id core.WonderID, rect pixel.Rect) {
	wondersTx[id].Draw(win, pixel.IM.Scaled(pixel.ZV, 0.5).Moved(pixel.V(wonderWidth/2, wonderHeight/2)).Moved(rect.Min))
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

var (
	cardIM = pixel.IM //.Moved(pixel.V(cardWidth/2, cardHeight/2))
)
