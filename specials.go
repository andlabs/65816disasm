// 31 may 2013
package main

import (
	"fmt"
)

// brk #nn
// TODO invalidate a?
func brk_immediate(pos uint32) (disassembled string, newpos uint32, done bool) {
	b, pos := getbyte(pos)
	return fmt.Sprintf("brk\t#$%02X", b), pos, false
}

// clc
func clc_noarguments(pos uint32) (disassembled string, newpos uint32, done bool) {
	env.carryflag.value = 0
	env.carryflag.known = true
	return fmt.Sprintf("clc"), pos, false
}

// sec
func sec_noarguments(pos uint32) (disassembled string, newpos uint32, done bool) {
	env.carryflag.value = 1
	env.carryflag.known = true
	return fmt.Sprintf("sec"), pos, false
}

// lda #nn
// TODO (also for pla) - do we save b if m=1?
func lda_immediate(pos uint32) (disassembled string, newpos uint32, done bool) {
	if !env.m.known {
		addcomment(pos - 1, "(!) cannot disassemble lda with immediate operand because size unknown")
		return fmt.Sprintf("lda\t???"), pos, true
	} else {
		if env.m.value == 0 {
			w, pos := getword(pos)
			env.a.value = w
			env.a.known = true
			return fmt.Sprintf("lda\t#$%04X", w), pos, false
		} else {
			b, pos := getbyte(pos)
			env.a.value = uint16(b)
			env.a.known = true
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
func stp_noarguments(pos uint32) (disassembled string, newpos uint32, done bool) {
	return fmt.Sprintf("stp"), pos, true
}

// wdm #nn
func wdm_immediate(pos uint32) (disassembled string, newpos uint32, done bool) {
	b, pos := getbyte(pos)
	return fmt.Sprintf("wdm\t#$%02X", b), pos, false
}

// xba
func xba_noarguments(pos uint32) (disassembled string, newpos uint32, done bool) {
	low := byte(env.a.value & 0xFF)		// whether a is known is irrelevant
	high := byte((env.a.value >> 8) & 0xFF)
	env.a.value = (uint16(low) << 8) | uint16(high)
	return fmt.Sprintf("xba"), pos, false
}

// xce
func xce_noarguments(pos uint32) (disassembled string, newpos uint32, done bool) {
	stop := false
	if !env.carryflag.known {
		addcomment(pos - 1, "(!) cannot swap in emulation mode flag because carry flag is not known, meaning we cannot set the m and x flags properly")
		stop = true
	} else {
		env.carryflag, env.e = env.e, env.carryflag
		env.x.value = env.e.value
		env.x.known = true
		env.m.value = env.e.value
		env.m.known = true
	}
	return fmt.Sprintf("xce"), pos, stop
}
