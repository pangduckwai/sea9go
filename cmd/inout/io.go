package main

import (
	"log"
	"os"

	"github.com/pangduckwai/sea9go/pkg/inout"
)

func inOut(in, out string) (err error) {
	inp := os.Stdin
	if in != "" {
		inp, err = os.Open(in)
		if err != nil {
			return
		}
		defer inp.Close()
	}

	opn := os.Stdout
	if out != "" {
		opn, err = os.Create(out)
		if err != nil {
			return
		}
		defer opn.Close()
	}

	read, err := inout.Read(inp, 16)
	if err != nil {
		return
	}

	err = inout.Write(opn, read)
	return
}

func main() {
	var err error
	var in, out string
	switch len(os.Args) {
	case 3:
		out = os.Args[2]
		fallthrough
	case 2:
		in = os.Args[1]
	default:
		log.Println("Usage: ./io [in file path] [out file path]")
		os.Exit(0)
	}
	err = inOut(in, out)
	if err != nil {
		log.Fatal(err)
	}
}
