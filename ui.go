// 1 june 2013
package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
)

// TODO better name than ui?

var commands = []struct {
	Name	string
	Desc		string
	Func		func(fields []string)
}{
	{ "help", "show this help", c_help },
	{ "doauto", "auto-analyze vectors", c_doauto },		// TODO keep "vectors"?
}

// the map key is a logical address
// TODO this will need to be made part of MemoryMap later if it branches out from SNES ROMs
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

func c_doauto(fields []string) {
	for addr, label := range vectorLocs {
		pos, inROM := memmap.Physical(addr)				// addr is a logical address
		if !inROM {
			errorf("sanity check failure: vector %s ($%06X) not in ROM (memmap.Physical() returned $%06X)", label, addr, pos)
		}
		posw, _ := getword(pos)
		pos, inROM = memmap.Physical(uint32(posw))		// always bank 0
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
		env.pbr = 0			// we execute from bank 0
		disassemble(pos)
	}
	fmt.Fprintf(os.Stderr, "finished auto-analyzing vectors\n")
}

var helptext string

func c_help(fields []string) {
	fmt.Fprintf(os.Stderr, "%s", helptext)
}

func init() {
	for _, v := range commands {
		helptext += fmt.Sprintf("%10s - %s\n", v.Name, v.Desc)
	}
}

func doui() {
	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		line := stdin.Text()
		fields := strings.Fields(line)
		// TODO this means comments cannot start in the middle of a token
		lastValid := 0					// strip comments
		for ; lastValid < len(fields); lastValid++ {
			if fields[lastValid][0] == '#' {
				break
			}
		}
		fields = fields[:lastValid]
		if len(fields) == 0 {				// blank line
			continue
		}
		command := fields[0]
		found := false
		for _, v := range commands {
			if command == v.Name {
				v.Func(fields[1:])
				found = true
				break
			}
		}
		if !found {
			fmt.Fprintf(os.Stderr, "command not found\n")
		}
	}
	if err := stdin.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading standard input: %v\n", err)
		os.Exit(1)
	}
}
