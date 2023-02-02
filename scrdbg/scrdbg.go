package scrdbg

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var Default = Scrdbg{
	Padding: 5,
}
var textHeight = 16 // based on ebiten's drawDebugText ch

type Scrdbg struct {
	lineNum int
	Padding int
	Screen  *ebiten.Image
}

func (scrdbg *Scrdbg) Printf(format string, args ...any) {
	if scrdbg.Screen == nil {
		return
	}
	if scrdbg.lineNum == 0 {
		Default.Screen.Clear()
	}

	lines := 1
	for _, ch := range format {
		if ch == '\n' {
			lines++
		}
	}

	y := scrdbg.lineNum * (textHeight + scrdbg.Padding)
	ebitenutil.DebugPrintAt(scrdbg.Screen, fmt.Sprintf(format, args...), 10, y)
	scrdbg.lineNum += lines
}

func (scrdbg *Scrdbg) Reset() {
	scrdbg.lineNum = 0
}

func Printf(format string, args ...any) {
	Default.Printf(format, args...)
}

func Reset() {
	Default.Reset()
}
