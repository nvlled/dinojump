package dino_enums

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
	"github.com/nvlled/carrot-example/vector"
)

var dinoMaxSpeed float64 = 20

type DinoState int
type DinoAnimation int
type JumpChargeState int

type JumpChargeData struct {
	n            float64
	size         vector.T
	pos          vector.T
	decreaseStep float64
	idle         int
	pressed      float64
}

const (
	DinoStateIdle DinoState = iota + 1
	DinoStateWalk
	DinoStateRun
	DinoStateBrake
	DinoStateJump
	DinoStateFall
	DinoStateBounce
	DinoStateFly
	DinoStateJumpCharge
)

const (
	JumpChargeState1 JumpChargeState = iota + 1
	JumpChargeState2
	JumpChargeState3
	JumpChargeState4
	JumpChargeState5
	JumpChargeStateEnd
)

const (
	AnimationNone DinoAnimation = iota
	AnimationIdle
	AnimationWalk
	AnimationRun
	AnimationBrake
	AnimationJump
	AnimationFly
	AnimationFall
	AnimationBounce
	AnimationOuchie
)

type Sprite struct {
	sprite.T

	Hit bitf.T

	Level *level.T

	animation       DinoAnimation
	animate         bool
	animationStep   int
	animationDelay  int
	animationFrames seqiter.Iterator[int]

	state DinoState

	turns   int
	jumps   int
	preJump bool

	maxJumps    int
	maxFlySpeed float64

	jumpCharge      int
	jumpChargeState JumpChargeState
	jumpChargeData  JumpChargeData
}

func New(level *level.T) *Sprite {
	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(assets.DinoSpriteData))
	if err != nil {
		panic(err)
	}

	dino := &Sprite{
		T:     *sprite.New(img, 24, 1),
		Level: level,

		turns: 0,
		jumps: 0,

		maxJumps:    3,
		maxFlySpeed: 20,
	}
	dino.T.CollisionScale.X = 0.7
	dino.T.CollisionScale.Y = 0.7

	dino.Actions.Add(dino.ApplyGravity)
	dino.Actions.Add(dino.CollideWithTile)
	dino.transition(DinoStateIdle)

	return dino
}

func (dino *Sprite) Update() {
	dino.T.Update()
	dino.updateController()
	dino.updateAnimation()

}

func (dino *Sprite) updateAnimation() {
	dino.animationStep++
	if !dino.animate || dino.animationStep%dino.animationDelay != 0 {
		return
	}
	dino.CurrentTileID = dino.animationFrames.Next()

	if dino.animation == AnimationRun {
		if dino.Vel.X >= float64(dinoMaxSpeed) {
			dino.animationDelay = 1
		} else {
			dino.animationDelay = 3
		}
	}
}

func (dino *Sprite) updateController() {
	switch dino.state {
	case DinoStateIdle:
		dino.updateIdle()
	case DinoStateWalk:
		dino.updateWalk()
	case DinoStateBrake:
		dino.updateBrake()
	case DinoStateBounce:
		dino.updateBounce()
	case DinoStateRun:
		dino.updateRun()
	case DinoStateJump:
		dino.updateJump()
	case DinoStateJumpCharge:
		dino.updateJumpCharge()
	case DinoStateFly:
		dino.updateFly()
	case DinoStateFall:
		dino.updateFall()
	}
}

func (dino *Sprite) SetAnimation(animation DinoAnimation) {
	dino.animate = true
	dino.animation = animation
	dino.animationStep = 0

	switch animation {
	case AnimationOuchie:
		dino.animationFrames = seqiter.CreateSeqIterator(14, 16)
		dino.animationDelay = 3
	case AnimationIdle:
		dino.animationFrames = seqiter.CreateSeqIterator(0, 1, 2, 3)
		dino.animationDelay = 7
	case AnimationWalk:
		dino.animationFrames = seqiter.CreateSeqIterator(3, 4, 5, 6, 7, 8)
		dino.animationDelay = 10
	case AnimationRun:
		dino.animationFrames = seqiter.CreateSeqIterator(common.RangeSlice(18, 23)...)
		dino.animationDelay = 3
	case AnimationFly:
		dino.animationFrames = seqiter.CreateSeqIterator(17, 18)
		dino.animationDelay = 10
	default:
		dino.animate = false
	}
}
func (dino *Sprite) RestrictInWorld(carrot.Void) {
	r := dino.Level.GetRect()
	cx, _ := dino.RestrictWithin(&r)
	if cx {
		dino.Vel.X = 0
	}
}

func (dino *Sprite) CollideWithTile(carrot.Void) {
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

func (dino *Sprite) ApplyGravity(carrot.Void) {
	if dino.Hit.Some(0b0001) {
		dino.Vel.Y = 0
	} else {
		dino.Vel.Y += 0.5

	}
	dino.Pos.Y += dino.Vel.Y
}

func (dino *Sprite) transition(nextState DinoState) DinoState {
	switch nextState {
	case DinoStateIdle:
		println("idle")
		dino.SetAnimation(AnimationIdle)
	case DinoStateWalk:
		println("walk")
		dino.Vel.X = 0.5
		dino.SetAnimation(AnimationWalk)
	case DinoStateBrake:
		println("brake")
		dino.SetAnimation(AnimationWalk)
	case DinoStateBounce:
		println("bounce")
		dino.SetAnimation(AnimationOuchie)
		dino.Vel.X *= -0.8
		dino.Vel.Y = -4.5
	case DinoStateRun:
		println("run")
		dino.SetAnimation(AnimationRun)
	case DinoStateJump:
		println("jump")
		dino.jumps++
		dino.Actions.Remove(dino.ApplyGravity)

		dino.CurrentTileID = 12
		dino.SetAnimation(AnimationNone)
		dino.preJump = true

	case DinoStateJumpCharge:
		println("jump charge")
		dino.Actions.Remove(dino.ApplyGravity)
		dino.Actions.Remove(dino.CollideWithTile)

		dino.jumpChargeData.n = 0.5
		dino.jumpChargeData.size = dino.DrawSize
		dino.jumpChargeData.pos = dino.Pos
		dino.jumpChargeData.decreaseStep = float64(1)
		dino.jumpChargeData.idle = 0
		dino.jumpChargeData.pressed = float64(0)
		dino.jumpChargeState = JumpChargeState1

	case DinoStateFly:
		println("fly")
		dino.SetAnimation(AnimationNone)
		dino.Actions.Remove(dino.ApplyGravity)
		dino.Rotation = 0
		dino.Vel.Scale(0.80)
	case DinoStateFall:
		println("fall")
		dino.CurrentTileID = 12
		dino.SetAnimation(AnimationNone)
		dino.Actions.Add(dino.ApplyGravity)
	}

	dino.state = nextState
	return nextState
}

func (dino *Sprite) updateIdle() DinoState {
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
		return dino.transition(DinoStateFall)
	}

	if walk {
		return dino.transition(DinoStateWalk)
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		return dino.transition(DinoStateJump)
	}

	return 0
}

func (dino *Sprite) updateWalk() DinoState {
	oldDir := numsign.Get(dino.Vel.X)
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		dino.Flip = 0b10
		numsign.Set(&dino.Vel.X, -1)
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		dino.Flip = 0b00
		numsign.Set(&dino.Vel.X, 1)
	} else {
		return dino.transition(DinoStateIdle)
	}

	if !dino.Hit.Some(0b0001) {
		return dino.transition(DinoStateFall)
	}

	if oldDir != numsign.Get(dino.Vel.X) {
		dino.Vel.X *= 0.5
	}

	dino.Pos.X += dino.Vel.X
	dino.Vel.X += 0.1 * numsign.Get(dino.Vel.X)

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		return dino.transition(DinoStateJump)
	}
	if math.Abs(dino.Vel.X) > 4.5 {
		return dino.transition(DinoStateRun)
	}

	return 0
}

func (dino *Sprite) updateBrake() DinoState {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		dino.Flip = 0b10
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		dino.Flip = 0b00
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		return dino.transition(DinoStateJump)
	}
	if dino.Hit.Some(0b1100) && math.Abs(dino.Vel.X) >= float64(dinoMaxSpeed)*0.55 {
		return dino.transition(DinoStateBounce)
	}
	if math.Abs(dino.Vel.X) <= 1 {
		return dino.transition(DinoStateIdle)
	}
	dino.Pos.X += dino.Vel.X

	dirX := numsign.Get(dino.Vel.X)
	if (ebiten.IsKeyPressed(ebiten.KeyLeft) && dirX > 0) ||
		(ebiten.IsKeyPressed(ebiten.KeyRight) && dirX < 0) {
		dino.Vel.X *= 0.90
	} else {
		dino.Vel.X *= 0.97

	}

	return 0
}

func (dino *Sprite) updateBounce() DinoState {
	dino.Pos.X += dino.Vel.X
	dino.Pos.Y += dino.Vel.Y
	dino.Vel.X *= 0.9

	if dino.Vel.Y < 0 {
		dino.Vel.Y += 0.25
	}

	if math.Abs(dino.Vel.X) < 1 {
		dino.Vel.Y = 0
		dino.Vel.X = 0
		return dino.transition(DinoStateIdle)
	}

	return 0
}

func (dino *Sprite) updateRun() DinoState {
	leftDown := ebiten.IsKeyPressed(ebiten.KeyLeft)
	rightDown := ebiten.IsKeyPressed(ebiten.KeyRight)
	noDown := !leftDown && !rightDown
	dirX := numsign.Get(dino.Vel.X)
	brake := noDown || (leftDown && dirX == 1) || (rightDown && dirX == -1)

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		return dino.transition(DinoStateJump)
	}

	if !dino.Hit.Some(0b0001) {
		return dino.transition(DinoStateFall)
	}

	if brake {
		return dino.transition(DinoStateBrake)
	}

	if leftDown {
		dino.Flip = 0b10
		numsign.Set(&dino.Vel.X, -1)
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		dino.Flip = 0b00
		numsign.Set(&dino.Vel.X, 1)
	}

	if dino.Hit.Some(0b1100) && math.Abs(dino.Vel.X) >= float64(dinoMaxSpeed)*0.55 {
		return dino.transition(DinoStateBounce)
	}

	dino.Pos.X += dino.Vel.X

	if dino.Vel.X < float64(dinoMaxSpeed) && !dino.Hit.Some(0b1100) {
		dino.Vel.X += 0.2 * numsign.Get(dino.Vel.X)
	}

	return 0
}

func (dino *Sprite) updateJump() DinoState {
	if dino.preJump {
		dino.preJump = false
		dino.Vel.Y = -7.5
		dino.jumpCharge = 0
		return 0
	}

	leftDown := ebiten.IsKeyPressed(ebiten.KeyLeft)
	rightDown := ebiten.IsKeyPressed(ebiten.KeyRight)

	dirX := numsign.Get(dino.Vel.X)
	if math.Abs(dino.Vel.X) < float64(dinoMaxSpeed) {
		if leftDown {
			dino.Flip = 0b10
			numsign.Set(&dino.Vel.X, -1)
		} else if rightDown {
			dino.Flip = 0b00
			numsign.Set(&dino.Vel.X, 1)
		} else {
			numsign.Set(&dino.Vel.X, 0)
		}
	}
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

	dino.jumpCharge++
	if dino.jumpCharge >= 40 && ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		return dino.transition(DinoStateJumpCharge)
	}

	if dino.Hit.Some(0b0010) {
		dino.Vel.Y *= -0.3
		return dino.transition(DinoStateFall)
	}

	if math.Abs(dino.Vel.Y) < 1 {
		return dino.transition(DinoStateFall)
	}

	if dino.jumps >= dino.maxJumps && math.Abs(dino.Vel.X) >= float64(dinoMaxSpeed) {
		return dino.transition(DinoStateFly)
	} else if dino.jumps >= 2 {
		dino.Rotation += dirX * float64(dino.jumps) / 10
	}

	if dino.Hit.Some(0b1100) && math.Abs(dino.Vel.X) >= float64(dinoMaxSpeed)*0.55 {
		return dino.transition(DinoStateBounce)
	}

	dino.turns++
	return 0
}

func (dino *Sprite) updateJumpCharge() DinoState {
	data := &dino.jumpChargeData

	switch dino.jumpChargeState {
	case JumpChargeState1:
		{
			if data.pressed >= 10 {
				dino.jumpChargeState = JumpChargeState2
				return dino.state
			}
			if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
				data.idle = 0
				data.pressed++
				data.n += 0.05
			} else {
				data.idle++
			}

			if data.idle > 150 {
				data.pressed = 0
				dino.jumpChargeState = JumpChargeStateEnd
				return dino.state
			}
			dino.Rotation += data.n
		}

	case JumpChargeState2:
		{
			if data.pressed >= 30 {
				dino.jumpChargeState = JumpChargeState3
				return dino.state
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) ||
				inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) ||
				inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) ||
				inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
				data.pressed++
				data.idle = 0
				dino.DrawSize.Set(data.size.X*(1+float64(data.pressed)/5), data.size.X*(1+float64(data.pressed)/5))
				dino.Rotation += -0.1 + rand.Float64()*0.2
			} else {
				data.idle++
			}
			if data.idle > 200 {
				dino.jumpChargeState = JumpChargeStateEnd
				return dino.state
			}

			dino.Pos.X += -0.2 + rand.Float64()*0.3
			dino.Pos.Y += -0.2 + rand.Float64()*0.3
		}

	case JumpChargeState3:
		{
			if dino.DrawSize.X > data.size.X && dino.DrawSize.Y > data.size.Y {
				dino.DrawSize.SubXY(data.decreaseStep, data.decreaseStep)
				data.decreaseStep += 0.5
			} else {
				dino.DrawSize = data.size
				dino.Rotation = 0
				dino.jumpChargeState = JumpChargeState4
				return dino.state
			}
			if dino.Rotation > 0 {
				dino.Rotation -= 0.55
			} else if dino.Rotation < 0 {
				dino.Rotation += 0.55
			}
		}
	case JumpChargeState4:
		{
			if ebiten.IsKeyPressed(ebiten.KeyLeft) {
				data.idle = 0
				dino.Flip |= 0b10
			} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
				data.idle = 0
				dino.Flip &^= 0b10
			}
			if ebiten.IsKeyPressed(ebiten.KeyUp) {
				data.idle = 0
				dino.Flip &^= 0b01
			} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
				data.idle = 0
				dino.Flip |= 0b01
			}

			dino.Pos.X += -0.2 + rand.Float64()*0.3
			dino.Pos.Y += -0.2 + rand.Float64()*0.3

			if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
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

				dino.jumpChargeState = JumpChargeState5
				return dino.state
			}

			data.idle++
			if data.idle > 100 {
				dino.jumpChargeState = JumpChargeStateEnd
				return dino.state
			}

		}
	case JumpChargeState5:
		{
			dino.Pos.Add(&dino.Vel)
			dino.Vel.Scale(0.99)
			if dino.Vel.Length() < 5 || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
				dino.jumpChargeState = JumpChargeStateEnd
				return dino.state
			}
		}
	case JumpChargeStateEnd:
		{
			dino.jumpCharge = 0
			dino.DrawSize = data.size
			dino.Pos = data.pos
			dino.Vel.Scale(0)
			dino.Flip &^= 01
			dino.Actions.Add(dino.CollideWithTile)
			return dino.transition(DinoStateFall)
		}
	}
	return 0
}

func (dino *Sprite) updateFly() DinoState {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		numsign.Set(&dino.Vel.X, -1)
		dino.Flip = 0b10
		dino.Vel.X--
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		numsign.Set(&dino.Vel.X, 1)
		dino.Flip = 0b00
		dino.Vel.X++
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		dino.Vel.Y += -1
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		dino.Vel.Y += 1
	}

	dino.Vel.Scale(0.97)
	dino.Vel.ClampXY(-dinoMaxSpeed, -dinoMaxSpeed, dinoMaxSpeed, dinoMaxSpeed)

	dino.Pos.X += dino.Vel.X
	dino.Pos.Y += dino.Vel.Y

	if ebiten.IsKeyPressed(ebiten.KeySpace) && ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		dino.Vel.Y = 0
		return dino.transition(DinoStateFall)
	}

	return 0
}

func (dino *Sprite) updateFall() DinoState {
	leftDown := ebiten.IsKeyPressed(ebiten.KeyLeft)
	rightDown := ebiten.IsKeyPressed(ebiten.KeyRight)

	dirX := numsign.Get(dino.Vel.X)
	if math.Abs(dino.Vel.X) < 4.5 {
		if leftDown {
			dino.Flip = 0b10
			numsign.Set(&dino.Vel.X, -1)
		} else if rightDown {
			dino.Flip = 0b00
			numsign.Set(&dino.Vel.X, 1)
		} else {
			numsign.Set(&dino.Vel.X, 0)
		}
	}

	swerve := (leftDown && dirX == 1) || (rightDown && dirX == -1)
	if swerve {
		dino.Vel.X *= 0.7
	}

	if dino.jumps >= 2 {
		dino.Rotation += dirX * float64(dino.jumps) / 10
	}

	if !dino.Hit.Some(0b1100) {
		dino.Pos.X += dino.Vel.X
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && dino.jumps < dino.maxJumps {
		return dino.transition(DinoStateJump)
	}

	if dino.Hit.Some(0b0001) {
		dino.Vel.Y = 0
		dino.jumps = 0
		dino.Rotation = 0
		if math.Abs(dino.Vel.X) > 4.5 {
			return dino.transition(DinoStateRun)
		} else {
			return dino.transition(DinoStateWalk)
		}
	}

	return 0
}
