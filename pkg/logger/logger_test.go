package logger

import (
	"fmt"
	"math"
	"testing"
)

var vals = []uint64{345, 90000, 2, 567890, 1, 1232345, math.MaxUint64, 34958769857, 5678}

func TestDigits(t *testing.T) {
	var ctrl, dcnt int
	for _, val := range vals {
		ctrl = int(math.Log10(float64(val))) + 1
		dcnt = DigitCount(val)
		if dcnt != ctrl {
			t.Fatalf("TestDigits() %20v -> %2v != %v\n", val, dcnt, ctrl)
		}
		fmt.Printf("TestDigits() %20v -> %2v == %v\n", val, dcnt, ctrl)
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

func TestLogger(t *testing.T) {
	lbls := []string{"this", "is", "a", "test!!!"}
	msgs := []string{"How", "are", "you", "today", "?"}
	log, err, _ := Init()

	log1 := AddPrefix(AddPrefix(log, "TestLogger() - "), lbls...)
	for _, m := range msgs {
		log1(" ==> %v\n", m)
	}

	loge, _ := AddLabels(AddPrefix(err, "TestLogger() - "), 1, 7, 0, 5, 2, 2)
	for i, l := range lbls {
		for j, m := range msgs {
			loge(" ==> IDX:%v\n", l, m, i, j)
		}
	}
}
