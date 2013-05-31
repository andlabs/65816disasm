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
	0x00FFFE:	"EmuIRQBRK",
	0x00FFFC:	"EmuRESET",
	0x00FFFA:	"EmuNMI",
	0x00FFF8:	"EmuABORT",
	0x00FFF6:	"EmuReserved1",
	0x00FFF4:	"EmuCOP",
	0x00FFF2:	"EmuReserved2",
	0x00FFF0:	"EmuReserved3",
	0x00FFEE:	"NativeIRQ",
	0x00FFEC:	"NativeReserved1",
	0x00FFEA:	"NativeNMI",
	0x00FFE8:	"NativeABORT",
	0x00FFE6:	"NativeBRK",
	0x00FFE4:	"NativeCOP",
	0x00FFE2:	"NativeReserved2",
	0x00FFE0:	"NativeReserved3",
}

func errorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s ROM mode\n", os.Args[0])
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
		pos, inROM := memmap.Physical(uint32(posw))		// always bank 0
		if !inROM {
			fmt.Fprintf(os.Stderr, "physical address for %s vector ($%06X) not in ROM\n", label, uint32(posw))
			continue
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
