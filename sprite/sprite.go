package sprite

import (
	"image"
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/nvlled/dinojump/action"
	"github.com/nvlled/dinojump/common"
	"github.com/nvlled/dinojump/ebitenx"
	"github.com/nvlled/dinojump/rect"
	"github.com/nvlled/dinojump/vector"
)

var Debug = false

var emptyTile = image.Rectangle{}

type f64 = float64

type T struct {
	Image   *ebiten.Image
	imageOp ebiten.GeoM

	CollisionScale vector.T
	DrawSize       vector.T
	tileSize       vector.T

	Pos vector.T
	Vel vector.T

	//Animation       string
	//AnimationFrames map[string][]image.Rectangle

	//TicksPerFrame int
	CurrentTileID int

	Flip     byte
	Rotation float64

	Rect rect.T
	//topLeft     vector.T
	//bottomRight vector.T

	cols int
	rows int

	tileIDRect image.Rectangle

	Actions *action.ActionSet[common.Void]
}

func New(atlas *ebiten.Image, numCols, numRows int /*, animationIndices map[string][]int*/) *T {
	imgW, imgH := atlas.Size()
	var tileW = imgW / numCols
	var tileH = imgH / numRows

	/*
		frames := map[string][]image.Rectangle{}
			firstAnimation := ""
				for animation, indices := range animationIndices {
					if firstAnimation == "" {
						firstAnimation = animation
					}

					var rects []image.Rectangle
					for _, i := range indices {
						y := (i / numCols) * tileH
						x := (i % numCols) * tileW
						rects = append(rects, image.Rect(x, y, x+tileW, y+tileH))
					}
					frames[animation] = rects
				}

				if _, ok := frames["default"]; !ok {
					frames["default"] = frames[firstAnimation]
				}
	*/

	size := vector.T{float64(tileW), float64(tileH)}

	return &T{
		Image:    atlas,
		DrawSize: size,
		tileSize: size,

		CollisionScale: vector.Unit,

		Pos: vector.Zero,
		Vel: vector.Zero,

		//Animation:       "default",
		//AnimationFrames: frames,
		//TicksPerFrame: 7,

		cols: numCols,
		rows: numCols,

		Actions: action.NewSet[common.Void](),
	}
}

func (sprite *T) GetTileCount() (cols, rows int) {
	return sprite.cols, sprite.rows
}

func (sprite *T) GetTileImage(index int) *ebiten.Image {
	return sprite.Image.SubImage(sprite.GetTile(index)).(*ebiten.Image)
}

func (sprite *T) GetTile(index int) image.Rectangle {
	if index < 0 || index >= sprite.cols*sprite.rows {
		return emptyTile
	}

	tileW, tileH := sprite.tileSize.XY_int()
	cols, _ := sprite.GetTileCount()
	y := (index / cols) * tileH
	x := (index % cols) * tileW
	return image.Rect(x, y, x+tileW, y+tileH)
}

func (sprite *T) CurrentTile() image.Rectangle {
	return sprite.GetTile(sprite.CurrentTileID)

	/*
		i := sprite.FrameIndex
		if frames, ok := sprite.AnimationFrames[sprite.Animation]; ok {
			return frames[i]
		}
		return sprite.defaultFrame
	*/
}

/*
func (sprite *Sprite) TickFrame() {
	sprite.ticks++
	if sprite.ticks >= sprite.TicksPerFrame {
		sprite.ticks = 0
		sprite.FrameIndex++

		if frames, ok := sprite.AnimationFrames[sprite.Animation]; ok {
			if sprite.FrameIndex >= len(frames) {
				sprite.FrameIndex = 0
			}
		}
	}
}
*/

func (sprite *T) Update() {
	sprite.Rect.Min = vector.AddXY(&sprite.Pos, -sprite.DrawSize.X/2, -sprite.DrawSize.Y/2)
	sprite.Rect.Max = vector.AddXY(&sprite.Pos, sprite.DrawSize.X/2, sprite.DrawSize.Y/2)

	sprite.Actions.Apply(common.None)
}

func (sprite *T) Draw(canvas *ebiten.Image) {
	subImg := sprite.Image.SubImage(sprite.CurrentTile()).(*ebiten.Image)

	vr := sprite.GetViewRect()
	cr := sprite.GetCollisionRect()
	mid := vr.Mid()

	//sprite.imageOp.Reset()
	//sprite.imageOp.Skew(0.01, 0)

	op := sprite.imageOp

	//op.Skew(-0.5, 0)
	//op.Translate(0.5*sprite.tileSize.X/2, 0)

	ebitenx.TransformImageFlip(subImg, &op, sprite.Flip)
	ebitenx.TransformImageRotate(subImg, &op, sprite.Rotation)
	ebitenx.TransformImageRect(subImg, &vr, &op)

	canvas.DrawImage(subImg, &ebiten.DrawImageOptions{GeoM: op})

	if Debug {
		ebitenx.DrawRectT(canvas, vr, color.RGBA{50, 50, 50, 30})
		ebitenutil.DrawRect(canvas, mid.X, mid.Y, 2, 2, color.RGBA{0, 255, 255, 255})

		ebitenx.DrawRectT(canvas, sprite.Rect, color.RGBA{150, 50, 50, 90})
		ebitenx.DrawRectT(canvas, sprite.GetCollisionRect(), color.RGBA{50, 150, 50, 90})

		_ = cr

		rect.Layout.Center(&cr, &vr)
		//ebitenx.DrawRect(canvas, cr, colorxt.Green)
		//ebitenutil.DrawRect(canvas, 50, 50, 20, 30, colorxt.Blue)
		//ebitenutilx.DrawPoint(canvas, cr.Mid(), colorxt.Green)
		{
			r := rect.FromImageRect(sprite.tileIDRect)

			rect.Layout.Center(&r, &vr)
			rect.Layout.Top(&r, &vr)
			x, y := r.Min.XY_int()
			ebitenutil.DebugPrintAt(canvas, strconv.Itoa(sprite.CurrentTileID), x, y)
		}
	}

}

func (sprite *T) GetViewRect() rect.T {
	r := sprite.Rect
	r.Scale(&sprite.CollisionScale)
	return r
}

func (sprite *T) GetCollisionRect() rect.T {
	return sprite.Rect
}

func (sprite *T) SetLeft(x f64) {
	sprite.Rect.SetLeft(x)
	sprite.Pos.Set(sprite.Rect.MidXY())
}

func (sprite *T) SetRight(x f64) {
	sprite.Rect.SetRight(x)
	sprite.Pos.Set(sprite.Rect.MidXY())
}
func (sprite *T) SetTop(y f64) {
	sprite.Rect.SetTop(y)
	sprite.Pos.Set(sprite.Rect.MidXY())
}
func (sprite *T) SetBottom(y f64) {
	sprite.Rect.SetBottom(y)
	sprite.Pos.Set(sprite.Rect.MidXY())
}

func (sprite *T) RestrictWithin(container *rect.T) (bool, bool) {
	changedX, changedY := rect.Layout.Restrict(&sprite.Rect, container)
	if changedX || changedY {
		sprite.Pos.Set(sprite.Rect.MidXY())
	}
	return changedX, changedY
}

func (sprite *T) DrawDebugImage(canvas *ebiten.Image, baseX, baseY float64) {
	tw, th := sprite.DrawSize.XY()
	w, h := float64(sprite.cols)*tw, float64(sprite.rows)*th
	baseX -= w / 2
	baseY -= h

	destRect := rect.CreateInt(0, 0, int(tw), int(th))
	for y := 0; y < sprite.rows; y++ {
		for x := 0; x < sprite.cols; x++ {
			i := y*sprite.cols + x
			tile := sprite.GetTileImage(i)
			destRect.SetTopLeftXY(baseX+float64(x)*tw, baseY+float64(y)*th)
			ebitenx.DrawImageAtRect(canvas, tile, &destRect)
			ebitenutil.DebugPrintAt(
				canvas, strconv.Itoa(i),
				int(destRect.MidX()), int(destRect.MidY()),
			)
		}
	}
}
