package lib

import (
	"testing"

	"github.com/gofiber/fiber/v2/utils"
)

func TestRound(t *testing.T) {
	var a, b float64 = 785333.8, 785334

	a = Round(a)
	utils.AssertEqual(t, b, a)
}

func TestPrettyRound(t *testing.T) {
	var zero, zeroExpected float64 = 0, 0
	zeroTest := PrettyRound(zero)
	utils.AssertEqual(t, zeroExpected, zeroTest)

	var normal, normalExpected float64 = 500000, 500000
	normalTest := PrettyRound(normal)
	utils.AssertEqual(t, normalExpected, normalTest)

	var comma, commaExpected float64 = 500000.33333, 500000.34
	commaTest := PrettyRound(comma)
	utils.AssertEqual(t, commaTest, commaExpected)

	var anotherComma, anotherCommaExpected float64 = 500000.666, 500000.67
	anotherCommaTest := PrettyRound(anotherComma)
	utils.AssertEqual(t, anotherCommaExpected, anotherCommaTest)

	var uglyLow, uglyLowExpected float64 = 25000.000001, 25000
	uglyLowTest := PrettyRound(uglyLow)
	utils.AssertEqual(t, uglyLowExpected, uglyLowTest)

	var uglyHigh, uglyHighExpected float64 = 19999.99996, 20000
	uglyHighTest := PrettyRound(uglyHigh)
	utils.AssertEqual(t, uglyHighExpected, uglyHighTest)

	var aHigh, aHighExpected float64 = 19999.1, 19999.1
	aHighTest := PrettyRound(aHigh)
	utils.AssertEqual(t, aHighExpected, aHighTest)

	var bHigh, bHighExpected float64 = 19999.11, 19999.11
	bHighTest := PrettyRound(bHigh)
	utils.AssertEqual(t, bHighExpected, bHighTest)

	var cHigh, cHighExpected float64 = 19999.111, 19999.12
	cHighTest := PrettyRound(cHigh)
	utils.AssertEqual(t, cHighExpected, cHighTest)

	var dHigh, dHighExpected float64 = 19999.991, 20000
	dHighTest := PrettyRound(dHigh)
	utils.AssertEqual(t, dHighExpected, dHighTest)

	var eHigh, eHighExpected float64 = 20000.01, 20000
	eHighTest := PrettyRound(eHigh)
	utils.AssertEqual(t, eHighExpected, eHighTest)
}

func TestPrettyRoundPtr(t *testing.T) {
	var n *float64
	res := PrettyRoundPtr(n)
	utils.AssertEqual(t, true, res == nil)

	n = Float64ptr(12000.2999)
	res = PrettyRoundPtr(n)
	utils.AssertEqual(t, true, res != nil)
}

func TestRoundPtr(t *testing.T) {
	var n *float64
	res := RoundPtr(n)
	utils.AssertEqual(t, true, res == nil)

	n = Float64ptr(12000.2999)
	res = RoundPtr(n)
	utils.AssertEqual(t, true, res != nil)
}
