package rand

// Usage:
// > cd .../pkg/fast
// > go test -bench .

import (
	"testing"

	"github.com/bytedance/gopkg/lang/fastrand"
)

func fastRand(n int) int {
	return int(fastrand.Uint32n(uint32(n)))
}

var rnd Rand = RandFast(0)
var rndFast RandFast = RandFast(1)

func Benchmark2(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = rnd.Intn(2)
	}
}

func Benchmark4(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = rnd.Intn(4)
	}
}

func BenchmarkIface(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = rnd.Intn(6)
	}
}

func BenchmarkDirect(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = rndFast.Intn(6)
	}
}

func BenchmarkCtrlIntn(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = fastrand.Intn(6)
	}
}

func BenchmarkCtrlFast(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = fastRand(6)
	}
}

func BenchmarkCtrlUint32n(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = fastrand.Uint32n(uint32(6))
	}
}
