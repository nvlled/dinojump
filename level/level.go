package level

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/nvlled/dinojump/assets"
	"github.com/nvlled/dinojump/ebitenx"
	"github.com/nvlled/dinojump/rect"
	"github.com/nvlled/dinojump/sprite"
)

type f64 = float64

type Tile struct {
	TileID int
	Flags  uint16
}

type TilePos struct {
	x, y     f64
	row, col int
}

type T struct {
	RenderTileSize int

	rows int
	cols int
	data []Tile

	Atlas *sprite.T

	bgImage *ebiten.Image
}

type NewOptions struct {
	AtlasFilename      string
	BackgroundFilename string

	TileMap        map[rune]Tile
	RenderTileSize int
}

func CreateTile(id int, flagsOpt ...uint16) Tile {
	var flags uint16 = 0
	if len(flagsOpt) > 0 {
		flags = flagsOpt[0]
	}
	return Tile{
		TileID: id,
		Flags:  flags,
	}
}

func NewLevel(options NewOptions, levelData string) *T {
	img, _, err := ebitenutil.NewImageFromFileSystem(assets.FS, options.AtlasFilename)
	if err != nil {
		panic(err)
	}

	sprite := sprite.New(img, 7, 8)

	levelData = strings.TrimSpace(levelData)
	lines := strings.Split(levelData, "\n")
	rows := len(lines)
	cols := 0

	for _, line := range lines {
		if cols < len(line) {
			cols = len(line)
		}
	}
	data := make([]Tile, rows*cols)

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			index := r*cols + c
			data[index] = Tile{TileID: -1}
		}
	}

	for r, line := range lines {
		for c, ch := range line {
			tile, ok := options.TileMap[ch]
			if ok {
				data[r*cols+c] = tile
			}
		}
	}

	level := &T{
		RenderTileSize: options.RenderTileSize,
		rows:           rows,
		cols:           cols,
		Atlas:          sprite,
		data:           data,

		bgImage: ebitenx.NewImageFromAssets(options.BackgroundFilename),
	}

	return level
}

/*
TODO:
left, right, top, bottom := level.GetTileIntersections(rect)
if left != nil {
	blah = *left
}
*/

func (level *T) GetTileIntersections(rect *rect.T) (
	hitFlag byte, left f64, right f64, top f64, bottom f64,
) {
	size := level.RenderTileSize
	mid := rect.Mid()

	mc := int(mid.X) / size
	mr := int(mid.Y) / size
	sl := int(rect.Left()+1) / size
	sr := int(rect.Right()-1) / size
	st := int(rect.Top()+1) / size
	sb := int(rect.Bottom()-1) / size

	if level.HasTileAt(sl, mr) {
		hitFlag |= 0b1000
	}
	if level.HasTileAt(sr, mr) {
		hitFlag |= 0b0100
	}
	if level.HasTileAt(mc, st) {
		hitFlag |= 0b0010
	}
	if level.HasTileAt(mc, sb) {
		hitFlag |= 0b0001
	}

	left = f64((sl + 1) * size)
	right = f64(sr * size)
	top = f64((st + 1) * size)
	bottom = f64(sb * size)

	return
}

func (level *T) GetTileOn(x, y f64, rect *rect.T) bool {
	size := level.RenderTileSize
	c := int(x) / size
	r := int(y) / size
	_, ok := level.GetTileAt(c, r)
	if !ok {
		return false
	}

	if rect != nil {
		rect.SetTopLeftXY(f64(c*level.RenderTileSize), f64(r*level.RenderTileSize))
	}

	return true
}

func (level *T) HasTileOn(x, y f64) bool {
	size := level.RenderTileSize
	c := int(x) / size
	r := int(y) / size
	tile, ok := level.GetTileAt(c, r)
	if !ok {
		return false
	}

	return tile.TileID >= 0
}

func (level *T) HasTileAt(c, r int) bool {
	i := r*level.cols + c
	if i < 0 || i >= len(level.data) {
		return false
	}
	return level.data[i].TileID >= 0
}

func (level *T) GetTileAt(c, r int) (Tile, bool) {
	i := r*level.cols + c
	if i < 0 || i >= len(level.data) {
		return Tile{}, false
	}
	return level.data[i], true
}

func (level *T) GetTileRectAt(c, r int) rect.T {
	s := float64(level.RenderTileSize)
	x := f64(c) * s
	y := f64(r) * s
	return rect.Create(x, y, s, s)
}

func (level *T) Size() (cols, rows int) {
	return level.cols, level.rows
}

func (level *T) TotalSize() (width, height int) {
	size := level.RenderTileSize
	return level.cols * size, level.rows * size
}

func (level *T) GetRect() rect.T {
	size := level.RenderTileSize
	return rect.CreateInt(0, 0, level.cols*size, level.rows*size)
}

func (level *T) drawBackground(canvas *ebiten.Image) {
	var nilImage *ebiten.Image
	if level.bgImage == nilImage {
		return
	}
	w, h := level.TotalSize()
	rect := rect.CreateInt(0, 0, w, h)
	ebitenx.DrawImageAtRect(canvas, level.bgImage, &rect)
}

func (level *T) Draw(canvas *ebiten.Image, view *rect.T) {
	sprite := level.Atlas
	tileSize := level.RenderTileSize

	level.drawBackground(canvas)

	ac, ar := view.Min.XY_int()
	bc, br := view.Max.XY_int()

	size := level.RenderTileSize
	ac, ar = ac/size, ar/size
	bc, br = bc/size, br/size

	destRect := rect.CreateInt(0, 0, tileSize, tileSize)
	for r := 0; r < level.rows; r++ {
		for c := 0; c < level.cols; c++ {
			if c < ac || c > bc || r < ar || r > br {
				continue
			}

			index := r*level.cols + c
			tile := level.data[index]
			if tile.TileID < 0 {
				continue
			}
			tileImg := sprite.GetTileImage(tile.TileID)
			destRect.SetTopLeftXY(f64(c*tileSize), float64(r*tileSize))
			ebitenx.DrawImageAtRect(canvas, tileImg, &destRect)
		}
	}
}
