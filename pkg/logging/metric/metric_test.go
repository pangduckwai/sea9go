package metric

import (
	"fmt"
	"math"
	"math/bits"
	"testing"
	"time"
)

var tESTS = []int64{
	56,
	5678,
	56789,
	567890,
	5678901,
	56789012,
	567890123,
	56789012345,
	56789012345678,
	56789012345678901,
	math.MaxInt64 - 3,
}

func TestDev(t *testing.T) {
	for _, val := range tESTS {
		for idx, hilo := range hHLO_DEC {
			q, _ := divmodDec(val, hilo)
			if q <= 0 {
				fmt.Printf("TestDev() [%2v] <- %v\n", idx, val)
				break
			}
		}
	}
}

// divDec reference: https://github.com/bmkessler/fastdiv/
func divDec(n int64, hilo []uint64) int64 {
	neg := false
	if n < 0 {
		n = -n
		neg = true
	}

	l1, _ := bits.Mul64(hilo[1], uint64(n))
	rst, l2 := bits.Mul64(hilo[0], uint64(n))
	_, c := bits.Add64(l1, l2, 0)
	rst, _ = bits.Add64(rst, 0, c)

	if neg {
		return -int64(rst)
	}
	return int64(rst)
}

// mod reference: https://github.com/bmkessler/fastdiv/
// func modDec(n int64, hilo []uint64) int64 {
// 	neg := false
// 	if n < 0 {
// 		n = -n
// 		neg = true
// 	}

// 	hi, lo := bits.Mul64(hilo[1], uint64(n))
// 	hi += hilo[0] * uint64(n)
// 	l1, _ := bits.Mul64(lo, hilo[2])
// 	rst, l2 := bits.Mul64(hi, hilo[2])
// 	_, c := bits.Add64(l1, l2, 0)
// 	rst, _ = bits.Add64(rst, 0, c)

// 	if neg {
// 		return -int64(rst)
// 	}
// 	return int64(rst)
// }

func _decimal(i int64, dec int) (o int64) {
	neg := false
	if i < 0 {
		neg = true
		i = -i
	}
	o = i

	var q, r int64
	if q = divDec(i, hHLO_DEC[dec-1]); q > 0 {
		for idx, hilo := range hHLO_DEC[dec:] {
			q = divDec(i, hilo)
			if q <= 0 {
				o, r = divmodDec(i, hHLO_DEC[idx])
				if r > int64(hHLO_DEC[idx][2]>>1) {
					o++
				}
				break
			}
		}
	}
	if neg {
		o = -o
	}
	return
}

func __decimal(i int64, dec int) (o int64) {
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

func TestDeci(t *testing.T) {
	for _, val := range tESTS {
		fmt.Printf("TestDeci() %5v <- %v\n", __decimal(val, 4), val)
	}
}

func TestMetric(t *testing.T) {
	// Metric
	for _, val := range tESTS {
		fmt.Printf("TestMetric() %19v : %v\n", val, Metric(val, 3))
	}

	now0 := time.Now()
	for range 1000000 {
		for _, val := range tESTS {
			Metric(val, 1)
		}
	}
	fmt.Printf("TestMetric() elapsed: %v\n", time.Since(now0))
}

func BenchmarkOnce3(b *testing.B) {
	b.ResetTimer()
	i := 0
	for b.Loop() {
		decimal(tESTS[i%11], 3)
	}
}

func BenchmarkOnce4(b *testing.B) {
	b.ResetTimer()
	i := 0
	for b.Loop() {
		decimal(tESTS[i%11], 4)
	}
}

func BenchmarkTwice3(b *testing.B) {
	b.ResetTimer()
	i := 0
	for b.Loop() {
		_decimal(tESTS[i%11], 3)
	}
}

func BenchmarkTwice4(b *testing.B) {
	b.ResetTimer()
	i := 0
	for b.Loop() {
		_decimal(tESTS[i%11], 4)
	}
}

func BenchmarkOnce5(b *testing.B) {
	b.ResetTimer()
	i := 0
	for b.Loop() {
		decimal(tESTS[i%11], 5)
	}
}

func BenchmarkOnce6(b *testing.B) {
	b.ResetTimer()
	i := 0
	for b.Loop() {
		decimal(tESTS[i%11], 6)
	}
}

func BenchmarkTwice5(b *testing.B) {
	b.ResetTimer()
	i := 0
	for b.Loop() {
		_decimal(tESTS[i%11], 5)
	}
}

func BenchmarkTwice6(b *testing.B) {
	b.ResetTimer()
	i := 0
	for b.Loop() {
		_decimal(tESTS[i%11], 6)
	}
}
