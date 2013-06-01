// 1 june 2013
package main

import (
	"fmt"
	"os"
	"strconv"
)

type specialsub func(pos uint32) uint32

var specialsubs map[uint32]specialsub

var specialsubslist = []struct {
	Name		string
	Format		string
	Help			string
	Generator		func(fields []string) specialsub
}{
	{ "stringafter", "stringafter n",
		"marks the function as having a null-terminated string immediately after the call, skilling n bytes first", ssgen_stringafter },
}

func ssgen_stringafter(fields []string) specialsub {
	if len(fields) != 1 {
		fmt.Fprintf(os.Stderr, "stringafter usage: specialsub addr stringafter n\n")
		return nil
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "stringafter error: invalid number %q: %v", fields[0], err)
		return nil
	}
	return func(pos uint32) uint32 {
		var b byte

		for i := 0; i < n; i++ {
			b, pos = getbyte(pos)
			instructions[pos - 1] = fmt.Sprintf("dc.b\t$%02X", b)
		}
		sb := make([]byte, 0, 128)
		origpos := pos
		for {
			b, pos = getbyte(pos)
			if b == 0 {
				break
			}
			sb = append(sb, b)
		}
		instructions[origpos] = fmt.Sprintf("dc.b\t%q, 0", sb)
		return pos
	}
}

// TODO help command

func c_specialsub(fields []string) {
	if len(fields) < 2 {
		fmt.Fprintf(os.Stderr, "specialsubs usage: specialsubs addr command [args]\n")
		return
	}

	// TODO addr must be a bare hex number with this
	addr64, err := strconv.ParseUint(fields[0], 16, 32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "specialsubs error: invalid address hex number %q: %v", fields[0], err)
		return
	}
	addr := uint32(addr64)

	command := fields[1]
	args := fields[2:]
	found := false
	for _, v := range specialsubslist {
		if command == v.Name {
			f := v.Generator(args)
			if f != nil {
				specialsubs[addr] = f
			}
			found = true
		}
	}
	if !found {
		fmt.Fprintf(os.Stderr, "specialsubs %s: command not found\n", command)
	}
}
