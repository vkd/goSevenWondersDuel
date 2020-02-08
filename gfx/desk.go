package gfx

import "github.com/faiface/pixel"

type deskRowGrid struct {
	itemWidth, xMargin float64
}

func (d deskRowGrid) genRowCenter(center float64, count int) (out []float64) {
	if count <= 0 {
		return
	}
	var totalWidth = d.itemWidth*float64(count) + d.xMargin*float64(count-1)
	var x = center - (totalWidth / 2)
	for i := 0; i < count; i++ {
		out = append(out, x)
		x += d.itemWidth + d.xMargin
	}
	return
}

var (
	ageGrid = deskGrid{
		deskRowGrid: deskRowGrid{
			itemWidth: cardWidth,
			xMargin:   deltaEpoh,
		},
		dy: -cardTitleHeight,
	}
)

type deskGrid struct {
	deskRowGrid
	dy float64
}

func (d deskGrid) genAgeIVecs(v pixel.Vec) (vs []pixel.Vec) {
	var y = v.Y
	//      0  1
	//     2  3  4
	//    5  6  7  8
	//   9 10 11 12 13
	// 14 15 16 17 18 19
	for i := 2; i <= 6; i++ {
		for _, x := range d.genRowCenter(v.X, i) {
			vs = append(vs, pixel.V(x, y))
		}
		y += d.dy
	}
	return
}

func (d deskGrid) genAgeIIVecs(v pixel.Vec) (vs []pixel.Vec) {
	var y = v.Y
	// 0  1  2  3  4  5
	//  6  7  8  9  10
	//   11 12 13 14
	//    15 16 17
	//     18 19
	for i := 6; i >= 2; i-- {
		for _, x := range d.genRowCenter(v.X, i) {
			vs = append(vs, pixel.V(x, y))
		}
		y += d.dy
	}
	return
}

func (d deskGrid) genAgeIIIVecs(v pixel.Vec) (vs []pixel.Vec) {
	var y = v.Y

	//   0 1
	//  2 3 4
	// 5 6 7 8
	for i := 2; i <= 4; i++ {
		for _, x := range ageGrid.genRowCenter(v.X, i) {
			vs = append(vs, pixel.V(x, y))
		}
		y += d.dy
	}

	// 9 _ 10
	row := ageGrid.genRowCenter(v.X, 3)
	vs = append(vs, pixel.V(row[0], y))
	vs = append(vs, pixel.V(row[2], y))
	y += d.dy

	// 11  12  13  14
	//   15  16  17
	//     18  19
	for i := 4; i >= 2; i-- {
		for _, x := range ageGrid.genRowCenter(v.X, i) {
			vs = append(vs, pixel.V(x, y))
		}
		y += d.dy
	}
	return
}

func genCardRects(vs []pixel.Vec) (out []pixel.Rect) {
	for _, v := range vs {
		out = append(out, pixel.R(v.X, v.Y, v.X+cardWidth, v.Y+cardHeight))
	}

	return
}
