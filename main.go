// 31 may 2013
// based on huc6280disasm
package main

import (
	"fmt"
	"os"
	"flag"
	"io/ioutil"
)

var bytes []byte
var instructions map[uint32]string
var labels map[uint32]string
var labelpriorities map[uint32]int
var labelplaces map[uint32]uint32
var comments map[uint32]string

const (
	lpLoc = iota
	lpLocret
	lpSub
	lpUser
)

func errorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

// command-line options
var (
	isolateSubs = flag.Bool("isolatesubs", false, "isolate subroutines in their own environment")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [options] ROM mode\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "mode must be one of:")
	for m := range memmaps {
		fmt.Fprintf(os.Stderr, " %s", m)
	}
	fmt.Fprintf(os.Stderr, "\n")
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	var err error

	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 2 {
		usage()
	}

	filename := flag.Arg(0)

	// TODO make case-insensitive?
	wantedmap := flag.Arg(1)
	if _, ok := memmaps[wantedmap]; !ok {
		fmt.Fprintf(os.Stderr, "unsupported memory map %q\n", wantedmap)
		usage()
	}
	memmap = memmaps[wantedmap]

	bytes, err = ioutil.ReadFile(filename)
	if err != nil {
		errorf("error reading input file %s: %v", filename, err)
	}
/*	if len(bytes) < 0x2000 {
		errorf("given input file %s does not provide a complete interrupt vector table (this restriction may be lifted in the future)", filename)
	}
	if len(bytes) >= 0x1F0000 {
		errorf("given input file %s too large (this restriction may be lifted in the future)", filename)
	}*/

	instructions = map[uint32]string{}
	labels = map[uint32]string{}
	labelpriorities = map[uint32]int{}
	labelplaces = map[uint32]uint32{}
	comments = map[uint32]string{}
	specialsubs = map[uint32]specialsub{}

	doui()

	print()
}
