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

	var tableCards = TableCards{
		Cards: gg.CardsState(),
		Rects: drawIEpoh(pixel.V(windowWidth/2, windowHeight-100)),
	}
	// log.Printf("%#v", cards[0])

	var fps = time.Tick(time.Second / 15)

	var boardState BoardState

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
		}

		// showCard = win.Pressed(pixelgl.KeyLeftShift) || win.Pressed(pixelgl.MouseButtonLeft)

		// drawCard(0, pixel.V(510, 10), win)
		// drawFirstEpoh(win, pixel.V(windowWidth/2, windowHeight-100))

		txt.Clear()
		txt.Color = colornames.Orange
		fmt.Fprintf(txt, "Left: %d\nBottom: %d\nTitle: %d\nDelta: %d\n\nMouse (%d;%d)\nActive player: %d", int(left), int(bottom), int(cardTitleHeight), int(deltaEpoh), int(win.MousePosition().X), int(win.MousePosition().Y), 0)
		txt.Draw(win, pixel.IM)

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
			// pic, err := loadPicture("../../textures/10.jpg")
			// if err != nil {
			// 	panic(err)
			// }
			// sp := pixel.NewSprite(pic, rectByWonder9(6))
			// sp.Draw(win, pixel.IM.Moved(pixel.V(500, 300)).sca)
			drawWonder(&wonder{Index: 0, Rect: pixel.R(0, 0, wonderWidth, wonderHeight)}, win)
		}

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
	selectedCardIndex int = -1
)

type TableCards struct {
	Rects []pixel.Rect
	Cards core.CardsState
}

func drawIEpoh(vi pixel.Vec) (out []pixel.Rect) {
	// if len(in) < 20 {
	// 	panic("wrong count of cards for first epoch")
	// }

	var cardRectFn = func(v pixel.Vec) pixel.Rect {
		return pixel.R(v.X, v.Y, v.X+cardWidth, v.Y+cardHeight)
	}

	vi.X -= cardWidth + deltaEpoh/2
	vi.Y -= cardHeight

	var v pixel.Vec
	for i := 0; i < 5; i++ {
		v = vi
		v.X -= float64(0+i) / 2 * (cardWidth + deltaEpoh)
		v.Y -= cardTitleHeight * float64(i)
		for j := 0; j < i+2; j++ {
			out = append(out, cardRectFn(v))
			v.X += cardWidth + deltaEpoh
		}
	}

	return
}

func drawIIIEpoh(topCenter pixel.Vec) (out []pixel.Rect) {
	var v = topCenter
	// y
	// ^
	// |
	// |
	// |
	// +---------> x
	v.Y -= cardHeight

	var width = cardWidth
	var width2 = width / 2
	var dx = deltaEpoh
	var dx2 = dx / 2
	var dy = cardTitleHeight

	var vs []pixel.Vec

	// 0 1
	vs = append(vs, pixel.V(v.X-width-dx2, v.Y))
	vs = append(vs, pixel.V(v.X+dx2, v.Y))

	// 2 3 4
	v.Y -= dy
	vs = append(vs, pixel.V(v.X-width-dx-width2, v.Y))
	vs = append(vs, pixel.V(v.X-width2, v.Y))
	vs = append(vs, pixel.V(v.X+width2+dx, v.Y))

	// 5 6 7 8
	v.Y -= dy
	vs = append(vs, pixel.V(v.X-width-dx-width-dx2, v.Y))
	vs = append(vs, pixel.V(v.X-width-dx2, v.Y))
	vs = append(vs, pixel.V(v.X+dx2, v.Y))
	vs = append(vs, pixel.V(v.X+dx2+width+dx, v.Y))

	// 9 _ 10
	v.Y -= dy
	vs = append(vs, pixel.V(v.X-width-dx-width2, v.Y))
	vs = append(vs, pixel.V(v.X+width2+dx, v.Y))

	// 11 12 13 14
	v.Y -= dy
	vs = append(vs, pixel.V(v.X-width-dx-width-dx2, v.Y))
	vs = append(vs, pixel.V(v.X-width-dx2, v.Y))
	vs = append(vs, pixel.V(v.X+dx2, v.Y))
	vs = append(vs, pixel.V(v.X+dx2+width+dx, v.Y))

	// 15 16 17
	v.Y -= dy
	vs = append(vs, pixel.V(v.X-width-dx-width2, v.Y))
	vs = append(vs, pixel.V(v.X-width2, v.Y))
	vs = append(vs, pixel.V(v.X+width2+dx, v.Y))

	// 18 19
	v.Y -= dy
	vs = append(vs, pixel.V(v.X-width-dx2, v.Y))
	vs = append(vs, pixel.V(v.X+dx2, v.Y))

	for _, v := range vs {
		out = append(out, pixel.R(v.X, v.Y, v.X+cardWidth, v.Y+cardHeight))
	}

	return
}

// func drawFirstEpohBak(win pixel.Target, vi pixel.Vec) {
// 	vi.X -= cardWidth
// 	vi.Y -= cardHeight

// 	var v pixel.Vec
// 	cardNum := 0
// 	for i := 0; i < 5; i++ {
// 		v = vi
// 		v.X -= float64(1+i) / 2 * (cardWidth + deltaEpoh)
// 		v.Y -= cardTitleHeight * float64(i)
// 		for j := 0; j < i+2; j++ {
// 			drawCard(cardNum, v, win)
// 			v.X += cardWidth + deltaEpoh
// 			cardNum++
// 		}
// 	}
// }

var (
	showCard bool

	atlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)
)

// func drawCard(i int, v pixel.Vec, win pixel.Target) {
// 	im := cardIM
// 	if !showCard {
// 		im = im.Scaled(pixel.ZV, 0.5).Moved(pixel.V(cardWidth/2, cardHeight/2))
// 	} else {
// 		im = im.Moved(pixel.V(cardWidth, cardHeight))
// 	}
// 	cards[i].Draw(win, im.Moved(v))
// }
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
	fmt.Fprintf(txt, "Index: %d", c.ID)
	txt.Draw(win, pixel.IM)
}

func drawWonder(w *wonder, win pixel.Target) {
	wondersTx[w.Index].Draw(win, pixel.IM.Scaled(pixel.ZV, 0.5).Moved(pixel.V(wonderWidth/2, wonderHeight/2)).Moved(w.Rect.Min))
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
