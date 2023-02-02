package common

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/nvlled/carrot"
)

type Void struct{}

var None = Void{}

func ImageRect(x, y, w, h float64) image.Rectangle {
	return image.Rect(int(x), int(y), int(x+w), int(y+h))
}

var colorxt = struct {
	Red   color.Color
	Blue  color.Color
	Green color.Color
}{
	Red:   color.RGBA{255, 0, 0, 128},
	Blue:  color.RGBA{0, 0, 255, 128},
	Green: color.RGBA{0, 255, 0, 128},
}

func RangeSlice(start, end int) []int {
	nums := make([]int, end-start+1)
	index := 0
	for n := start; n <= end; n++ {
		nums[index] = n
		index++
	}
	return nums
}

func AwaitKey(ctrl carrot.Control, key ebiten.Key, abortOpt ...*bool) {
	var abort *bool
	if len(abortOpt) > 0 {
		abort = abortOpt[0]
	}

	for abort == nil || !*abort {
		if inpututil.IsKeyJustPressed(key) {
			break
		}
		ctrl.Yield()
	}
}

func RuhOh(err error) {
	if err != nil {
		panic(err)
	}
}

func Throttle(fn func(), tickRate int) func() {
	n := 0
	return func() {
		if n%tickRate == 0 {
			fn()
		}
		n++
	}
}
