// 1 june 2013
package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"strconv"
)

// TODO better name than ui?

var commands = []struct {
	Name	string
	Desc		string
	Func		func(fields []string)
}{
	{ "help", "show this help", c_help },
	{ "doauto", "auto-analyze vectors", c_doauto },		// TODO keep "vectors"?
	{ "specialsub", "mark a subroutine as doing something special", c_specialsub },
	{ "dowordptr", "mark a word as a pointer to code (with a pre-existing environment)", c_dowordptr },
	{ "dowordmanual", "mark a word as code (with a pre-existing environment)", c_dowordmanual },
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

func c_dowordptr(fields []string) {
	if len(fields) != 2 {
		fmt.Fprintf(os.Stderr, "dowordptr usage: dowordptr word-address env-address\n")
		return
	}

	// TODO addr and envaddr must be bare hex numbers with this
	addr64, err := strconv.ParseUint(fields[0], 16, 32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dowordptr error: invalid address hex number %q: %v", fields[0], err)
		return
	}
	addr := uint32(addr64)

	envaddr64, err := strconv.ParseUint(fields[1], 16, 32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "specialsubs error: invalid environment address hex number %q (for $%06X): %v", fields[1], addr, err)
		return
	}
	envaddr := uint32(envaddr64)

	if addr + 1 > uint32(len(bytes)) {
		fmt.Fprintf(os.Stderr, "specialsubs error: address $%06X not in ROM\n", addr)
		return
	}
	npbase := (uint16(bytes[addr + 1]) << 8) | uint16(bytes[addr])

	if env, ok := savedenvs[envaddr]; ok {
		nplogical := (uint32(env.pbr) << 16) | uint32(npbase)
		npaddr, inROM := memmap.Physical(nplogical)
		if !inROM {
			fmt.Fprintf(os.Stderr, "dowordptr error: new address $%06X (from $%06X) not in ROM\n", nplogical, addr)
			return
		}
		mklabel(npaddr, "loc", lpLoc)
		restoreenv(env)
		disassemble(npaddr)
		labelplaces[addr] = npaddr
		instructions[addr] = "dc.w\t(%s & 0xFFFF)"
		instructions[addr + 1] = operandString
	} else {
		fmt.Fprintf(os.Stderr, "dowordptr error: no environment available for environment $%06X (from $%06X)\n", envaddr, addr)
		return
	}
}

func c_dowordmanual(fields []string) {
	if len(fields) != 2 {
		fmt.Fprintf(os.Stderr, "dowordmanual usage: dowordptr sub-address env-address\n")
		return
	}

	// TODO addr and envaddr must be bare hex numbers with this
	addr64, err := strconv.ParseUint(fields[0], 16, 32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dowordmanual error: invalid address hex number %q: %v", fields[0], err)
		return
	}
	addr := uint32(addr64)

	envaddr64, err := strconv.ParseUint(fields[1], 16, 32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "specialsubs error: invalid environment address hex number %q (for $%06X): %v", fields[1], addr, err)
		return
	}
	envaddr := uint32(envaddr64)

	if addr + 1 > uint32(len(bytes)) {
		fmt.Fprintf(os.Stderr, "specialsubs error: address $%06X not in ROM\n", addr)
		return
	}

	if env, ok := savedenvs[envaddr]; ok {
		cplogical := (uint32(env.pbr) << 16) | addr
		cpaddr, inROM := memmap.Physical(cplogical)
		if !inROM {
			fmt.Fprintf(os.Stderr, "dowordmanual error: new address $%06X ($%06X physical) not in ROM\n", cplogical, cpaddr)
			return
		}
		fmt.Fprintf(os.Stderr, "dowordmanual info: new address added to code: $%06X ($%06X physical)\n", cplogical, cpaddr)
		mklabel(cpaddr, "sub", lpSub)
		restoreenv(env)
		disassemble(cpaddr)
		labelplaces[addr] = cpaddr
	} else {
		fmt.Fprintf(os.Stderr, "dowordmanual error: no environment available for environment $%06X (from $%06X)\n", envaddr, addr)
		return
	}
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
