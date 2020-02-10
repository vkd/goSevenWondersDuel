package gfx

import (
	_ "image/jpeg"
	"reflect"
	"testing"

	"github.com/faiface/pixel"
)

func TestRectByCard(t *testing.T) {
	tests := []struct {
		i    int
		want pixel.Vec
	}{
		// TODO: Add test cases.
		{0, pixel.V(90, 1280)},
		{3, pixel.V(887, 1280)},
		{6, pixel.V(621, 878)},
		{12, pixel.V(90, 75)},
		{13, pixel.V(355, 75)},
		{14, pixel.V(621, 75)},
		{15, pixel.V(887, 75)},
	}
	for _, tt := range tests {
		want := pixel.R(tt.want.X, tt.want.Y, tt.want.X+cardWidth*2, tt.want.Y+cardHeight*2)
		if got := rectByCard(tt.i); !reflect.DeepEqual(got, want) {
			t.Errorf("%d. rectByCard() = %v, want %v", tt.i, got, want)
		}
	}
}
