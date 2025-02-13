 package lib

import (
	"math"
)

// Round Ceil returns the least integer value greater than or equal to x
func Round(a float64) (b float64) {
	return math.Ceil(a)
}

// PrettyRound will limit after comma to 2 digits, and remove unneeded 0000 trail if exists (check in test case)
func PrettyRound(a float64) (b float64) {
	a = a * 10
	floored := math.Floor(a)
	if (a)-floored > 0 {
		if (a-floored)*100 >= 99 {
			return (floored + 1) / 10
		} else if (a-floored)*10 < 1 {
			return (floored / 10)
		} else {
			// round with 2 after comma precision
			return (math.Ceil(a*10) / 100)
		}
	}
	// number already floored
	return a / 10
}

func PrettyRoundPtr(a *float64) (b *float64) {
	if a == nil {
		return nil
	}

	rounded := PrettyRound(*a)
	return &rounded
}

func RoundPtr(a *float64) (b *float64) {
	if a == nil {
		return nil
	}
	rounded := Round(*a)
	return &rounded
}