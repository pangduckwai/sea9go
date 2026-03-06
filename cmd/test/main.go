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

// // fastRand returns a uniform value in [0,n)
//
//	func fastRand(n int) int {
//		return int(fastrand.Uint32n(uint32(n)))
//	}
func simFast(n, run int) (lps time.Duration, cnt []int, nmz []float32) {
	n32 := uint32(n)
	cnt = make([]int, 0)
	for range n {
		cnt = append(cnt, 0)
	}
	now := time.Now()
	for range run {
		cnt[int(fastrand.Uint32n(n32))] += 1
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
	var rnd rand.Rand = rand.RandFast(id)
	cnt = make([]int, 0)
	for range n {
		cnt = append(cnt, 0)
	}
	now := time.Now()
	for range run {
		cnt[rnd.Intn(n)] += 1
	}
	lps = time.Since(now)
	nmz = make([]float32, 0)
	for _, c := range cnt {
		nmz = append(nmz, float32(c)/float32(run))
	}
	return
}

func simDirect(id, n, run int) (lps time.Duration, cnt []int, nmz []float32) {
	var rnd rand.RandFast = rand.RandFast(id)
	cnt = make([]int, 0)
	for range n {
		cnt = append(cnt, 0)
	}
	now := time.Now()
	for range run {
		cnt[rnd.Intn(n)] += 1
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
		fmt.Println("sea9go test rand control")
		lps, cnt, nmz = simCtrl(rng, run)
	case 1:
		fmt.Println("sea9go test rand iface")
		lps, cnt, nmz = simIface(1, rng, run)
	case 2:
		fmt.Println("sea9go test rand direct")
		lps, cnt, nmz = simDirect(1, rng, run)
	case 3:
		fmt.Println("sea9go test rand fast")
		lps, cnt, nmz = simFast(rng, run)
	}

	var buf strings.Builder
	for i, v := range cnt {
		fmt.Fprintf(&buf, " %3v: %v (%.4f%%)\n", i, v, nmz[i]*100)
	}
	fmt.Printf(" %v simulations with [0,%v) range, elapsed time: %12v (%v per op)\n%v", run, rng, lps, lps/time.Duration(run), buf.String())
}

func main() {
	// go build -pgo=pgo/cpu.pprof
	////////////////// pprof /////////////////////
	dir, _ := filepath.Split(os.Args[0])
	fcpu, e := os.Create(filepath.Join(dir, "pgo", "cpu.pprof.pend"))
	if e != nil {
		if !os.IsNotExist(e) {
			log.Fatal("Failed to create CPU profile", e)
			// } else {
			// 	log.Printf("Unable to create CPU profile: %v", e)
		}
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
	var cmd string
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "io":
		log.Println("Usage: ./test io ***WIP")
	case "rand":
		run := 1000000000 // 1,000,000,000
		rng := 6
		typ := 0
		switch len(os.Args) {
		case 5:
			run, err = strconv.Atoi(os.Args[4])
			if err != nil {
				log.Fatal(err)
			}
			fallthrough
		case 4:
			rng, err = strconv.Atoi(os.Args[3])
			if err != nil {
				log.Fatal(err)
			}
			fallthrough
		case 3:
			switch os.Args[2] {
			case "iface":
				typ = 1
			case "direct":
				typ = 2
			case "fast":
				typ = 3
			}
			fallthrough
		case 2:
			random(typ, run, rng)
		default:
			log.Println("Usage: ./test rand [ctrl|iface|direct|fast] [range] [num-of-runs]")
		}

	case "traverse":
		log.Println("Usage: ./test traverse ***WIP")
	default:
		log.Println("Usage: ./test [io|rand|traverse] {options...}")
	}

	////////////////// pprof /////////////////////
	fmem, e := os.Create(filepath.Join(dir, "pgo", "mem.pprof.pend"))
	if e != nil {
		if !os.IsNotExist(e) {
			log.Fatal("Failed to create Memory profile", e)
			// } else {
			// 	log.Printf("Unable to create Memory profile: %v", e)
		}
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
