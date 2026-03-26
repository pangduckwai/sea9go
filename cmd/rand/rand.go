package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"

	"github.com/bytedance/gopkg/lang/fastrand"
	"github.com/pangduckwai/sea9go/pkg/rand"
)

func simDirect(n, run int) (lps time.Duration, cnt []int, nmz []float32) {
	o := uint32(n)
	cnt = make([]int, 0)
	for range n {
		cnt = append(cnt, 0)
	}
	now := time.Now()
	for range run {
		cnt[fastrand.Uint32n(o)] += 1
	}
	lps = time.Since(now)
	nmz = make([]float32, 0)
	for _, c := range cnt {
		nmz = append(nmz, float32(c)/float32(run))
	}
	return
}

func simCtrl(n, run int) (lps time.Duration, cnt []int, nmz []float32) {
	cnt = make([]int, 0)
	for range n {
		cnt = append(cnt, 0)
	}
	now := time.Now()
	for range run {
		cnt[fastrand.Intn(n)] += 1
	}
	lps = time.Since(now)
	nmz = make([]float32, 0)
	for _, c := range cnt {
		nmz = append(nmz, float32(c)/float32(run))
	}
	return
}

func simIface(id, n, run int) (lps time.Duration, cnt []int, nmz []float32) {
	o := uint32(n)
	var rnd rand.Rand = rand.New(id)
	cnt = make([]int, 0)
	for range n {
		cnt = append(cnt, 0)
	}
	now := time.Now()
	for range run {
		cnt[rnd.Uint32n(o)] += 1
	}
	lps = time.Since(now)
	nmz = make([]float32, 0)
	for _, c := range cnt {
		nmz = append(nmz, float32(c)/float32(run))
	}
	return
}

func simStru(id, n, run int) (lps time.Duration, cnt []int, nmz []float32) {
	o := uint32(n)
	rnd := rand.Temp(id)
	cnt = make([]int, 0)
	for range n {
		cnt = append(cnt, 0)
	}
	now := time.Now()
	for range run {
		cnt[rnd.Uint32n(o)] += 1
	}
	lps = time.Since(now)
	nmz = make([]float32, 0)
	for _, c := range cnt {
		nmz = append(nmz, float32(c)/float32(run))
	}
	return
}

func random(typ, run, rng int) {
	var lps time.Duration
	var cnt []int
	var nmz []float32

	switch typ {
	case 0:
		fmt.Println("sea9go rand control")
		lps, cnt, nmz = simCtrl(rng, run)
	case 1:
		fmt.Println("sea9go rand iface")
		lps, cnt, nmz = simIface(1, rng, run)
	case 2:
		fmt.Println("sea9go rand stru")
		lps, cnt, nmz = simStru(1, rng, run)
	case 3:
		fmt.Println("sea9go rand direct")
		lps, cnt, nmz = simDirect(rng, run)
	}

	var buf strings.Builder
	for i, v := range cnt {
		fmt.Fprintf(&buf, " %3v: %v (%.4f%%)\n", i, v, nmz[i]*100)
	}
	fmt.Printf(" %v simulations with [0,%v) range, elapsed time: %12v (%.4fns per op)\n%v", run, rng, lps, float64(lps)/float64(run), buf.String())
}

func main() {
	// go build -pgo=pgo/cpu.pprof
	////////////////// pprof /////////////////////
	dir, _ := filepath.Split(os.Args[0])
	fcpu, e := os.Create(filepath.Join(dir, "pgo", "cpu.pprof.pend"))
	if e != nil {
		if !os.IsNotExist(e) {
			log.Fatal("Failed to create CPU profile", e)
		} // proceed if failed to create CPU profile
	}
	defer fcpu.Close()
	if e == nil {
		if e = pprof.StartCPUProfile(fcpu); e != nil {
			log.Fatal("Failed to start CPU profiling", e)
		}
		defer pprof.StopCPUProfile()
	}
	////////////////// pprof ///////////////////*/

	var err error
	run := 10000000000 // 10,000,000,000
	rng := 6
	switch len(os.Args) {
	case 4:
		run, err = strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatal(err)
		}
		fallthrough
	case 3:
		rng, err = strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		fallthrough
	case 2:
		switch os.Args[1] {
		case "all":
			if len(os.Args) < 4 {
				run = 1000000000 // 1,000,000,000
			}
			idcs := fastrand.Perm(3)
			random(0, run, rng)
			for _, idx := range idcs {
				random(idx+1, run, rng)
			}
		case "iface":
			random(1, run, rng)
		case "stru":
			random(2, run, rng)
		case "direct":
			random(3, run, rng)
		default:
			random(0, run, rng)
		}
	default:
		log.Println("Usage: ./rand [ctrl|iface|stru|direct|all] [range] [num-of-runs]")
	}

	////////////////// pprof /////////////////////
	fmem, e := os.Create(filepath.Join(dir, "pgo", "mem.pprof.pend"))
	if e != nil {
		if !os.IsNotExist(e) {
			log.Fatal("Failed to create Memory profile", e)
		} // do not panic if failed to create CPU profile
	}
	defer fmem.Close()
	if e == nil {
		runtime.GC() // get up-to-date statistics
		if e = pprof.WriteHeapProfile(fmem); e != nil {
			log.Fatal("Failed to write Memory profile", e)
		}
	}
	////////////////// pprof ///////////////////*/
}
