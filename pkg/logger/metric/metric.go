package metric

import (
	"fmt"
	"math/bits"

	logging "github.com/pangduckwai/sea9go/pkg/logger"
)

// hILO_DEC to generate the hi/lo values for dividing by 10^x:
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
var hILO_DEC = [][]uint64{
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

// sUFFIX list of metric suffices of values can be fit into an int64
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
	return
}

// divDec reference: https://github.com/bmkessler/fastdiv/
func divDec(n int64, hilo []uint64) int64 {
	l1, _ := bits.Mul64(hilo[1], uint64(n))
	rst, l2 := bits.Mul64(hilo[0], uint64(n))
	_, c := bits.Add64(l1, l2, 0)
	rst, _ = bits.Add64(rst, 0, c)
	return int64(rst)
}

// modDec reference: https://github.com/bmkessler/fastdiv/
func modDec(n int64, hilo []uint64) int64 {
	hi, lo := bits.Mul64(hilo[1], uint64(n))
	hi += hilo[0] * uint64(n)
	l1, _ := bits.Mul64(lo, hilo[2])
	rst, l2 := bits.Mul64(hi, hilo[2])
	_, c := bits.Add64(l1, l2, 0)
	rst, _ = bits.Add64(rst, 0, c)
	return int64(rst)
}

// decimal round off the remainder
func decimal(i int64, dec int) (q int64) {
	q = i
	var r int64
	idx := logging.DigitCount(uint64(i)) - 1
	if idx >= dec {
		q, r = divmodDec(i, hILO_DEC[idx-dec])
		if r > int64(hILO_DEC[idx-dec][2]>>1) {
			q++
		}
	}

	idx = logging.DigitCount(uint64(q)) - 2
	for ; idx >= 0; idx-- {
		p := divDec(q, hILO_DEC[idx])
		if q == p*int64(hILO_DEC[idx][2]) {
			return p
		}
	}
	return
}

// round decide if the given value >= 5xx...
func round(i int64) int64 {
	idx := logging.DigitCount(uint64(i)) - 1
	r := modDec(i, hILO_DEC[idx])
	if r > int64(hILO_DEC[idx][2]>>1) {
		return 1 // round up
	}
	return 0 // round down
}

func __metric(inp int64, dec, idx int) (int64, int64) {
	q, r := divmodDec(inp, hILO_DEC[sUFFIX[idx].i])
	if dec > 0 {
		if r > 0 {
			return q, decimal(r, dec)
		}
		return q, -1
	}
	q += round(r)
	return q, -1
}

func _metric(neg string, inp int64, dec, idx int) string {
	q, r := __metric(inp, dec, idx)
	if r >= 0 {
		return fmt.Sprintf("%v%v.%v %v", neg, q, r, sUFFIX[idx].s)
	}
	return fmt.Sprintf("%v%v %v", neg, q, sUFFIX[idx].s)
}

// Metric convert input to metric suffix with the given decimal places.
func Metric(inp int64, dec int) string {
	if dec < 0 {
		panic("No. of decimal points cannot be negative")
	}

	neg := ""
	if inp < 0 {
		neg = "-"
		inp = -inp
	}

	i, k := logging.DigitCount(uint64(inp))-2, len(sUFFIX)-1
	if i < sUFFIX[0].i {
		return fmt.Sprintf("%v%v", neg, inp)
	}

	for j, s := range sUFFIX[1:] {
		if i < s.i {
			return _metric(neg, inp, dec, j)
		}
	}
	return _metric(neg, inp, dec, k)
}
