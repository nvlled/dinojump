package numsign

import "math"

func Get(x float64) float64 {
	if x == 0 {
		return 0
	}
	return x / math.Abs(x)
}

func Set(x *float64, sign float64) {
	*x = math.Abs(*x) * sign
}
