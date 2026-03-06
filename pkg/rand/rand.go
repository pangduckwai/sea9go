package rand

import (
	v2 "math/rand/v2"
	"time"

	"github.com/bytedance/gopkg/lang/fastrand"
)

// Seeded get a specific value from the given seed(s).
func Seeded(n int, seed ...uint64) int {
	if len(seed) < 1 {
		panic("Seed missing")
	} else if len(seed) == 1 {
		return v2.New(v2.NewPCG(uint64(time.Now().UnixNano()), uint64(seed[0]))).IntN(n)
	} else {
		return v2.New(v2.NewPCG(uint64(seed[1]), uint64(seed[0]))).IntN(n)
	}
}

// Rand pseudo-random number generator in a `games.Game`
type Rand interface {
	// Uint64 returns a pseudo-random 64-bit value as a uint64.
	Uint64() uint64

	// Intn returns, as an int, a non-negative pseudo-random number in [0,n). It panics if n <= 0.
	Intn(int) int

	// Perm returns, as a slice of n ints, a pseudo-random permutation of the integers in the half-open interval [0,n).
	Perm(int) []int
}

// RandFast uses bytedance's `fastrand`
type RandFast uint8

func New(id int) RandFast {
	return RandFast(id)
}

func (m RandFast) Uint64() uint64 {
	return fastrand.Uint64()
}

func (m RandFast) Intn(n int) int {
	return int(fastrand.Uint32n(uint32(n)))
}

func (m RandFast) Perm(n int) []int {
	return fastrand.Perm(n)
}
