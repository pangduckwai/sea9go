package rand

// Usage:
// $ cd .../rand
// $ go test -bench . -benchmem

import (
	"testing"

	"github.com/bytedance/gopkg/lang/fastrand"
)

func fastRand(n int) int {
	return int(fastrand.Uint32n(uint32(n)))
}

var rnd Rand = randFast(0)
var rndFast randFast = randFast(1)

func Benchmark2(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = rnd.Uint32n(2)
	}
}

func Benchmark4(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = rnd.Uint32n(4)
	}
}

func BenchmarkIface(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = rnd.Uint32n(6)
	}
}

func BenchmarkStru(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = rndFast.Uint32n(6)
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
		_ = fastrand.Uint32n(6)
	}
}
