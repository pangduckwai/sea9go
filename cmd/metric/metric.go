package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strconv"
	"time"

	"github.com/bytedance/gopkg/lang/fastrand"
	"github.com/pangduckwai/sea9go/pkg/logger/metric"
)

var hILO_DEC = [][]uint64{
	{1844674407370955161, 11068046444225730970, 10},
	{184467440737095516, 2951479051793528259, 100},
	{18446744073709551, 11363194349405083796, 1000}, // k 10^3
	{1844674407370955, 2980993842311463542, 10000},
	{184467440737095, 9521471421085922163, 100000},
	{18446744073709, 10175519178963368025, 1000000}, // M 10^6
	{1844674407370, 17619621584234933257, 10000000},
	{184467440737, 1761962158423493326, 100000000},
	{18446744073, 13088917067439035464, 1000000000}, // G 10^9
	{1844674407, 6842914928856769032, 10000000000},
	{184467440, 13597012344482363035, 100000000000},
	{18446744, 1359701234448236304, 1000000000000}, // T 10^12
	{1844674, 7514667752928644277, 10000000000000},
	{184467, 8130164404776685075, 100000000000000},
	{18446, 13725737292074354639, 1000000000000000}, // P 10^15
	{1844, 12440620173433166434, 10000000000000000},
	{184, 8622759646827137290, 100000000000000000},
	{18, 8240973594166534376, 1000000000000000000}, // E 10^18
	{1, 15581492618384294731, 10000000000000000000},
}

// sUFFIX list of metric suffices of values can be fit into an int64
var sUFFIX = []struct {
	s string
	i int
}{
	{"k", 2},  // kilo
	{"M", 5},  // mega
	{"G", 8},  // giga
	{"T", 11}, // tera
	{"P", 14}, // peta
	{"E", 17}, // exa
}

var tESTS = []int64{
	math.MinInt64 + 2,
	-56789012,
	-5678901,
	-567890,
	-56789,
	-5678,
	-567,
	-56,
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

func metricCtrl(inp int64, dec int) string {
	neg := ""
	if inp < 0 {
		neg = "-"
		inp = -inp
	}

	if inp < int64(hILO_DEC[sUFFIX[0].i][2]) {
		return fmt.Sprintf("%v%v", neg, inp)
	}
	var rnd uint64 = 1
	for j := 0; j < dec; j++ {
		rnd *= 10
	}
	round := float64(rnd)
	i := 1
	for ; i < len(sUFFIX); i++ {
		if (uint64(inp) / hILO_DEC[sUFFIX[i].i][2]) <= 0 {
			return fmt.Sprintf("%v%v %v", neg, math.Round((float64(inp)/float64(hILO_DEC[sUFFIX[i-1].i][2]))*round)/round, sUFFIX[i-1].s)
		}
	}
	return fmt.Sprintf("%v%v %v", neg, math.Round((float64(inp)/float64(hILO_DEC[sUFFIX[i-1].i][2]))*round)/round, sUFFIX[i-1].s)
}

func metricSimu(run int) {
	fmt.Println("sea9go metric")
	var lps time.Duration
	var n uint32 = uint32(len(tESTS))
	now := time.Now()
	for range run {
		// idx := fastrand.Uint32n(n)
		// dec := int(fastrand.Uint32n(7)) + 3
		// tmp := metric.Metric(tESTS[idx], dec)
		// log.Printf(" Metric(%20v, %v) -> %16v\n", tESTS[idx], dec, tmp)
		metric.Metric(tESTS[fastrand.Uint32n(n)], int(fastrand.Uint32n(7))+3) // decmial point range 3 to 9
	}
	lps = time.Since(now)
	fmt.Printf(" %v simulations elapsed time: %12v (%.4fns per op)\n", run, lps, float64(lps)/float64(run))
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
	run := 100000000 // 100,000,000
	switch len(os.Args) {
	case 3:
		run, err = strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		fallthrough
	case 2:
		switch os.Args[1] {
		case "all":
			log.Println("WIP: all")
		case "metric":
			metricSimu(run)
		default:
			log.Println("WIP: control")
		}
	default:
		log.Println("Usage: ./metric [ctrl|metric|all] [num-of-runs]")
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
