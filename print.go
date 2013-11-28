// 31 may 2013
// based on huc6280disasm
package main

import (
	"fmt"
)

func print() {
	lbu32 := uint32(len(bytes))
	for i := uint32(0); i < lbu32; i++ {
		if label, ok := labels[i]; ok {
			logical, mirror := memmap.Logical(i)
			if logical != mirror {
				fmt.Printf("%s: ; $%06X / $%06X\n", label, logical, mirror)
			} else {
				fmt.Printf("%s: ; $%06X\n", label, logical)
			}
		}
		if instruction, ok := instructions[i]; ok && instruction != operandString {
			if labelpos, ok := labelplaces[i]; ok {		// need to add a label
				if labels[labelpos] == "" {
					labels[labelpos] = fmt.Sprintf("<no label for $%X>", labelpos)
				}
				instruction = fmt.Sprintf(instruction, labels[labelpos])
			}
			fmt.Printf("\t%-20s\t; $%X", instruction, i)
			if comment, ok := comments[i]; ok {
				fmt.Printf(" | %s", comment)
			}
			fmt.Println()
		} else if !ok && *showAll {
			has := false
			fmt.Printf("\tdc.b\t$%02X\t\t; $%X", bytes[i], i)
			if i + 1 < lbu32 {
				fmt.Printf(" | w")
				has = true
				if i + 2 < lbu32 {
					fmt.Printf("/l $(%02X)",
						bytes[i + 2])
				} else {
					fmt.Printf(" $")
				}
				fmt.Printf("%02X%02X",
					bytes[i + 1], bytes[i])
			}
			if bytes[i] >= 0x20 && bytes[i] < 0x7F {		// ASCII
				if !has {
					fmt.Printf(" |")
					has = true
				}
				fmt.Printf(" '%c'", bytes[i])
			}
			if comment, ok := comments[i]; ok {
				fmt.Printf(" | %s", comment)
			}
			fmt.Println()
		}
	}
}
