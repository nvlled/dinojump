package bitf

import "strconv"

type T byte

func (x T) All(y byte) bool  { return byte(x)&y == y }
func (x T) Some(y byte) bool { return byte(x)&y != 0 }

func (x T) String() string {
	return strconv.FormatInt(int64(x), 2)
}
