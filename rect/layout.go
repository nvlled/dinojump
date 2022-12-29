package rect

type LayoutType struct{}

var Layout LayoutType

func (LayoutType) Left(rect, container *T)   { rect.SetLeft(container.Min.X) }
func (LayoutType) Right(rect, container *T)  { rect.SetRight(container.Max.X) }
func (LayoutType) Top(rect, container *T)    { rect.SetTop(container.Min.Y) }
func (LayoutType) Bottom(rect, container *T) { rect.SetBottom(container.Max.Y) }
func (LayoutType) Center(rect, container *T) {
	mid := container.Mid()
	rect.SetMid(&mid)
}

func (LayoutType) Restrict(rect, container *T) (bool, bool) {
	changedX, changedY := false, false
	if rect.Min.X < container.Min.X {
		rect.SetLeft(container.Min.X)
		changedX = true
	}
	if rect.Min.Y < container.Min.Y {
		rect.SetTop(container.Min.Y)
		changedY = true
	}
	if rect.Max.X > container.Max.X {
		rect.SetRight(container.Max.X)
		changedX = true
	}
	if rect.Max.Y > container.Max.Y {
		rect.SetBottom(container.Max.Y)
		changedY = true
	}
	return changedX, changedY
}

// Aligns rect with the container.
// The flags is as follows:
// 0b1000 - align left
// 0b1000 - align right
// 0b1100 - center horizontally
// 0b0010 - align top
// 0b0001 - align bottom
// 0b0011 - center vertically
//
// examples:
//
//	Align(r, container, 0b1000) - aligns left
//	Align(r, container, 0b1010) - aligns top-left
//	Align(r, container, 0b1110) - aligns top-center-X
func Align(rect, container *T, flags byte) {
	var centerX byte = 0b1100
	var centerY byte = 0b1100

	if flags&centerX == centerX {
		rect.SetMidX(container.MidX())
	} else {
		if flags&0b1000 != 0 {
			rect.SetLeft(container.Min.X)
		}
		if flags&0b0100 != 0 {
			rect.SetRight(container.Max.X)
		}
	}

	if flags&centerY == centerY {
		rect.SetMidY(container.MidY())
	} else {
		if flags&0b0010 != 0 {
			rect.SetTop(container.Min.Y)
		}
		if flags&0b0001 != 0 {
			rect.SetBottom(container.Max.Y)
		}
	}
}
