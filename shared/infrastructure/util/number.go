package util

import (
	"math"
	"strconv"
)

func DigitGrouping(n int64) string {

	in := strconv.FormatInt(n, 10)
	numOfDigits := len(in)
	if n < 0 {
		numOfDigits-- // First character is the - sign (not a digit)
	}
	numOfCommas := (numOfDigits - 1) / 3

	out := make([]byte, len(in)+numOfCommas)
	if n < 0 {
		in, out[0] = in[1:], '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = '.'
		}
	}
}

func RoundUpDiscount(disc float64) float64 {
	x := math.Round(disc)

	return x
}

func RoundDownDiscount(disc float64) float64 {
	x := math.Floor(disc)

	if x == 0 {
		if disc > 0 {
			x = 1
		}
	}

	return x
}
