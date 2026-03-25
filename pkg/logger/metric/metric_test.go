package metric

import (
	"fmt"
	"math"
	"testing"
)

var tESTS = []int64{
	56,
	345,
	5678,
	7890,
	8900,
	9000,
	56789,
	567890,
	5678901,
	56789012,
	567890123,
	2345678901,
	56789012345,
	234567890123,
	56789012345678,
	2345678901234567,
	56789012345678901,
	math.MaxInt64 - 3,
}

// TestDev for development
// func TestDev(t *testing.T) {
// 	for _, val := range tESTS {
// 		for idx, hilo := range hHLO_DEC {
// 			q, _ := divmodDec(val, hilo)
// 			if q <= 0 {
// 				fmt.Printf("TestDev() [%2v] <- %v\n", idx, val)
// 				break
// 			}
// 		}
// 	}
// }

// _decimal uses division twice
func _decimal(i int64, dec int) (o int64) {
	neg := false
	if i < 0 {
		neg = true
		i = -i
	}
	o = i

	var q, r int64
	if q = divDec(i, hILO_DEC[dec-1]); q > 0 {
		for idx, hilo := range hILO_DEC[dec:] {
			q = divDec(i, hilo)
			if q <= 0 {
				o, r = divmodDec(i, hILO_DEC[idx])
				if r > int64(hILO_DEC[idx][2]>>1) {
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

func _matched(si, c int64) (m bool) {
	if si == c {
		return true
	}

	s := fmt.Sprintf("%v", si)
	t := s[0 : len(s)-2]
	i := len(t)
loop:
	for {
		switch t[i-1] {
		case '0':
			i--
		case '.':
			fallthrough
		default:
			m = true
			break loop
		}
	}
	return
}

func TestDecimals(t *testing.T) {
	for _, val := range tESTS {
		c := _decimal(val, 1)
		d := decimal(val, 1)
		if !_matched(c, d) { //c != d {
			t.Fatalf("TestDecimals() 1 '%v' and '%v' mismatched", d, c)
		}
	}

	for _, val := range tESTS {
		c := _decimal(val, 4)
		d := decimal(val, 4)
		if !_matched(c, d) {
			t.Fatalf("TestDecimals() 4 '%v' and '%v' mismatched", d, c)
		}
		// fmt.Printf("TestDecimals() 4 %5v and %v matched\n", d, c)
	}

	for _, val := range tESTS {
		c := _decimal(val, 7)
		d := decimal(val, 7)
		if !_matched(c, d) {
			t.Fatalf("TestDecimals() 7 '%v' and '%v' mismatched", d, c)
		}
	}

	fmt.Println("TestDecimals() test successful")
}

func _control(inp int64, dec int) string {
	if inp < int64(hILO_DEC[sUFFIX[0].i][2]) {
		return fmt.Sprintf("%v", inp)
	}
	var rnd uint64 = 1
	for j := 0; j < dec; j++ {
		rnd *= 10
	}
	round := float64(rnd)
	i := 1
	for ; i < len(sUFFIX); i++ {
		if (uint64(inp) / hILO_DEC[sUFFIX[i].i][2]) <= 0 {
			return fmt.Sprintf("%v %v", math.Round((float64(inp)/float64(hILO_DEC[sUFFIX[i-1].i][2]))*round)/round, sUFFIX[i-1].s)
		}
	}
	return fmt.Sprintf("%v %v", math.Round((float64(inp)/float64(hILO_DEC[sUFFIX[i-1].i][2]))*round)/round, sUFFIX[i-1].s)
}

// func TestMetrics(t *testing.T) {
// 	var c, m string
// 	for _, val := range tESTS {
// 		c = _control(val, 0)
// 		m = Metric(val, 0)
// 		if c != m {
// 			t.Fatalf("TestMetrics() 0 '%v' and '%v' mismatched", m, c)
// 		}
// 		fmt.Printf("TestMetrics() 0 %19v -> %11v and %11v matched\n", val, m, c) // TEMP
// 	}

// 	for _, val := range tESTS {
// 		c = _control(val, 3)
// 		m = Metric(val, 3)
// 		if c != m {
// 			t.Fatalf("TestMetrics() 3 '%v' and '%v' mismatched", m, c)
// 		}
// 		fmt.Printf("TestMetrics() 3 %19v -> %11v and %11v matched\n", val, m, c)
// 	}

// 	for _, val := range tESTS {
// 		c = _control(val, 5)
// 		m = Metric(val, 5)
// 		if c != m {
// 			t.Fatalf("TestMetrics() 5 '%v' and '%v' mismatched", m, c)
// 		}
// 		fmt.Printf("TestMetrics() 5 %19v -> %11v and %11v matched\n", val, m, c)
// 	}

// 	fmt.Println("TestMetrics() test successful")
// }

func BenchmarkOnce3(b *testing.B) {
	b.ResetTimer()
	i := 0
	for b.Loop() {
		decimal(tESTS[i%11], 3)
		i++
	}
}

func BenchmarkOnce4(b *testing.B) {
	b.ResetTimer()
	i := 0
	for b.Loop() {
		decimal(tESTS[i%11], 4)
		i++
	}
}

func BenchmarkOnce5(b *testing.B) {
	b.ResetTimer()
	i := 0
	for b.Loop() {
		decimal(tESTS[i%11], 5)
		i++
	}
}

func BenchmarkOnce6(b *testing.B) {
	b.ResetTimer()
	i := 0
	for b.Loop() {
		decimal(tESTS[i%11], 6)
		i++
	}
}

func BenchmarkTwice3(b *testing.B) {
	b.ResetTimer()
	i := 0
	for b.Loop() {
		_decimal(tESTS[i%11], 3)
		i++
	}
}

func BenchmarkTwice4(b *testing.B) {
	b.ResetTimer()
	i := 0
	for b.Loop() {
		_decimal(tESTS[i%11], 4)
		i++
	}
}

func BenchmarkTwice5(b *testing.B) {
	b.ResetTimer()
	i := 0
	for b.Loop() {
		_decimal(tESTS[i%11], 5)
		i++
	}
}

func BenchmarkTwice6(b *testing.B) {
	b.ResetTimer()
	i := 0
	for b.Loop() {
		_decimal(tESTS[i%11], 6)
		i++
	}
}

func BenchmarkControl3(b *testing.B) {
	b.ResetTimer()
	i := 0
	for b.Loop() {
		_control(tESTS[i%11], 3)
		i++
	}
}

func BenchmarkMetric3(b *testing.B) {
	b.ResetTimer()
	i := 0
	for b.Loop() {
		Metric(tESTS[i%11], 3)
		i++
	}
}

func BenchmarkControl7(b *testing.B) {
	b.ResetTimer()
	i := 0
	for b.Loop() {
		_control(tESTS[i%11], 7)
		i++
	}
}

func BenchmarkMetric7(b *testing.B) {
	b.ResetTimer()
	i := 0
	for b.Loop() {
		Metric(tESTS[i%11], 7)
		i++
	}
}
