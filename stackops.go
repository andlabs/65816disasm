// 31 may 2013
package main

import (
	"fmt"
)

// pea hhll
func pea_absolute(pos uint32) (disassembled string, newpos uint32, done bool) {
	w, pos := getword(pos)
	pushword(w, true)
	addDBRComment(pos - 3, w)		// just in case
	return fmt.Sprintf("pea\t$%04X", w), pos, false
}

// pei (nn)
func pei_indirect(pos uint32) (disassembled string, newpos uint32, done bool) {
	b, pos := getbyte(pos)
	pushbyte(0, false)				// push dummy
	addDirectComment(pos - 2, b)
	return fmt.Sprintf("pei\t($%02X)", b), pos, false
}

// per hhll
func per_pcrelativeword(pos uint32) (disassembled string, newpos uint32, done bool) {
	w, pos := getword(pos)
	// TODO does this properly handle crossing banks?
	ealong := uint32(int32(uint16(w)) + int32(pos))
	pushword(uint16(ealong), true)	// just push the low word without the bank
	return fmt.Sprintf("pea\t$%04X", w), pos, false
}

// pha
func pha_noarguments(pos uint32) (disassembled string, newpos uint32, done bool) {
	stop := false
	if !env.m.known {
		addcomment(pos - 1, "(!) cannot push a because size is not known")
		stop = true
	} else {
		if env.m.value == 0 {
			pushword(env.a.value, env.a.known)
		} else {
			pushbyte(byte(env.a.value & 0xFF), env.a.known)
		}
	}
	return fmt.Sprintf("pha"), pos, stop
}

// phb
func phb_noarguments(pos uint32) (disassembled string, newpos uint32, done bool) {
	pushbyte(env.dbr.value, env.dbr.known)
	return fmt.Sprintf("phb"), pos, false
}

// phd
func phd_noarguments(pos uint32) (disassembled string, newpos uint32, done bool) {
	pushbyte(env.direct.value, env.direct.known)
	return fmt.Sprintf("phd"), pos, false
}

// phk
func phk_noarguments(pos uint32) (disassembled string, newpos uint32, done bool) {
	// TODO does this account for crossed bank?
	pushbyte(byte((pos & 0xFF0000) >> 16), true)
	return fmt.Sprintf("phk"), pos, false
}

// php
func php_noarguments(pos uint32) (disassembled string, newpos uint32, done bool) {
	p, err := getp()
	pushbyte(p, err == nil)		// if there is an error, then make p unknown
	return fmt.Sprintf("php"), pos, false
}

// phx
func phx_noarguments(pos uint32) (disassembled string, newpos uint32, done bool) {
	stop := false
	if !env.x.known {
		addcomment(pos - 1, "(!) cannot push x because size is not known")
		stop = true
	} else {
		if env.x.value == 0 {			// we don't track x so push dummies
			pushword(0, false)
		} else {
			pushbyte(0, false)
		}
	}
	return fmt.Sprintf("phx"), pos, stop
}

// phy
func phy_noarguments(pos uint32) (disassembled string, newpos uint32, done bool) {
	stop := false
	if !env.x.known {
		addcomment(pos - 1, "(!) cannot push y because size is not known")
		stop = true
	} else {
		if env.x.value == 0 {			// we don't track y so push dummies
			pushword(0, false)
		} else {
			pushbyte(0, false)
		}
	}
	return fmt.Sprintf("phy"), pos, stop
}

// pla
func pla_noarguments(pos uint32) (disassembled string, newpos uint32, done bool) {
	stop := false
	if !env.m.known {
		addcomment(pos - 1, "(!) cannot pop a because size is not known")
		stop = true
	} else {
		if env.m.value == 0 {
			env.a.value, env.a.known = popword()
		} else {
			v, k := popbyte()
			env.a.value = uint16(v)
			env.a.known = k
		}
	}
	return fmt.Sprintf("pla"), pos, stop
}

// plb
func plb_noarguments(pos uint32) (disassembled string, newpos uint32, done bool) {
	env.dbr.value, env.dbr.known = popbyte()
	return fmt.Sprintf("plb"), pos, false
}

// pld
func pld_noarguments(pos uint32) (disassembled string, newpos uint32, done bool) {
	env.direct.value, env.direct.known = popbyte()
	return fmt.Sprintf("pld"), pos, false
}

// plp
// TODO does this alter m and x in emulation mode?
func plp_noarguments(pos uint32) (disassembled string, newpos uint32, done bool) {
	setp(popbyte())
	return fmt.Sprintf("plp"), pos, false
}

// plx
func plx_noarguments(pos uint32) (disassembled string, newpos uint32, done bool) {
	stop := false
	if !env.x.known {
		addcomment(pos - 1, "(!) cannot pop x because size is not known")
		stop = true
	} else {
		if env.x.value == 0 {			// we don't track x so discared popped value
			popword()
		} else {
			popbyte()
		}
	}
	return fmt.Sprintf("plx"), pos, stop
}

// ply
func ply_noarguments(pos uint32) (disassembled string, newpos uint32, done bool) {
	stop := false
	if !env.x.known {
		addcomment(pos - 1, "(!) cannot pop y because size is not known")
		stop = true
	} else {
		if env.x.value == 0 {			// we don't track y so discared popped value
			popword()
		} else {
			popbyte()
		}
	}
	return fmt.Sprintf("ply"), pos, stop
}
