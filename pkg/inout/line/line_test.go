package line

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

const BUFFER = 64 //1048576

func TestLine(t *testing.T) {
	f, err := os.Open("../../../README.md")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	lf, pg := fmt.Sprintf("%c", 10), fmt.Sprintf("%c", 185)
	read := make([]string, 0)
	ctrl := make([]string, 0)
	line := make([]byte, 0)

	rcnt, lcnt, err := ReadLine(
		f, BUFFER,
		func(s string) error {
			read = append(read, s)
			return nil
		},
		func(n, i0, i1 int, b []byte) {
			dspy := strings.ReplaceAll(string(b), lf, pg)
			mask := make([]byte, n)
			for i := range mask {
				if i >= i0 && i < i0+i1 {
					mask[i] = '\''
				} else {
					mask[i] = ' '
				}
			}
			if i1 < n {
				fmt.Printf("TestLine()+ [%3v:%3v] \"%v\" -> \"%s\"\n", i0, i1, dspy, b[i0:i0+i1])
				line = append(line, b[i0:i0+i1]...)
				ctrl = append(ctrl, string(line))
				line = line[:0]
			} else {
				fmt.Printf("TestLine()¤ [%3v:%3v] \"%v\" -> \"%s\"\n", i0, i1, dspy, b[i0:])
				line = append(line, b[i0:]...)
			}
			fmt.Printf("                       %s\n", mask)
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(line) > 0 {
		ctrl = append(ctrl, string(line))
	}

	if len(read) != len(ctrl) {
		t.Fatalf("TestLine() expected %v lines but got %v", len(ctrl), len(read))
	}
	if len(read) < 1 {
		t.Fatal("TestLine() no line read")
	}
	for i := range read {
		if read[i] != ctrl[i] {
			t.Fatalf("TestLine() line %v  expected '%s' but got '%s'", i, ctrl[i], read[i])
		}
	}
	fmt.Printf("TestLine() - Read %v x %v bytes, %v lines\n", rcnt, BUFFER, lcnt)
}
