package metric

import (
	"fmt"
	"math/bits"
)

// hHLO_DEC to generate the hi/lo values for dividing by 10^x:
//
//	const oVER float64 = float64(math.MaxUint64) / float64(10)
//	var a, h, l, r, c uint64
//	fmt.Println("[][]uint64{")
//	for a = 10; ; a *= 10 {
//		h, l, c = 0, 0, 1
//		h, r = math.MaxUint64/a, math.MaxUint64%a
//		l, _ = bits.Div64(r, math.MaxUint64, a)
//		if a&(a-1) == 0 {
//			c++
//		}
//		l, c = bits.Add64(l, c, 0)
//		h, _ = bits.Add64(h, 0, c)
//		fmt.Printf("\t{%v, %v, %v},\n", h, l, a)
//		if float64(a) > oVER {
//			break
//		}
//	}
//	fmt.Println("}")
var hHLO_DEC = [][]uint64{
	{1844674407370955161, 11068046444225730970, 10},
	{184467440737095516, 2951479051793528259, 100},
	{18446744073709551, 11363194349405083796, 1000}, // k 10^3
	{1844674407370955, 2980993842311463542, 10000},
	{184467440737095, 9521471421085922163, 100000},
	{18446744073709, 10175519178963368025, 1000000}, // M 10^6
	{1844674407370, 17619621584234933257, 10000000},
	{184467440737, 1761962158423493326, 100000000},
	{18446744073, 13088917067439035464, 1000000000}, // G 10^9
	{1844674407, 6842914928856769032, 10000000000},
	{184467440, 13597012344482363035, 100000000000},
	{18446744, 1359701234448236304, 1000000000000}, // T 10^12
	{1844674, 7514667752928644277, 10000000000000},
	{184467, 8130164404776685075, 100000000000000},
	{18446, 13725737292074354639, 1000000000000000}, // P 10^15
	{1844, 12440620173433166434, 10000000000000000},
	{184, 8622759646827137290, 100000000000000000},
	{18, 8240973594166534376, 1000000000000000000}, // E 10^18
	{1, 15581492618384294731, 10000000000000000000},
}

var sUFFIX = []struct {
	s string
	i int
}{
	{"k", 2},  // kilo
	{"M", 5},  // mega
	{"G", 8},  // giga
	{"T", 11}, // tera
	{"P", 14}, // peta
	{"E", 17}, // exa
}

// divmodDec reference: https://github.com/bmkessler/fastdiv/
func divmodDec(n int64, hilo []uint64) (q, r int64) {
	neg := false
	if n < 0 {
		n = -n
		neg = true
	}
	u := uint64(n)

	d1, lo := bits.Mul64(hilo[1], u)
	div, d2 := bits.Mul64(hilo[0], u)
	hi, c := bits.Add64(d1, d2, 0)
	div, _ = bits.Add64(div, 0, c)
	q = int64(div)

	m1, _ := bits.Mul64(lo, hilo[2])
	mod, m2 := bits.Mul64(hi, hilo[2])
	_, c = bits.Add64(m1, m2, 0)
	mod, _ = bits.Add64(mod, 0, c)
	r = int64(mod)

	if neg {
		q = -q
		r = -r
	}
	return
}

// decimal round off the remainder
func decimal(i int64, dec int) (o int64) {
	neg := false
	if i < 0 {
		neg = true
		i = -i
	}
	o = i
	if i >= int64(hHLO_DEC[dec-1][2]) {
		var q, r int64
		qt, rt, ot := make([]int64, 0), make([]int64, 0), make([]uint64, 0)
		for idx, hilo := range hHLO_DEC {
			q, r = divmodDec(i, hilo)
			if q <= 0 {
				if rt[0] > int64(ot[0])>>1 {
					o = qt[0] + 1
				} else {
					o = qt[0]
				}
				break
			}
			if idx < dec {
				qt, rt, ot = append(qt, q), append(rt, r), append(ot, hilo[2])
			} else {
				qt, rt, ot = append(qt[1:], q), append(rt[1:], r), append(ot[1:], hilo[2])
			}
		}
	}
	if neg {
		o = -o
	}
	return
}

// Metric convert input to metric suffix with the given decimal places.
func Metric(inp int64, dec int) string {
	if inp < int64(hHLO_DEC[sUFFIX[0].i][2]) {
		return fmt.Sprintf("%v", inp)
	}

	var q, r int64
	var st string = sUFFIX[0].s
	qt, rt := divmodDec(inp, hHLO_DEC[sUFFIX[0].i])
	for _, s := range sUFFIX[1:] {
		q, r = divmodDec(inp, hHLO_DEC[s.i])
		if q <= 0 {
			break
		}
		qt, rt, st = q, r, s.s
	}

	// if dec > 0 {
	return fmt.Sprintf("%v.%v %v", qt, decimal(rt, dec), st)
	// }
}
