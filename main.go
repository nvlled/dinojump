package main

import (
	// ---------------------------------------------------------
	// Uncomment only one of the following:
	dino "github.com/nvlled/carrot-example/dino_coroutine"
	//dino "github.com/nvlled/carrot-example/dino_enums"
	//dino "github.com/nvlled/carrot-example/dino_func"
	// ---------------------------------------------------------

	"image/color"
	"log"
	"math"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/nvlled/carrot"
	"github.com/nvlled/carrot-example/action"
	"github.com/nvlled/carrot-example/common"

	"github.com/nvlled/carrot-example/level"
	"github.com/nvlled/carrot-example/rect"
	"github.com/nvlled/carrot-example/scrdbg"
	"github.com/nvlled/carrot-example/vector"

	_ "image/jpeg"
)

var renderTileSize = 50

var Debug = false

var initialized sync.Once

type Game struct {
	viewSize  vector.T
	worldSize vector.T

	viewRect  rect.T
	worldRect rect.T

	dino *dino.Sprite

	canvas *ebiten.Image

	camera *Camera

	startTime time.Time
	endTime   time.Time

	level *level.T

	renderTileSize int

	UpdateActions action.ActionSet[carrot.Void]
	DrawActions   action.ActionSet[*ebiten.Image]
}

func createLevel(renderTileSize int) *level.T {
	return level.NewLevel(level.NewOptions{
		RenderTileSize:     renderTileSize,
		AtlasFilename:      "lemcraft-tiles.png",
		BackgroundFilename: "Cielo pixelado.png",
		TileMap: map[rune]level.Tile{
			'v': level.CreateTile(28),
			'^': level.CreateTile(14),
			'*': level.CreateTile(12),
			'|': level.CreateTile(11),
		},
	}, `
|vvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv|
|                                                                         |
|       **                                                                |
|      ****                                                               |
|                   *  *                                                  |
|                ***  *                                                   |
|    **        * *      * *                                               |
|   **                                                                    |
|     **************************                                          |
| *                                  ********                             |
| **  *   *****                      *                                    |
| *                  *  **                                                |
| **  * * * *  *******    *                                               |
| *           *                                                           |
|                                                                         |
|^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^|
`)
}

func NewGame() *Game {
	level := createLevel(renderTileSize)

	levelW, levelH := level.TotalSize()
	worldW, worldH := float64(levelW), float64(levelH)

	var viewW, viewH float64 = 500, 400

	game := &Game{
		renderTileSize: renderTileSize,

		viewSize:  vector.Create(viewW, viewH),
		worldSize: vector.Create(worldW, worldH),
		worldRect: rect.Create(0, 0, worldW, worldH),
		viewRect:  rect.Create(0, 0, viewW, viewH),

		canvas: ebiten.NewImage(int(worldW), int(worldH)),

		UpdateActions: *action.NewSet[carrot.Void](),
		DrawActions:   *action.NewSet[*ebiten.Image](),

		level: level,

		camera: NewCamera(
			viewW/2, viewH/2,
			viewW, viewH,
			0.6,
		),
	}

	game.dino = dino.New(level)

	return game
}

func (g *Game) Initialize() {
	dino := g.dino
	viewW, viewH := g.viewSize.XY()

	dino.Pos = vector.Create(200, 200)
	dino.DrawSize.X = float64(g.renderTileSize / 2)
	dino.DrawSize.Y = float64(g.renderTileSize / 2)
	dino.CollisionScale = vector.Create(1.9, 1.9)

	scrdbg.Default.Screen = ebiten.NewImage(int(viewW), int(viewH))
}

func (g *Game) Update() error {
	g.startTime = time.Now()
	g.dino.Update()
	initialized.Do(g.Initialize)

	g.camera.Follow(&g.dino.Pos)
	rect.Layout.Restrict(&g.camera.Rect, &g.worldRect)

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}
	g.endTime = time.Now()

	return nil
}

func (g *Game) QueueDraw(fn func(*ebiten.Image)) {
	g.DrawActions.Add(fn)
	g.DrawActions.ClearNextApply()
}

var maxDuration float64

func (g *Game) Draw(screen *ebiten.Image) {
	g.canvas.Fill(color.RGBA{32, 82, 82, 0xff})

	g.level.Draw(g.canvas, &g.camera.Rect)
	g.dino.Draw(g.canvas)

	subCanvas := g.camera.Render(g.canvas)
	screen.DrawImage(subCanvas, &ebiten.DrawImageOptions{})
	screen.DrawImage(scrdbg.Default.Screen, &ebiten.DrawImageOptions{})

	t := float64(g.endTime.Sub(g.startTime).Milliseconds())
	if t != 0 {
		println(t)
	}
	maxDuration = math.Max(t, maxDuration)

	if Debug {
		x, y := float64(g.viewSize.X), float64(g.viewSize.Y)
		ebitenutil.DrawRect(screen, x/2, y/2, 2, 2, color.RGBA{255, 0, 0, 255})
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(g.viewSize.X), int(g.viewSize.Y)
}

func main() {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(1500, 1000)
	ebiten.SetFullscreen(false)
	ebiten.SetTPS(60)

	if os.Getenv("EXIT_ON_MODIFY") == "1" {
		go handleFileChange()
	}

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}

func handleFileChange() {
	w, err := fsnotify.NewWatcher()
	common.RuhOh(err)

	w.Add(".")

	files, err := os.ReadDir(".")
	common.RuhOh(err)
	for _, f := range files {
		if f.IsDir() && f.Name() != ".git" {
			println("watching", f.Name())
			w.Add(f.Name())
		}
	}

	for e := range w.Events {
		if filepath.Ext(e.Name) != ".go" {
			continue
		}
		if e.Op == fsnotify.Write || e.Op == fsnotify.Create {
			os.Exit(0)
		}
	}
}
