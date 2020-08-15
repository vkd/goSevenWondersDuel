package gfx

import (
	"fmt"
	"image"

	// Register image's formats
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/faiface/pixel"
)

var (
	cardsTx     [16*4 + 9]*pixel.Sprite
	cardsTxBack [4]*pixel.Sprite

	wondersTx     [13]*pixel.Sprite
	wondersTxBack *pixel.Sprite

	progressTx [10]*pixel.Sprite
)

const (
	pathPrefix = ""
)

func loadTextures() error { //nolint: funlen
	var nextTxToLoadIndex int
	for _, ll := range []struct {
		filename string
		count    int
	}{
		{pathPrefix + "textures/1.jpg", 16},
		{pathPrefix + "textures/3.jpg", 16},
		{pathPrefix + "textures/5.jpg", 16},
		{pathPrefix + "textures/7.jpg", 16},
		{pathPrefix + "textures/9.jpg", 9},
	} {
		pic, err := loadPicture(ll.filename)
		if err != nil {
			return fmt.Errorf("error on load texture (%s): %v", ll.filename, err)
		}
		from := nextTxToLoadIndex
		to := nextTxToLoadIndex + ll.count
		for i := from; i < to; i++ {
			cardsTx[nextTxToLoadIndex] = pixel.NewSprite(pic, rectByCard(i))
			nextTxToLoadIndex++
		}
	}

	pic, err := loadPicture(pathPrefix + "textures/4.jpg")
	if err != nil {
		return fmt.Errorf("error on load texture: %v", err)
	}
	cardsTxBack[0] = pixel.NewSprite(pic, rectByCard(0))
	cardsTxBack[1] = pixel.NewSprite(pic, rectByCard(4))

	pic, err = loadPicture(pathPrefix + "textures/10.jpg")
	if err != nil {
		return fmt.Errorf("error on load texture: %v", err)
	}
	cardsTxBack[2] = pixel.NewSprite(pic, rectByCard(2))
	cardsTxBack[3] = pixel.NewSprite(pic, rectByCard(0))
	wondersTxBack = pixel.NewSprite(pic, rectByWonder9(6))

	pic, err = loadPicture(pathPrefix + "textures/9.jpg")
	if err != nil {
		return fmt.Errorf("error on load texture: %v", err)
	}
	wondersTx[12] = pixel.NewSprite(pic, rectByWonder9(5))
	wondersTx[0] = pixel.NewSprite(pic, rectByWonder9(6))
	wondersTx[1] = pixel.NewSprite(pic, rectByWonder9(7))

	pic, err = loadPicture(pathPrefix + "textures/11.jpg")
	if err != nil {
		return fmt.Errorf("error on load texture: %v", err)
	}
	for i := 2; i < 10; i++ {
		wondersTx[i] = pixel.NewSprite(pic, rectByWonder(i-2))
	}

	pic, err = loadPicture(pathPrefix + "textures/13.jpg")
	if err != nil {
		return fmt.Errorf("error on load texture: %v", err)
	}
	for i := 10; i < 12; i++ {
		wondersTx[i] = pixel.NewSprite(pic, rectByWonder(i-10))
	}

	pic, err = loadPicture(pathPrefix + "textures/progress_tokens.png")
	if err != nil {
		return fmt.Errorf("error on load textures: %v", err)
	}
	for i := 0; i < 10; i++ {
		var x = float64(159 * i)
		var y float64 = 0
		progressTx[i] = pixel.NewSprite(pic, pixel.R(x, y, x+texturePTokenWidth, y+texturePTokenWidth))
	}

	return nil
}

const (
	textureCardWidth  float64 = 264
	textureCardHeight float64 = 400

	textureWonderWidth  float64 = 588
	texturewonderHeight float64 = 382

	texturePTokenWidth float64 = 159
)

var (
	cardLefts   = [4]float64{90, 355, 621, 887}
	cardBottoms = [4]float64{75, 476, 878, 1280}
)

func rectByCard(i int) pixel.Rect {
	i %= 16
	x := cardLefts[i%4]
	y := cardBottoms[3-(i/4)]
	return pixel.R(x, y, x+textureCardWidth, y+textureCardHeight)
}

func rectByWonder(i int) pixel.Rect {
	i %= 8
	x := wonderLefts[i%2]
	y := wonderBottoms[3-(i/2)]
	return pixel.R(x, y, x+textureWonderWidth, y+texturewonderHeight)
}

func rectByWonder9(i int) pixel.Rect {
	i %= 8
	x := wonderLefts[i%2]
	y := []float64{57, 441}[3-(i/2)]
	return pixel.R(x, y, x+textureWonderWidth, y+texturewonderHeight)
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
	wonderLefts   = []float64{30, 621}
	wonderBottoms = []float64{110, 494, 878, 1262}
)
