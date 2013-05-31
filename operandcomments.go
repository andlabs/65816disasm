// 31 may 2013
package main

import (
	"fmt"
)

func addoperandcomment(pos uint32, logical uint32, addr string) {
	physcal, inROM := memmap.Physical(logical)
	if inROM {
		addcomment(pos, "%s $%06X -> ROM $%06X",
			addr, logical, physical)
		return
	}
	addcomment(pos, "%s $%06X",
		addr, logical)
}

// hhll and the like
func addDBRComment(pos uint32, addr uint16) {
	if !env.dbr.known {
		addcomment(pos, "$%04X - cannot get physical address because dbr is unknown at time of disassembly", addr)
		return
	}
	logical := (uint32(env.dbr.value) << 16) | uint32(addr)
	addoperandcomment(pos, logical, fmt.Sprintf("$%04X ->", addr))
}

// hhllmm and the like
func addLongComment(pos uint32, logical uint32) {
	addoperandcomment(pos, logical, "")
}

// bank numbers in transfer instructions
func addBankComment(pos uint32, bank byte, what string) {
	addcomment(pos, "%s: %s",
		what, memmap.BankComment(bank))
}

// nn, (nn), and the like
func addDirectComment(pos uint32, addr byte) {
	if !env.direct.known {
		addcomment(pos, "$%02X - cannot get physical address because d is unknown at time of disassembly", addr)
		return
	}
	logical := (uint32(env.direct.value) + uint32(addr)) & 0xFFFF		// keep in bank 0
	addoperandcomment(pos, logical,
		fmt.Sprintf("$%02X + d=$%02X ->", addr, env.direct.value))
}

// [nn] and the like
func addDirectLongComment(pos uint32, addr byte) {
	if !env.direct.known {
		addcomment(pos, "$%02X - cannot get physical address because d is unknown at time of disassembly", addr)
		return
	}
	logical := (uint32(env.direct.value) + uint32(addr)) & 0xFFFF		// keep in bank
	if !env.dbr.known {
		addcomment(pos, "$%02X - cannot get physical address because dbr is unknown at time of disassembly", addr)
		return
	}
	logical |= uint32(env.dbr.value) << 16
	addoperandcomment(pos, logical,
		fmt.Sprintf("$%02X + d=$%02X + dbr=$%02X ->",
			addr, env.direct.value, env.dbr.value))
}

// used for (nn,s),y to inform what the dbr is for the indirection
// TODO add a proper addStackComment and addIndirectStackComment instead
func addDBRReminderComment(pos uint32) {
	if !env.dbr.known {
		addcomment(pos, "$%04X - cannot get physical address because dbr is unknown at time of disassembly", addr)
		return
	}
	addcomment(pos, "dbr=$%02X", env.dbr.value)
}