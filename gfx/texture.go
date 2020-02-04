package gfx

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/faiface/pixel"
)

var (
	cardsTx     [16*4 + 9]*pixel.Sprite
	cardsTxBack [4]*pixel.Sprite

	wondersTx     [12]*pixel.Sprite
	wondersTxBack *pixel.Sprite

	progressTx [10]*pixel.Sprite
)

func loadTextures() error {
	var nextTxToLoadIndex int
	for _, ll := range []struct {
		filename string
		count    int
	}{
		{"../../textures/1.jpg", 16},
		{"../../textures/3.jpg", 16},
		{"../../textures/5.jpg", 16},
		{"../../textures/7.jpg", 16},
		{"../../textures/9.jpg", 9},
	} {
		pic, err := loadPicture(ll.filename)
		if err != nil {
			return fmt.Errorf("Error on load texture (%s): %v", ll.filename, err)
		}
		from := nextTxToLoadIndex
		to := nextTxToLoadIndex + ll.count
		for i := from; i < to; i++ {
			cardsTx[nextTxToLoadIndex] = pixel.NewSprite(pic, rectByCard(i))
			nextTxToLoadIndex++
		}
	}

	pic, err := loadPicture("../../textures/4.jpg")
	if err != nil {
		return fmt.Errorf("Error on load texture: %v", err)
	}
	cardsTxBack[0] = pixel.NewSprite(pic, rectByCard(0))
	cardsTxBack[1] = pixel.NewSprite(pic, rectByCard(4))

	pic, err = loadPicture("../../textures/10.jpg")
	if err != nil {
		return fmt.Errorf("Error on load texture: %v", err)
	}
	cardsTxBack[2] = pixel.NewSprite(pic, rectByCard(2))
	cardsTxBack[3] = pixel.NewSprite(pic, rectByCard(0))
	wondersTxBack = pixel.NewSprite(pic, rectByWonder9(6))

	pic, err = loadPicture("../../textures/9.jpg")
	if err != nil {
		return fmt.Errorf("Error on load texture: %v", err)
	}
	// wondersTx[13] = pixel.NewSprite(pic, rectByWonder9(5))
	wondersTx[0] = pixel.NewSprite(pic, rectByWonder9(6))
	wondersTx[1] = pixel.NewSprite(pic, rectByWonder9(7))

	pic, err = loadPicture("../../textures/11.jpg")
	if err != nil {
		return fmt.Errorf("Erroron load texture: %v", err)
	}
	for i := 2; i < 10; i++ {
		wondersTx[i] = pixel.NewSprite(pic, rectByWonder(i))
	}

	pic, err = loadPicture("../../textures/13.jpg")
	if err != nil {
		return fmt.Errorf("Error on load texture: %v", err)
	}
	for i := 10; i < 12; i++ {
		wondersTx[i] = pixel.NewSprite(pic, rectByWonder(i))
	}

	pic, err = loadPicture("../../textures/progress_tokens.png")
	if err != nil {
		return fmt.Errorf("Error on load textures: %v", err)
	}
	for i := 0; i < 10; i++ {
		y := float64(159 * i)
		progressTx[i] = pixel.NewSprite(pic, pixel.R(0, y, 159, y+159))
	}

	return nil
}

func rectByCard(i int) pixel.Rect {
	i = i % 16
	x := cardLefts[i%4]
	y := cardBottoms[3-(i/4)]
	return pixel.R(x, y, x+cardWidth*2, y+cardHeight*2)
}

func rectByWonder(i int) pixel.Rect {
	i = i % 8
	x := wonderLefts[i%2]
	y := wonderBottoms[3-(i/2)]
	return pixel.R(x, y, x+wonderWidth*2, y+wonderHeight*2)
}

func rectByWonder9(i int) pixel.Rect {
	i = i % 8
	x := wonderLefts[i%2]
	y := []float64{57, 441}[3-(i/2)]
	return pixel.R(x, y, x+wonderWidth*2, y+wonderHeight*2)
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

var (
	cardLefts   = []float64{90, 355, 621, 887}
	cardBottoms = []float64{75, 476, 878, 1280}

	wonderLefts   = []float64{30, 621}
	wonderBottoms = []float64{110, 494, 878, 1262}
)
