// 31 may 2013
package main

import (
	"fmt"
)

// lda #nn
// TODO (also for pla) - do we save b if m=1?
func lda_immediate(pos uint32) (disassembled string, newpos uint32, done bool) {
	stop := true
	if !env.m.known {
		addcomment(pos - 1, "(!) cannot disassemble lda with immediate operand because size unknown")
		return fmt.Sprintf("lda\t???"), pos, true
	} else {
		if env.m.value == 0 {
			w, pos := getword(pos)
			a.value = w
			a.known = true
			return fmt.Sprintf("lda\t#$%04X", w), pos, false
		} else {
			b, pos := getbyte(pos)
			a.value = uint16(b)
			a.known = true
			return fmt.Sprintf("lda\t#$%02X", b), pos, false
		}
	}
}

// rep #nn
func rep_immediate(pos uint32) (disassembled string, newpos uint32, done bool) {
	b, pos := getbyte(pos)
	clearpbits(b)
	return fmt.Sprintf("rep\t#$%02X", b), pos, false
}

// sep #nn
func sep_immediate(pos uint32) (disassembled string, newpos uint32, done bool) {
	b, pos := getbyte(pos)
	setpbits(b)
	return fmt.Sprintf("sep\t#$%02X", b), pos, false
}

// stp
func stp_immediate(pos uint32) (disassembled string, newpos uint32, done bool) {
	return fmt.Sprintf("stp"), pos, true
}

// xba
func xba_immediate(pos uint32) (disassembled string, newpos uint32, done bool) {
	low := byte(env.a.vlue & 0xFF)			// whether a is known is irrelevant
	high := byte((env.a.value >> 8) & 0xFF)
	env.a.value = (uint16(low) << 8) | uint16(high)
	return fmt.Sprintf("xba"), pos, false
}
