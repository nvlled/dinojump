package dino

import (
	"bytes"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/nvlled/carrot"
	"github.com/nvlled/carrot-example/assets"
	"github.com/nvlled/carrot-example/bitf"
	"github.com/nvlled/carrot-example/common"
	"github.com/nvlled/carrot-example/level"
	"github.com/nvlled/carrot-example/numsign"
	"github.com/nvlled/carrot-example/seqiter"
	"github.com/nvlled/carrot-example/sprite"
)

var dinoMaxSpeed = 20

type Sprite struct {
	sprite.T

	Hit bitf.T

	Level *level.T

	animationScript  *carrot.Script
	controllerScript *carrot.Script
	animations       seqiter.Iterator[carrot.Coroutine]
}

func New(level *level.T) *Sprite {
	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(assets.DinoSpriteData))
	if err != nil {
		panic(err)
	}

	dino := &Sprite{
		T:                *sprite.New(img, 24, 1),
		Level:            level,
		animationScript:  carrot.Create(),
		controllerScript: carrot.Create(),
	}
	dino.T.CollisionScale.X = 0.7
	dino.T.CollisionScale.Y = 0.7

	dino.animations = seqiter.CreateSeqIterator(
		dino.AnimateIdle,
		dino.AnimateWalk,
		dino.AnimateRun,
		dino.AnimatePreview,
	)

	dino.controllerScript.Transition(dino.ControllerCoroutine)

	return dino
}

func (dino *Sprite) Update() {
	dino.T.Update()
	dino.controllerScript.Update()
	dino.animationScript.Update()

	if inpututil.IsKeyJustPressed(ebiten.KeyF8) {
		dino.animationScript.Cancel()
		dino.controllerScript.Transition(dino.ControlTestFrame)
		sprite.Debug = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF9) {
		dino.animationScript.Transition(dino.AnimateIdle)
		dino.controllerScript.Transition(dino.ControllerCoroutine)
	}
}

func (dino *Sprite) SetAnimation(coroutine carrot.Coroutine) {
	dino.animationScript.Transition(coroutine)
}

func (dino *Sprite) SetController(coroutine carrot.Coroutine) {
	dino.controllerScript.Transition(coroutine)
}

func (dino *Sprite) AnimateOuchie(ctrl *carrot.Control) {
	frames := seqiter.CreateSeqIterator(14, 16)
	for {
		dino.CurrentTileID = frames.Next()
		ctrl.Delay(3)
	}
}

func (dino *Sprite) AnimateIdle(ctrl *carrot.Control) {
	frames := seqiter.CreateSeqIterator(0, 1, 2, 3)
	for {
		dino.CurrentTileID = frames.Next()
		ctrl.Delay(7)
	}
}

func (dino *Sprite) AnimateWalk(ctrl *carrot.Control) {
	frames := seqiter.CreateSeqIterator(3, 4, 5, 6, 7, 8)
	for {
		dino.CurrentTileID = frames.Next()
		ctrl.Delay(10)
	}
}

func (dino *Sprite) AnimateRun(ctrl *carrot.Control) {
	frames := seqiter.CreateSeqIterator(common.RangeSlice(18, 23)...)
	for {
		dino.CurrentTileID = frames.Next()
		if dino.Vel.X >= float64(dinoMaxSpeed) {
			ctrl.Delay(1)
		} else {
			ctrl.Delay(3)
		}
	}
}

func (dino *Sprite) AnimateFly(ctrl *carrot.Control) {
	frames := seqiter.CreateSeqIterator(17, 18)
	for {
		dino.CurrentTileID = frames.Next()
		ctrl.Delay(10)
	}
}

func (dino *Sprite) AnimatePreview(ctrl *carrot.Control) {
	for {
		ctrl.Yield()
	}
}

func (dino *Sprite) ControlTestPosition(ctrl *carrot.Control) {
	rect := dino.Level.GetTileRectAt(3, 3)
	dino.SetTop(rect.Top())
	dino.SetLeft(rect.Left())

	for {

		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
			dino.SetTop(rect.Top())
			dino.SetLeft(rect.Right())
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
			dino.SetTop(rect.Top())
			dino.SetRight(rect.Left())
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			dino.SetTop(rect.Bottom())
			dino.SetLeft(rect.Left())
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			dino.SetBottom(rect.Top())
			dino.SetLeft(rect.Left())
		}

		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			dino.controllerScript.Transition(dino.ControllerCoroutine)
		}

		ctrl.Yield()
	}
}

func (dino *Sprite) ControlTestFrame(ctrl *carrot.Control) {
	dino.Actions.ClearNextApply()

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		anim := dino.animations.Next()
		dino.SetAnimation(anim)
	}
	ids := common.RangeSlice(0, 23)
	frames := seqiter.CreateSeqIterator(ids...)
	for {
		ctrl.Delay(1)
		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
			dino.CurrentTileID = frames.Prev()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
			dino.CurrentTileID = frames.Next()
		}
	}
}

func (dino *Sprite) RestrictInWorld(common.Void) {
	r := dino.Level.GetRect()
	cx, _ := dino.RestrictWithin(&r)
	if cx {
		dino.Vel.X = 0
	}
}

func (dino *Sprite) CollideWithTile(common.Void) {
	level := dino.Level
	colRect := dino.GetCollisionRect()
	byte, l, r, t, b := level.GetTileIntersections(&colRect)
	hit := bitf.T(byte)

	dino.Hit = hit

	if hit.Some(0b1000) {
		dino.SetLeft(l - 1)
	} else if hit.Some(0b0100) {
		dino.SetRight(r)
	}

	if hit.Some(0b0010) {
		dino.SetTop(t - 1)
	} else if hit.Some(0b0001) {
		dino.SetBottom(b + 1)
	}
}

func (dino *Sprite) ApplyGravity(common.Void) {
	if dino.Hit.Some(0b0001) {
		dino.Vel.Y = 0
	} else {
		dino.Vel.Y += 0.25

	}
	dino.Pos.Y += dino.Vel.Y
}

func (dino *Sprite) ControllerCoroutine(ctrl *carrot.Control) {
	turns := 0
	jumps := 0
	maxJumps := 3
	jumpCharge := 0

	dino.Actions.Add(dino.ApplyGravity)
	dino.Actions.Add(dino.CollideWithTile)

IDLE:
	{ // ---------------------------------------------------------
		println("idle")
		dino.SetAnimation(dino.AnimateIdle)
		for {
			walk := false
			if ebiten.IsKeyPressed(ebiten.KeyLeft) {
				dino.Flip = 0b10
				numsign.Set(&dino.Vel.X, -1)
				walk = true
			} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
				dino.Flip = 0b00
				numsign.Set(&dino.Vel.X, 1)
				walk = true
			}

			if !dino.Hit.Some(0b0001) {
				goto FALL
			}

			if walk {
				goto WALK
			}
			if ebiten.IsKeyPressed(ebiten.KeySpace) {
				goto JUMP
			}

			ctrl.Yield()
		}
	} // ---------------------------------------------------------

WALK:
	{ // ---------------------------------------------------------
		println("walk")
		dino.Vel.X = 0.5
		dino.SetAnimation(dino.AnimateWalk)
		for {
			oldSign := numsign.Get(dino.Vel.X)
			if ebiten.IsKeyPressed(ebiten.KeyLeft) {
				dino.Flip = 0b10
				numsign.Set(&dino.Vel.X, -1)
			} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
				dino.Flip = 0b00
				numsign.Set(&dino.Vel.X, 1)
			} else {
				goto IDLE
			}

			if !dino.Hit.Some(0b0001) {
				goto FALL
			}

			if oldSign != numsign.Get(dino.Vel.X) {
				dino.Vel.X *= 0.5
			}

			dino.Pos.X += dino.Vel.X
			dino.Vel.X += 0.1 * numsign.Get(dino.Vel.X)

			if ebiten.IsKeyPressed(ebiten.KeySpace) {
				goto JUMP
			}
			if math.Abs(dino.Vel.X) > 4.5 {
				goto RUN
			}

			ctrl.Yield()
		}
	} // ---------------------------------------------------------

BRAKE:
	{
		println("brake")
		dino.SetAnimation(dino.AnimateWalk)
		for {
			if ebiten.IsKeyPressed(ebiten.KeyLeft) {
				dino.Flip = 0b10
			} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
				dino.Flip = 0b00
			}
			if ebiten.IsKeyPressed(ebiten.KeySpace) {
				goto JUMP
			}
			if dino.Hit.Some(0b1100) && math.Abs(dino.Vel.X) >= float64(dinoMaxSpeed)*0.55 {
				goto BOUNCE
			}
			if math.Abs(dino.Vel.X) <= 1 {
				goto IDLE
			}
			dino.Pos.X += dino.Vel.X

			dirX := numsign.Get(dino.Vel.X)
			if (ebiten.IsKeyPressed(ebiten.KeyLeft) && dirX > 0) ||
				(ebiten.IsKeyPressed(ebiten.KeyRight) && dirX < 0) {
				dino.Vel.X *= 0.90
			} else {
				dino.Vel.X *= 0.97

			}

			ctrl.Yield()
		}
	} // ---------------------------------------------------------

BOUNCE:
	{ // ---------------------------------------------------------
		println("bounce")
		dino.SetAnimation(dino.AnimateOuchie)
		dino.Vel.X *= -0.8
		dino.Vel.Y = -4.5

		for {
			dino.Pos.X += dino.Vel.X
			dino.Pos.Y += dino.Vel.Y
			dino.Vel.X *= 0.9

			if dino.Vel.Y < 0 {
				dino.Vel.Y += 0.25
			}

			if math.Abs(dino.Vel.X) < 1 {
				dino.Vel.Y = 0
				dino.Vel.X = 0
				goto IDLE
			}

			ctrl.Yield()
		}
	} // ---------------------------------------------------------

RUN:
	{ // ---------------------------------------------------------
		println("run")
		dino.SetAnimation(dino.AnimateRun)
		for {
			dirX := numsign.Get(dino.Vel.X)
			leftDown := ebiten.IsKeyPressed(ebiten.KeyLeft)
			rightDown := ebiten.IsKeyPressed(ebiten.KeyRight)
			noDown := !leftDown && !rightDown
			brake := noDown || (leftDown && dirX == 1) || (rightDown && dirX == -1)

			if ebiten.IsKeyPressed(ebiten.KeySpace) {
				goto JUMP
			}

			if !dino.Hit.Some(0b0001) {
				goto FALL
			}

			if brake {
				goto BRAKE
			}

			if leftDown {
				dino.Flip = 0b10
				numsign.Set(&dino.Vel.X, -1)
			} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
				dino.Flip = 0b00
				numsign.Set(&dino.Vel.X, 1)
			}

			if dino.Hit.Some(0b1100) && math.Abs(dino.Vel.X) >= float64(dinoMaxSpeed)*0.55 {
				goto BOUNCE
			}

			dino.Pos.X += dino.Vel.X

			if dino.Vel.X < float64(dinoMaxSpeed) && !dino.Hit.Some(0b1100) {
				dino.Vel.X += 0.2 * numsign.Get(dino.Vel.X)
			}

			ctrl.Yield()
		}
	} // ---------------------------------------------------------

JUMP:
	{ // ---------------------------------------------------------
		println("jump")
		jumps++
		dino.Actions.Remove(dino.ApplyGravity)

		dino.CurrentTileID = 11
		ctrl.Yield()
		dino.CurrentTileID = 12

		dino.animationScript.Cancel()
		ctrl.YieldUntil(dino.animationScript.IsDone)

		dino.Vel.Y = -7.5
		jumpCharge = 0

		for {
			leftDown := ebiten.IsKeyPressed(ebiten.KeyLeft)
			rightDown := ebiten.IsKeyPressed(ebiten.KeyRight)

			if math.Abs(dino.Vel.X) < float64(dinoMaxSpeed) {
				if leftDown {
					dino.Flip = 0b10
					numsign.Set(&dino.Vel.X, -1)
					dino.Pos.X -= 1.0
				} else if rightDown {
					dino.Flip = 0b00
					numsign.Set(&dino.Vel.X, 1)
					dino.Pos.X += 1.0
				} else {
					numsign.Set(&dino.Vel.X, 0)
				}
			}

			dirX := numsign.Get(dino.Vel.X)
			swerve := (leftDown && dirX == 1) || (rightDown && dirX == -1)
			if swerve {
				dino.Vel.X *= 0.7
			}

			if !dino.Hit.Some(0b1100) {
				dino.Pos.X += dino.Vel.X
			}
			if ebiten.IsKeyPressed(ebiten.KeySpace) {
				dino.Vel.Y *= 0.95
			} else {
				dino.Vel.Y *= 0.55
			}
			dino.Pos.Y += dino.Vel.Y

			jumpCharge++
			if jumpCharge >= 40 && ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
				goto JUMP_CHARGE
			}

			if dino.Hit.Some(0b0010) {
				dino.Vel.Y *= -0.3
				goto FALL
			}

			if math.Abs(dino.Vel.Y) < 1 {
				goto FALL
			}

			if jumps >= maxJumps && math.Abs(dino.Vel.X) >= float64(dinoMaxSpeed) {
				goto FLY
			} else if jumps >= 2 {
				dino.Rotation += dirX * float64(jumps) / 10
			}

			if dino.Hit.Some(0b1100) && math.Abs(dino.Vel.X) >= float64(dinoMaxSpeed)*0.55 {
				goto BOUNCE
			}

			turns++
			ctrl.Yield()
		}
	} // ---------------------------------------------------------

JUMP_CHARGE:
	{ // ---------------------------------------------------------
		println("jump charge")
		dino.Actions.Remove(dino.ApplyGravity)
		dino.Actions.Remove(dino.CollideWithTile)

		ctrl.Yield()
		n := 0.5
		size := dino.DrawSize
		pos := dino.Pos
		decreaseStep := float64(1)

		idle := 0
		pressed := float64(0)

		for pressed < 10 {
			if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
				idle = 0
				pressed++
				n += 0.05
			} else {
				idle++
			}

			if idle > 150 {
				goto END
			}
			dino.Rotation += n

			ctrl.Yield()
		}

		pressed = 0
		for pressed < 30 {
			if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) ||
				inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) ||
				inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) ||
				inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
				pressed++
				idle = 0
				dino.DrawSize.Set(size.X*(1+float64(pressed)/5), size.X*(1+float64(pressed)/5))
				dino.Rotation += -0.1 + rand.Float64()*0.2
			} else {
				idle++
			}
			if idle > 200 {
				goto END
			}

			dino.Pos.X += -0.2 + rand.Float64()*0.3
			dino.Pos.Y += -0.2 + rand.Float64()*0.3

			ctrl.Yield()
		}

		for {
			if dino.DrawSize.X > size.X && dino.DrawSize.Y > size.Y {
				dino.DrawSize.SubXY(decreaseStep, decreaseStep)
				decreaseStep += 0.5
			} else {
				break
			}
			if dino.Rotation > 0 {
				dino.Rotation -= 0.55
			} else if dino.Rotation < 0 {
				dino.Rotation += 0.55
			}
			ctrl.Yield()
		}

		dino.DrawSize = size
		dino.Rotation = 0

		for {
			if ebiten.IsKeyPressed(ebiten.KeyLeft) {
				idle = 0
				dino.Flip |= 0b10
			} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
				idle = 0
				dino.Flip &^= 0b10
			}
			if ebiten.IsKeyPressed(ebiten.KeyUp) {
				idle = 0
				dino.Flip &^= 0b01
			} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
				idle = 0
				dino.Flip |= 0b01
			}

			dino.Pos.X += -0.2 + rand.Float64()*0.3
			dino.Pos.Y += -0.2 + rand.Float64()*0.3

			if ebiten.IsKeyPressed(ebiten.KeySpace) {
				if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
					dino.Vel.X = -70
				} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
					dino.Vel.X = 70
				}
				if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
					dino.Vel.Y = -70
				} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
					dino.Vel.Y = 70
				}

				break
			}

			idle++
			if idle > 100 {
				goto END
			}

			ctrl.Yield()
		}

		ctrl.Yield()

		for {
			dino.Pos.Add(&dino.Vel)
			dino.Vel.Scale(0.99)
			if dino.Vel.Length() < 5 || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
				break
			}
			ctrl.Yield()
		}

	END:
		jumpCharge = 0
		dino.DrawSize = size
		dino.Pos = pos
		dino.Vel.Scale(0)
		dino.Flip &^= 01
		dino.Actions.Add(dino.CollideWithTile)
		goto FALL
	} // ---------------------------------------------------------

FLY:
	{ // ---------------------------------------------------------
		println("fly")
		dino.animationScript.Transition(dino.AnimateFly)
		dino.Actions.Remove(dino.ApplyGravity)
		dino.Rotation = 0
		maxSpeed := float64(10)
		dino.Vel.Scale(0.80)
		for {
			if ebiten.IsKeyPressed(ebiten.KeyLeft) {
				dino.Flip = 0b10
				numsign.Set(&dino.Vel.X, -1)
				dino.Vel.X--
			} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
				dino.Flip = 0b00
				numsign.Set(&dino.Vel.X, 1)
				dino.Vel.X++
			} else {
				numsign.Set(&dino.Vel.X, 0)
			}

			if ebiten.IsKeyPressed(ebiten.KeyUp) {
				dino.Vel.Y += -1
			} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
				dino.Vel.Y += 1
			}

			dino.Vel.Scale(0.97)
			dino.Vel.ClampXY(-maxSpeed, -maxSpeed, maxSpeed, maxSpeed)

			dino.Pos.X += dino.Vel.X
			dino.Pos.Y += dino.Vel.Y

			if ebiten.IsKeyPressed(ebiten.KeySpace) && ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
				dino.Vel.Y = 0
				goto FALL
			}

			ctrl.Yield()

		}
	} // ---------------------------------------------------------

FALL:
	{ // ---------------------------------------------------------
		println("fall")
		dino.CurrentTileID = 12
		dino.animationScript.Cancel()
		dino.Actions.Add(dino.ApplyGravity)
		ctrl.Yield()

		for {
			leftDown := ebiten.IsKeyPressed(ebiten.KeyLeft)
			rightDown := ebiten.IsKeyPressed(ebiten.KeyRight)

			dirX := numsign.Get(dino.Vel.X)

			if math.Abs(dino.Vel.X) < 4.5 {
				if leftDown {
					dino.Flip = 0b10
					numsign.Set(&dino.Vel.X, -1)
					dino.Pos.X -= 1.0
				} else if rightDown {
					dino.Flip = 0b00
					numsign.Set(&dino.Vel.X, 1)
					dino.Pos.X += 1.0
				} else {
					numsign.Set(&dino.Vel.X, 0)
				}
			}

			swerve := (leftDown && dirX == 1) || (rightDown && dirX == -1)
			if swerve {
				dino.Vel.X *= 0.7
			}

			if jumps >= 2 {
				dino.Rotation += dirX * float64(jumps) / 10
			}

			if !dino.Hit.Some(0b1100) {
				dino.Pos.X += dino.Vel.X
			}

			if inpututil.IsKeyJustPressed(ebiten.KeySpace) && jumps < maxJumps {
				goto JUMP
			}

			if dino.Hit.Some(0b0001) {
				dino.Vel.Y = 0
				jumps = 0
				dino.Rotation = 0
				if math.Abs(dino.Vel.X) > 4.5 {
					goto RUN
				} else {
					goto WALK
				}
			}
			ctrl.Yield()
		}
	} // ---------------------------------------------------------

}
