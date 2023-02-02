package ebitenx

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/nvlled/dinojump/assets"
	"github.com/nvlled/dinojump/rect"
	"github.com/nvlled/dinojump/vector"
)

func DrawRectT(canvas *ebiten.Image, r rect.T, color color.Color) {
	ebitenutil.DrawRect(canvas, r.X(), r.Y(), r.Width(), r.Height(), color)
}

func DrawRect(canvas *ebiten.Image, x, y, w, h float64, color color.Color) {
	ebitenutil.DrawRect(canvas, x, y, w, h, color)
}

func DrawPoint(canvas *ebiten.Image, v vector.T, c color.Color) {
	ebitenutil.DrawCircle(canvas, v.X, v.Y, 2, c)
}

func DrawPointXY(canvas *ebiten.Image, x, y float64, c color.Color) {
	ebitenutil.DrawCircle(canvas, x, y, 2, c)
}

func DrawCircle(canvas *ebiten.Image, x, y float64, c color.Color) {
	ebitenutil.DrawCircle(canvas, x, y, 2, c)
}
func DrawImageAtRect(canvas, img *ebiten.Image, r *rect.T) {
	op := ebiten.GeoM{}
	TransformImageRect(img, r, &op)
	canvas.DrawImage(img, &ebiten.DrawImageOptions{
		GeoM: op,
	})
}

func TransformImageFlip(img *ebiten.Image, op *ebiten.GeoM, flags byte) {
	if flags == 0 {
		return
	}
	sx, sy := 1.0, 1.0
	tx, ty := 1.0, 1.0

	b := img.Bounds().Size()
	if flags&0b10 != 0 {
		sx = -1
		tx = float64(b.X)
	}
	if flags&0b01 != 0 {
		sy = -1
		ty = float64(b.Y)
	}

	op.Scale(sx, sy)
	op.Translate(tx, ty)
}

func TransformImageRotate(img *ebiten.Image, op *ebiten.GeoM, theta float64) {
	bounds := img.Bounds().Size()
	op.Translate(-float64(bounds.X)/2, -float64(bounds.Y)/2)
	op.Rotate(theta)
	op.Translate(float64(bounds.X)/2, float64(bounds.Y)/2)
}

func TransformImageRect(img *ebiten.Image, rect *rect.T, op *ebiten.GeoM) {
	bounds := img.Bounds().Size()
	scaleX := rect.Width() / float64(bounds.X)
	scaleY := rect.Height() / float64(bounds.Y)

	op.Scale(scaleX, scaleY)
	op.Translate(rect.X(), rect.Y())
}

func NewImageFromAssets(filename string) *ebiten.Image {
	if filename == "" {
		return nil
	}

	img, _, err := ebitenutil.NewImageFromFileSystem(assets.FS, filename)
	if err != nil {
		panic(err)
	}

	return img
}
