// modified version of github.com/flywave/go3/float64/vec2
package rect

import (
	"fmt"
	"image"
	"math"

	"github.com/nvlled/carrot-example/vector"
)

// T is a coordinate system aligned rectangle defined by a Min and Max vector.
type T struct {
	// TODO: I don't think I will ever need to modify these
	// without deforming the rectangle, so no need to
	// expose these.
	Min vector.T
	Max vector.T
}

func Create(x, y, w, h float64) (rect T) {
	return T{
		Min: vector.Create(x, y),
		Max: vector.Create(x+w, y+h),
	}
}

func CreateInt(x, y, w, h int) (rect T) {
	return T{
		Min: vector.CreateInt(x, y),
		Max: vector.CreateInt(x+w, y+h),
	}
}

// New creates a Rect from two points.
func New(a, b *vector.T) (rect T) {
	rect.Min = vector.Min(a, b)
	rect.Max = vector.Max(a, b)
	return rect
}

// ParseRect parses a Rect from a string. See also String()
func ParseRect(s string) (r T, err error) {
	_, err = fmt.Sscan(s, &r.Min.X, &r.Min.Y, &r.Max.X, &r.Max.Y)
	return r, err
}

func FromImageRect(r image.Rectangle) T {
	return T{
		Min: vector.CreateInt(r.Min.X, r.Min.Y),
		Max: vector.CreateInt(r.Max.X, r.Max.Y),
	}
}

func (rect *T) Width() float64 {
	return rect.Max.X - rect.Min.X
}

func (rect *T) Height() float64 {
	return rect.Max.Y - rect.Min.Y
}

func (rect *T) Dimension() (width float64, height float64) {
	width = rect.Max.X - rect.Min.X
	height = rect.Max.Y - rect.Min.Y
	return
}

func (rect *T) Size() float64 {
	width := rect.Width()
	height := rect.Height()
	return math.Max(width, height)
}

// Slice returns the elements of the vector as slice.
func (rect *T) Slice() []float64 {
	return rect.Array()[:]
}

func (rect *T) Array() *[4]float64 {
	return &[...]float64{
		rect.Min.X, rect.Min.Y,
		rect.Max.X, rect.Max.Y,
	}
}

// String formats Rect as string. See also ParseRect().
func (rect *T) String() string {
	return rect.Min.String() + " " + rect.Max.String()
}

// ContainsPoint returns if a point is contained within the rectangle.
func (rect *T) ContainsPoint(p *vector.T) bool {
	return p.X >= rect.Min.X && p.X <= rect.Max.X &&
		p.Y >= rect.Min.Y && p.Y <= rect.Max.Y
}

// Contains returns if other Rect is contained within the rectangle.
func (rect *T) Contains(other *T) bool {
	return rect.Min.X <= other.Min.X &&
		rect.Min.Y <= other.Min.Y &&
		rect.Max.X >= other.Max.X &&
		rect.Max.Y >= other.Max.Y
}

// Area calculates the area of the rectangle.
func (rect *T) Area() float64 {
	return (rect.Max.X - rect.Min.X) * (rect.Max.Y - rect.Min.Y)
}

func (rect *T) Intersects(other *T) bool {
	return other.Max.X >= rect.Min.X &&
		other.Min.X <= rect.Max.X &&
		other.Max.Y >= rect.Min.Y &&
		other.Min.Y <= rect.Max.Y
}

// Join enlarges this rectangle to contain also the given rectangle.
func (rect *T) Join(other *T) {
	rect.Min = vector.Min(&rect.Min, &other.Min)
	rect.Max = vector.Max(&rect.Max, &other.Max)
}

func (rect *T) Extend(p *vector.T) {
	rect.Min = vector.Min(&rect.Min, p)
	rect.Max = vector.Max(&rect.Max, p)
}

func (rect *T) Mid() vector.T {
	w, h := rect.Dimension()
	return vector.AddXY(&rect.Min, w/2, h/2)
}
func (rect *T) MidX() float64 { return rect.Min.X + rect.Width()/2 }
func (rect *T) MidY() float64 { return rect.Min.Y + rect.Height()/2 }
func (rect *T) MidXY() (float64, float64) {
	return rect.Min.X + rect.Width()/2, rect.Min.Y + rect.Height()/2
}

func (rect *T) XY() (float64, float64) {
	return rect.Min.X, rect.Min.Y
}
func (rect *T) X() float64 { return rect.Min.X }
func (rect *T) Y() float64 { return rect.Min.Y }

// Scales the rectangle by vector.T{scaleX, scaleY}
// The rectangle's center point remains the same after scaling.
func (rect *T) Scale(pos *vector.T) {
	w, h := rect.Dimension()
	dx := (w*pos.X - w) / 2
	dy := (h*pos.Y - h) / 2
	rect.Min.AddXY(-dx, -dy)
	rect.Max.AddXY(dx, dy)
}
func (rect *T) ScaleXY(x, y float64) {
	w, h := rect.Dimension()
	dx := (w*x - w) / 2
	dy := (h*y - h) / 2
	rect.Min.AddXY(-dx, -dy)
	rect.Max.AddXY(dx, dy)
}

func (rect *T) Left() float64   { return rect.Min.X }
func (rect *T) Right() float64  { return rect.Max.X }
func (rect *T) Top() float64    { return rect.Min.Y }
func (rect *T) Bottom() float64 { return rect.Max.Y }

func (rect *T) Left_int() int   { return int(rect.Min.X) }
func (rect *T) Right_int() int  { return int(rect.Max.X) }
func (rect *T) Top_int() int    { return int(rect.Min.Y) }
func (rect *T) Bottom_int() int { return int(rect.Max.Y) }

func (rect *T) SetMid(pos *vector.T) {
	rect.SetMidX(pos.X)
	rect.SetMidY(pos.Y)
}

func (rect *T) SetMidXY(x, y float64) {
	rect.SetMidX(x)
	rect.SetMidY(y)
}

func (rect *T) SetMidX(x float64) {
	w := rect.Width()
	rect.Min.X = x - w/2
	rect.Max.X = x + w/2
}

func (rect *T) SetMidY(y float64) {
	h := rect.Height()
	rect.Min.Y = y - h/2
	rect.Max.Y = y + h/2
}

func (rect *T) SetLeft(x float64) {
	w := rect.Width()
	rect.Min.X = x
	rect.Max.X = x + w
}

func (rect *T) SetRight(x float64) {
	w := rect.Width()
	rect.Min.X = x - w
	rect.Max.X = x
}

func (rect *T) SetTop(y float64) {
	h := rect.Height()
	rect.Min.Y = y
	rect.Max.Y = y + h
}

func (rect *T) SetBottom(y float64) {
	h := rect.Height()
	rect.Min.Y = y - h
	rect.Max.Y = y
}

func (rect *T) SetTopLeft(v *vector.T) {
	rect.SetLeft(v.X)
	rect.SetTop(v.Y)
}

func (rect *T) SetTopLeftXY(x, y float64) {
	rect.SetLeft(x)
	rect.SetTop(y)
}

func (rect *T) SetBottomRight(v *vector.T) {
	rect.SetRight(v.X)
	rect.SetBottom(v.Y)
}

func (rect *T) Copy() T {
	copy := *rect
	return copy
}

// Joined returns the minimal rectangle containing both a and b.
func Joined(a, b *T) (rect T) {
	rect.Min = vector.Min(&a.Min, &b.Min)
	rect.Max = vector.Max(&a.Max, &b.Max)
	return rect
}
