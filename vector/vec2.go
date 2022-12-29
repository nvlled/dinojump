// Package vec2 contains a 2D float64 vector type T and functions.
// modified version vec2 from github.com/flywave/go3d
// changes:
// - uses struct XY instead of slice
// - added Div()
package vector

import (
	"fmt"
	"math"
)

var (
	// Zero holds a zero vector.
	Zero = T{}

	// UnitX holds a vector with X set to one.
	UnitX = T{1, 0}
	// UnitY holds a vector with Y set to one.
	UnitY = T{0, 1}
	// Unit holds a vector with X and Y set to one.
	Unit = T{1, 1}

	// MinVal holds a vector with the smallest possible component values.
	MinVal = T{-math.MaxInt, -math.MaxInt}
	// MaxVal holds a vector with the highest possible component values.
	MaxVal = T{+math.MaxInt, +math.MaxInt}
)

// T represents a 2D vector.
type T struct {
	X float64
	Y float64
}

func Create(x, y float64) T {
	return T{X: x, Y: y}
}

func CreateInt(x, y int) T {
	return T{X: float64(x), Y: float64(y)}
}

// Parse parses T from a string. See also String()
func Parse(s string) (r T, err error) {
	_, err = fmt.Sscan(s, &r.X, &r.Y)
	return r, err
}

func (vec *T) Set(x, y float64) {
	vec.X = x
	vec.Y = y
}

// String formats T as string. See also Parse().
func (vec *T) String() string {
	return fmt.Sprint(vec.X, vec.Y)
}

func (vec *T) XY() (float64, float64) {
	return vec.X, vec.Y
}

func (vec *T) XY_int() (int, int) {
	return int(vec.X), int(vec.Y)
}

// Rows returns the number of rows of the vector.
func (vec *T) Rows() int {
	return 2
}

// Cols returns the number of columns of the vector.
func (vec *T) Cols() int {
	return 1
}

// Size returns the number elements of the vector.
func (vec *T) Size() int {
	return 2
}

// Slice returns the elements of the vector as slice.
func (vec *T) Slice() []float64 {
	return []float64{vec.X, vec.Y}
}

// Get returns one element of the vector.
func (vec *T) Get(col, row int) float64 {
	if row == 0 {
		return vec.X
	}
	if row == 1 {
		return vec.Y
	}
	return 0
}

// IsZero checks if all elements of the vector are zero.
func (vec *T) IsZero() bool {
	return vec.X == 0 && vec.Y == 0
}

// Length returns the length of the vector.
// See also LengthSqr and Normalize.
func (vec *T) Length() float64 {
	return math.Hypot(float64(vec.X), float64(vec.Y))
}

// LengthSqr returns the squared length of the vector.
// See also Length and Normalize.
func (vec *T) LengthSqr() float64 {
	return vec.X*vec.X + vec.Y*vec.Y
}

// Scale multiplies all element of the vector by f and returns vec.
func (vec *T) Scale(f float64) *T {
	vec.X = float64(vec.X) * f
	vec.Y = float64(vec.Y) * f
	return vec
}

func (vec *T) ScaleXY(x, y float64) {
	vec.X = float64(vec.X) * x
	vec.Y = float64(vec.Y) * y
}

// Scaled returns a copy of vec with all elements multiplies by f.
func (vec *T) Scaled(f float64) T {
	copy := *vec
	copy.Scale(f)
	return copy
}

// Invert inverts the vector.
func (vec *T) Invert() *T {
	vec.X = -vec.X
	vec.Y = -vec.Y
	return vec
}

// Inverted returns an inverted copy of the vector.
func (vec *T) Inverted() T {
	return T{-vec.X, -vec.Y}
}

// Normalize normalizes the vector to unit length.
func (vec *T) Normalize() *T {
	sl := vec.LengthSqr()
	if sl == 0 || sl == 1 {
		return vec
	}
	return vec.Scale(1 / math.Sqrt(sl))
}

// Normalized returns a unit length normalized copy of the vector.
func (vec *T) Normalized() T {
	v := *vec
	v.Normalize()
	return v
}

// Add adds another vector to vec.
func (vec *T) Add(v *T) *T {
	vec.X += v.X
	vec.Y += v.Y
	return vec
}

func (vec *T) AddXY(dx, dy float64) *T {
	vec.X += dx
	vec.Y += dy
	return vec
}

func (vec *T) SubXY(x, y float64) {
	vec.X -= x
	vec.Y -= y
}

// Sub subtracts another vector from vec.
func (vec *T) Sub(v *T) *T {
	vec.X -= v.X
	vec.Y -= v.Y
	return vec
}

// Mul multiplies the components of the vector with the respective components of v.
func (vec *T) Mul(v *T) *T {
	vec.X *= v.X
	vec.Y *= v.Y
	return vec
}

// Div divides the components of the vector with the respective components of v.
func (vec *T) Div(v *T) *T {
	vec.X /= v.X
	vec.Y /= v.Y
	return vec
}

// Same as Div(), but checks for zero division.
// A component will be set to 0 if divided by zero.
func (vec *T) DivSafe(v *T) *T {
	if v.X != 0 {
		vec.X /= v.X
	} else {
		vec.X = 0
	}

	if v.Y != 0 {
		vec.Y /= v.Y
	} else {
		vec.Y = 0
	}

	return vec
}

// Rotate rotates the vector counter-clockwise by angle.
func (vec *T) Rotate(angle float64) *T {
	*vec = vec.Rotated(angle)
	return vec
}

// Rotated returns a counter-clockwise rotated copy of the vector.
func (vec *T) Rotated(angle float64) T {
	sinus := math.Sin(angle)
	cosinus := math.Cos(angle)
	return T{
		vec.X*cosinus - vec.Y*sinus,
		vec.X*sinus + vec.Y*cosinus,
	}
}

// RotateAroundPoint rotates the vector counter-clockwise around a point.
func (vec *T) RotateAroundPoint(point *T, angle float64) *T {
	return vec.Sub(point).Rotate(angle).Add(point)
}

// Rotate90DegLeft rotates the vector 90 degrees left (counter-clockwise).
func (vec *T) Rotate90DegLeft() *T {
	temp := vec.X
	vec.X = -vec.Y
	vec.Y = temp
	return vec
}

// Rotate90DegRight rotates the vector 90 degrees right (clockwise).
func (vec *T) Rotate90DegRight() *T {
	temp := vec.X
	vec.X = vec.Y
	vec.Y = -temp
	return vec
}

// Angle returns the counter-clockwise angle of the vector from the x axis.
func (vec *T) Angle() float64 {
	return math.Atan2(vec.Y, vec.X)
}

// Add returns the sum of two vectors.
func Add(a, b *T) T {
	return T{a.X + b.X, a.Y + b.Y}
}

func AddXY(a *T, dx, dy float64) T {
	return T{a.X + dx, a.Y + dy}
}

// Sub returns the difference of two vectors.
func Sub(a, b *T) T {
	return T{a.X - b.X, a.Y - b.Y}
}

// Mul returns the component wise product of two vectors.
func Mul(a, b *T) T {
	return T{a.X * b.X, a.Y * b.Y}
}

// Dot returns the dot product of two vectors.
func Dot(a, b *T) float64 {
	return a.X*b.X + a.Y*b.Y
}

// Cross returns the cross product of two vectors.
func Cross(a, b *T) T {
	return T{
		a.Y*b.X - a.X*b.Y,
		a.X*b.Y - a.Y*b.X,
	}
}

// Angle returns the angle between two vectors.
func Angle(a, b *T) float64 {
	v := Dot(a, b) / (a.Length() * b.Length())
	// prevent NaN
	if v > 1. {
		v = v - 2
	} else if v < -1. {
		v = v + 2
	}
	return math.Acos(v)
}

// IsLeftWinding returns if the angle from a to b is left winding.
func IsLeftWinding(a, b *T) bool {
	ab := b.Rotated(-a.Angle())
	return ab.Angle() > 0
}

// IsRightWinding returns if the angle from a to b is right winding.
func IsRightWinding(a, b *T) bool {
	ab := b.Rotated(-a.Angle())
	return ab.Angle() < 0
}

// Min returns the component wise minimum of two vectors.
func Min(a, b *T) T {
	min := *a
	if b.X < min.X {
		min.X = b.X
	}
	if b.Y < min.Y {
		min.Y = b.Y
	}
	return min
}

// Max returns the component wise maximum of two vectors.
func Max(a, b *T) T {
	max := *a
	if b.X > max.X {
		max.X = b.X
	}
	if b.Y > max.Y {
		max.Y = b.Y
	}
	return max
}

// Interpolate interpolates between a and b at t (0,1).
func Interpolate(a, b *T, t float64) T {
	t1 := 1 - t
	return T{
		a.X*t1 + b.X*t,
		a.Y*t1 + b.Y*t,
	}
}

func (vec *T) ClampXY(minX, minY, maxX, maxY float64) {
	if vec.X < minX {
		vec.X = minX
	} else if vec.X > maxX {
		vec.X = maxX
	}
	if vec.Y < minY {
		vec.Y = minY
	} else if vec.Y > maxY {
		vec.Y = maxY
	}
}

// Clamp clamps the vector's components to be in the range of min to max.
func (vec *T) Clamp(min, max *T) *T {
	if vec.X < min.X {
		vec.X = min.X
	} else if vec.X > max.X {
		vec.X = max.X
	}
	if vec.Y < min.Y {
		vec.Y = min.Y
	} else if vec.Y > max.Y {
		vec.Y = max.Y
	}
	return vec
}

// Clamped returns a copy of the vector with the components clamped to be in the range of min to max.
func (vec *T) Clamped(min, max *T) T {
	result := *vec
	result.Clamp(min, max)
	return result
}

// Clamp01 clamps the vector's components to be in the range of 0 to 1.
func (vec *T) Clamp01() *T {
	return vec.Clamp(&Zero, &Unit)
}

// Clamped01 returns a copy of the vector with the components clamped to be in the range of 0 to 1.
func (vec *T) Clamped01() T {
	result := *vec
	result.Clamp01()
	return result
}

func (vec *T) SetMin(c T) {
	if c.X < vec.X {
		vec.X = c.X
	}
	if c.Y < vec.Y {
		vec.Y = c.Y
	}
}
func (vec *T) SetMax(c T) {
	if c.X > vec.X {
		vec.X = c.X
	}
	if c.Y > vec.Y {
		vec.Y = c.Y
	}
}
