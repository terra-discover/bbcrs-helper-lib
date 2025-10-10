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

// RoundSmart will limit comma to n digits
func RoundSmart(val float64) float64 {
	diff := math.Abs(val - math.Round(val))

	// number behind - is n
	if diff < 1e-4 {
		return math.Floor(val)
	}

	// might have to changed the pow used as well (not sure yet)
	factor := math.Pow(10, 4)
	return math.Round(val*factor) / factor
}

// will times percent with 100 and add the remaining 0 to b
func NormalizePercent(percent float64) (a, b float64) {
	a = percent * 100
	b = math.Pow(10, 4)
	return math.Round(a), b
}
