// Package rand is a wrapper of fast pseudo random values generation.
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

	// Uint32n returns, as an uint32, a non-negative pseudo-random number in [0,n). It panics if n <= 0.
	Uint32n(uint32) uint32

	// Perm returns, as a slice of n ints, a pseudo-random permutation of the integers in the half-open interval [0,n).
	Perm(int) []int

	// Intn returns, as an int, a non-negative pseudo-random number in [0,n). It panics if n <= 0.
	Intn(int) int
}

// randFast uses bytedance's `fastrand`
type randFast uint8

func New(id int) Rand {
	return randFast(id)
}

// Temp return the struct instead of the interface for testing
func Temp(id int) randFast {
	return randFast(id)
}

func (m randFast) Uint64() uint64 {
	return fastrand.Uint64()
}

func (m randFast) Uint32n(n uint32) uint32 {
	return fastrand.Uint32n(n)
}

func (m randFast) Perm(n int) []int {
	return fastrand.Perm(n)
}

func (m randFast) Intn(n int) int {
	return int(fastrand.Uint32n(uint32(n)))
}
