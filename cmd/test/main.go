package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bytedance/gopkg/lang/fastrand"
)

// fastRand returns a uniform value in [0,n)
func fastRand(n int) int {
	return int(fastrand.Uint32n(uint32(n)))
}

func sim(id, n, run int) (lps time.Duration, cnt []int, nmz []float32) {
	cnt = make([]int, 0)
	for range n {
		cnt = append(cnt, 0)
	}

	now := time.Now()
	for range run {
		cnt[fastRand(n)] += 1
	}
	lps = time.Since(now)

	nmz = make([]float32, 0)
	for _, c := range cnt {
		nmz = append(nmz, float32(c)/float32(run))
	}

	return
}

func main() {
	var err error
	run := 1000000000 // 1,000,000,000
	rng := 6
	switch len(os.Args) {
	case 3:
		run, err = strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		fallthrough
	case 2:
		rng, err = strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
	case 1:
	default:
		log.Println("Usage: cmd/fast [range] [num-of-runs]")
	}

	lps, cnt, nmz := sim(1, rng, run)

	var buf strings.Builder
	for i, v := range cnt {
		fmt.Fprintf(&buf, " %3v: %v (%.4f%%)\n", i, v, nmz[i]*100)
	}

	// prt := message.NewPrinter(language.English)
	fmt.Printf("[FAST] %v simulations with [0,%v) range, elapsed time: %12v (%v per op)\n%v", run, rng, lps, lps/time.Duration(run), buf.String())
}
