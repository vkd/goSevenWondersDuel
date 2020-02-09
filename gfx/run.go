package gfx

import (
	"fmt"
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
	Desk BoardState = iota
	Wonders
)

func run() error {
	gg, err := core.NewGame(core.WithSeed(0))
	if err != nil {
		return err
	}

	wonders, _, ok := gg.Init()
	if !ok {
		return fmt.Errorf("cannot init game")
	}
	currentWonder = 0

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
	statsRPlayer := text.New(pixel.V(windowWidth-230, 100), atlas)
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

	var fps = time.Tick(time.Second / 15)

	var boardState BoardState

	for i, idx := range [8]int{3, 0, 1, 2, 5, 4, 7, 6} {
		wonderTaken[i] = true
		wonderChosen[i] = wonders[idx]
		currentWonder++
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

		if win.Pressed(pixelgl.KeyW) {
			boardState = Wonders
		}
		if win.Pressed(pixelgl.KeyD) {
			boardState = Desk
		}

		if win.JustPressed(pixelgl.KeyQ) || win.JustPressed(pixelgl.KeyEscape) {
			win.SetClosed(true)
		}

		mouse := win.MousePosition()

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			if selectedCardIndex > -1 {
				tableCards.Cards, err = gg.Build(tableCards.Cards[selectedCardIndex].ID)
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
		}

		// showCard = win.Pressed(pixelgl.KeyLeftShift) || win.Pressed(pixelgl.MouseButtonLeft)

		// drawCard(0, pixel.V(510, 10), win)
		// drawFirstEpoh(win, pixel.V(windowWidth/2, windowHeight-100))

		war.Clear()
		war.Color = colornames.Red
		warText := []string{
			"9 . . 6 . . 3 . . 0 . . 3 . . 6 . . 9\n",
			strings.Repeat(" ", int(0)*2+18) + "*\n", // 18 = 0
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
		switch boardState {
		case Desk:
			for i, c := range tableCards.Cards {
				if !c.Exists {
					continue
				}
				drawCard(c, tableCards.Rects[i], win)
				if !c.Covered && tableCards.Rects[i].Contains(mouse) {
					selectedCardIndex = i
				}
				idx := text.New(tableCards.Rects[i].Max, atlas)
				idx.Color = colornames.Lightgreen
				idx.Dot.X -= idx.BoundsOf(strconv.Itoa(i)).W()
				idx.Dot.Y -= 10
				idx.WriteString(strconv.Itoa(i))
				idx.Draw(win, pixel.IM)
			}

			if selectedCardIndex > -1 {
				drawCardBorder(tableCards.Rects[selectedCardIndex], win)
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
				drawCardBorder(wonderRects[selectedWonderIndex%4], win)
			}
		}

		txt.Clear()
		txt.Color = colornames.Orange
		fmt.Fprintf(txt,
			"Left: %d\nBottom: %d\nTitle: %d\nDelta: %d\n\nMouse (%d;%d)\nActive player: %d\nWonderTaken: %v\nWonderChosen: %v\nCurrent wonder: %d",
			int(left),
			int(bottom),
			int(cardTitleHeight),
			int(deltaEpoh),
			int(win.MousePosition().X),
			int(win.MousePosition().Y),
			0,
			wonderTaken,
			wonderChosen,
			currentWonder,
		)
		txt.Draw(win, pixel.IM)

		statsLPlayer.Clear()
		fmt.Fprintf(statsLPlayer, debugPlayerInfo(gg.Player(0)))
		statsLPlayer.Draw(win, pixel.IM)

		statsRPlayer.Clear()
		fmt.Fprintf(statsRPlayer, debugPlayerInfo(gg.Player(1)))
		statsRPlayer.Draw(win, pixel.IM)

		win.Update()

		<-fps
	}
	return nil
}

func debugPlayerInfo(p core.Player) string {
	return fmt.Sprintf("Money: %d\nResources: %v\nVP: %d\nChains: %v\nScience: %v", p.Coins, p.Resources, p.VP, p.Chains.Strings(), p.ScientificSymbols)
}

var (
	selectedCardIndex   int = -1
	selectedWonderIndex int = -1
)

type TableCards struct {
	Rects []pixel.Rect
	Cards core.CardsState
}

var (
	showCard bool

	atlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)
)

func drawCard(c core.CardState, r pixel.Rect, win pixel.Target) {
	if !c.Exists {
		return
	}

	im := cardIM
	if !showCard {
		im = im.Scaled(pixel.ZV, 0.5).Moved(pixel.V(cardWidth/2, cardHeight/2))
	} else {
		im = im.Moved(pixel.V(cardWidth, cardHeight))
	}
	if c.FaceUp {
		cardsTx[c.ID].Draw(win, im.Moved(r.Min))
	} else {
		cardsTxBack[0].Draw(win, im.Moved(r.Min))
	}
	// if r.Contains(mouse) && c.IsOnTop() {
	// 	drawCardBorder(c, win)
	// }

	txt := text.New(pixel.V(r.Min.X, r.Max.Y-10), atlas)
	txt.Color = colornames.Lightblue
	if c.FaceUp {
		fmt.Fprintf(txt, "Index: %d", c.ID)
	}
	txt.Draw(win, pixel.IM)
}

func drawWonder(win pixel.Target, id core.WonderID, rect pixel.Rect) {
	wondersTx[id].Draw(win, pixel.IM.Scaled(pixel.ZV, 0.5).Moved(pixel.V(wonderWidth/2, wonderHeight/2)).Moved(rect.Min))
}

func drawCardBorder(c pixel.Rect, win pixel.Target) {
	imd := imdraw.New(nil)
	imd.Color = colornames.Yellow
	imd.EndShape = imdraw.RoundEndShape
	imd.Push(c.Min, c.Max)
	imd.Rectangle(4)
	imd.Draw(win)
}

var (
	cardIM = pixel.IM //.Moved(pixel.V(cardWidth/2, cardHeight/2))
)
