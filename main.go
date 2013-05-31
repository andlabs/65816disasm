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

var vectorLocs = map[uint32]string{
	0x1FFE:	"EntryPoint",
	0x1FFC:	"NMI",
	0x1FFA:	"TimerInterrupt",
	0x1FF8:	"IRQ1",
	0x1FF6:	"IRQ2_BRK",
}

func errorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

// command-line options
var (
	useStack = flag.Bool("stack", false, "follow stack for tam/tma values (may fix some broken disassemblies but breaks if some subroutine breaks the push/pop system (TODO add jsr and rts))")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [-stack] ROM\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	var err error

	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 1 {
		usage()
	}

	filename := flag.Arg(0)

	bytes, err = ioutil.ReadFile(filename)
	if err != nil {
		errorf("error reading input file %s: %v", filename, err)
	}
	if len(bytes) < 0x2000 {
		errorf("given input file %s does not provide a complete interrupt vector table (this restriction may be lifted in the future)", filename)
	}
	if len(bytes) >= 0x1F0000 {
		errorf("given input file %s too large (this restriction may be lifted in the future)", filename)
	}

	instructions = map[uint32]string{}
	labels = map[uint32]string{}
	labelpriorities = map[uint32]int{}
	labelplaces = map[uint32]uint32{}
	comments = map[uint32]string{}

	// autoanalyze vectors
	for addr, label := range vectorLocs {
		posw, _ := getword(addr)
		pos, err := physical(posw)
		if err != nil {
			errorf("internal error: could not get physical address for %s vector (meaning something is up with the paging or the game actually does have the vector outside page 7): %v\n", label, err)
		}
		if labels[pos] != "" {		// if already defined as a different vector, concatenate the labels to make sure everything is represented
			// TODO because this uses a map, it will not be in vector order
			labels[pos] = labels[pos] + "_" + label
		} else {
			labels[pos] = label
		}
		labelpriorities[pos] = lpSub
		disassemble(pos)
	}

	// TODO read additional starts from standard input

	print()
}
