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
			fmt.Printf("%s:\n", label)
		}
		if instruction, ok := instructions[i]; ok && instruction != operandString {
			if labelpos, ok := labelplaces[i]; ok {		// need to add a label
				if labels[labelpos] == "" {
					labels[labelpos] = fmt.Sprintf("<no label for $%X>", labelpos)
				}
				instruction = fmt.Sprintf(instruction, labels[labelpos])
			}
			fmt.Printf("\t%s\t\t; $%X", instruction, i)
			if comment, ok := comments[i]; ok {
				fmt.Printf(" | %s", comment)
			}
			fmt.Println()
		}
	}
}
