package logging

import (
	"fmt"
	"math"
	"testing"
)

var vals = []uint64{345, 90000, 2, 567890, 1, 1232345, math.MaxUint64, 34958769857, 5678}

func TestDigits(t *testing.T) {
	for _, val := range vals {
		fmt.Printf("TestDigits() %14v -> %v\n", val, DigitCount(val))
	}
}

func TestLog10(t *testing.T) {
	for _, val := range vals {
		fmt.Printf("TestLog10()   %14v -> %v / %v\n", val, int(math.Log10(float64(val)))+1, DigitCount(val))
	}
}

func BenchmarkDigits(b *testing.B) {
	b.ResetTimer()
	i := 0
	for b.Loop() {
		DigitCount(vals[i%9])
		i++
	}
}

func BenchmarkLog10(b *testing.B) {
	b.ResetTimer()
	i := 0
	for b.Loop() {
		_ = int(math.Log10(float64(vals[i%9]))) + 1
		i++
	}
}
