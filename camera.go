package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nvlled/dinojump/common"
	"github.com/nvlled/dinojump/rect"
	"github.com/nvlled/dinojump/vector"
)

type Camera struct {
	Rect      rect.T
	InnerRect rect.T
}

func NewCamera(centerX, centerY, width, height, innerSize float64) *Camera {
	r := rect.Create(0, 0, width, height)
	inner := rect.Create(0, 0, width*innerSize, height*float64(innerSize))

	r.SetMidXY(centerX, centerY)
	inner.SetMidXY(centerX, centerY)

	return &Camera{
		Rect:      r,
		InnerRect: inner,
	}
}

func (camera *Camera) CenterAt(pos *vector.T) {
	camera.Rect.SetMid(pos)
}

func (camera *Camera) Follow(pos *vector.T) {
	r := &camera.InnerRect
	if pos.X < r.Left() {
		r.SetLeft(pos.X)
	} else if pos.X > r.Right() {
		r.SetRight(pos.X)
	}
	if pos.Y < r.Top() {
		r.SetTop(pos.Y)
	} else if pos.Y > r.Bottom() {
		r.SetBottom(pos.Y)
	}

	camera.Rect.SetMidXY(camera.InnerRect.MidXY())
}

func (camera *Camera) Render(world *ebiten.Image) *ebiten.Image {
	x, y := camera.Rect.XY()
	w, h := camera.Rect.Dimension()
	return world.SubImage(common.ImageRect(x, y, w, h)).(*ebiten.Image)
}
