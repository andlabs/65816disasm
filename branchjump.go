// 31 may 2013
// based on huc6280disasm
package main

import (
	"fmt"
	"os"
)

func dobranch(pos uint32) (labelpos uint32, newpos uint32) {
	origpos := pos - 1
	b, pos := getbyte(pos)
	offset := int32(int8(b))
	bpos := uint32(int32(pos) + offset)
	
	banksize := memmap.BankSize()
	if (origpos & banksize) != (bpos & banksize)  {
		fmt.Fprintf(os.Stderr, "cannot follow branch at $%X as it crosses a bank boundary\n", pos)
		mklabel(origpos, "invalid", lpUser)
		return origpos, pos
	}

	mklabel(bpos, "loc", lpLoc)
	if bpos != origpos {		// avoid endless recursion on branch to self
		// should this depend on -isolatesubs?
		var saved *envt = saveenv()
		disassemble(bpos)
		restoreenv(saved)
	}
	return bpos, pos
}

func dolongbranch(pos uint32) (labelpos uint32, newpos uint32) {
	origpos := pos - 1
	w, pos := getword(pos)
	offset := int32(int16(w))
	bpos := uint32(int32(pos) + offset)

	banksize := memmap.BankSize()
	if (origpos & banksize) != (bpos & banksize)  {
		fmt.Fprintf(os.Stderr, "cannot follow branch at $%X as it crosses a bank boundary\n", pos)
		mklabel(origpos, "invalid", lpUser)
		return origpos, pos
	}

	mklabel(bpos, "loc", lpLoc)
	if bpos != origpos {		// avoid endless recursion on branch to self
		// should this depend on -isolatesubs?
		var saved *envt = saveenv()
		disassemble(bpos)
		restoreenv(saved)
	}
	return bpos, pos
}

// xxx label
func op_branch(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		labelpos, pos := dobranch(pos)
		labelplaces[pos - 2] = labelpos
		// don't disassemble past unconditional branches
		// (since it'll either try to disassemble non-code or code that would ideally
		// still get disassembled anyway)
		done = m == "bra"
		return fmt.Sprintf("%s\t%%s", m), pos, done
	}
}

// brl label
func brl_pcrelativeword(pos uint32) (disassembled string, newpos uint32, done bool) {
	labelpos, pos := dolongbranch(pos)
	labelplaces[pos - 3] = labelpos
	return fmt.Sprintf("brl\t%%s"), pos, true // don't disassemble past unconditional branches
}

// jmp hhll
// page 113 of the programming manual says this is within the current program bank
// TODO this (all of these, actually, including jsr_absolute) does (do) not account for crossing banks
func jmp_absolute(pos uint32) (disassembled string, newpos uint32, done bool) {
	w, pos := getword(pos)
	logical := (uint32(env.pbr) << 16) | uint32(w)
	phys, inROM := memmap.Physical(logical)
	if !inROM {
		addPBRComment(pos - 3, pos, w)
		return fmt.Sprintf("jmp\t$%04X", w), pos, true
	}
	mklabel(phys, "loc", lpLoc)
	labelplaces[pos - 3] = phys
	if phys != (pos - 3) {		// avoid endless recursion on jump to self
		disassemble(phys)
	}
	return fmt.Sprintf("jmp\t%%s"), pos, true
}

// jmp (hhll)
func jmp_absoluteindirect(pos uint32) (disassembled string, newpos uint32, done bool) {
	w, pos := getword(pos)
	addPBRComment(pos - 3, pos, w)
	savedenvs[pos - 3] = saveenv()
	return fmt.Sprintf("jmp\t($%04X)", w), pos, true
}

// jmp (hhll,x)
func jmp_absoluteindirectx(pos uint32) (disassembled string, newpos uint32, done bool) {
	w, pos := getword(pos)
	addPBRComment(pos - 3, pos, w)
	savedenvs[pos - 3] = saveenv()
	return fmt.Sprintf("jmp\t($%04X,x)", w), pos, true
}

// jmp hhllmm
func jmp_absolutelong(pos uint32) (disassembled string, newpos uint32, done bool) {
	logical, pos := getlong(pos)
	phys, inROM := memmap.Physical(logical)
	if !inROM {
		addLongComment(pos - 4, logical)
		return fmt.Sprintf("jmp\t$%06X", logical), pos, true
	}
	mklabel(phys, "loc", lpLoc)
	labelplaces[pos - 4] = phys
	if phys != (pos - 4) {		// avoid endless recursion on jump to self
		oldpbr := env.pbr
		env.pbr = byte((logical >> 16) & 0xFF)
		disassemble(phys)
		env.pbr = oldpbr
	}
	return fmt.Sprintf("jmp\t%%s"), pos, true
}

// jmp [hhllmm]
func jmp_absolutelongindirect(pos uint32) (disassembled string, newpos uint32, done bool) {
	l, pos := getlong(pos)
	addLongComment(pos - 3, l)
	savedenvs[pos - 3] = saveenv()
	return fmt.Sprintf("jmp\t($%04X)", l), pos, true
}

// jsr hhll
func jsr_absolute(pos uint32) (disassembled string, newpos uint32, done bool) {
	w, pos := getword(pos)
	logical := (uint32(env.pbr) << 16) | uint32(w)
	phys, inROM := memmap.Physical(logical)
	if !inROM {
		addPBRComment(pos - 3, pos, w)
		return fmt.Sprintf("jmp\t$%04X", w), pos, true
	}
	mklabel(phys, "sub", lpSub)
	labelplaces[pos - 3] = phys
	if phys != (pos - 3) {		// avoid endless recursion on call to self
		var saved *envt

		if *isolateSubs {
			saved = saveenv()
		}
		disassemble(phys)
		if *isolateSubs {
			restoreenv(saved)
		}
	}
	if f, ok := specialsubs[phys]; ok {		// handle special subroutine action
		pos = f(pos)
	}
	return fmt.Sprintf("jsr\t%%s"), pos, false
}

// jsr (hhll,x)
func jsr_absoluteindirectx(pos uint32) (disassembled string, newpos uint32, done bool) {
	w, pos := getword(pos)
	addPBRComment(pos - 3, pos, w)
	savedenvs[pos - 3] = saveenv()
	return fmt.Sprintf("jsr\t($%04X,x)", w), pos, false
}

// jsr hhllmm
func jsr_absolutelong(pos uint32) (disassembled string, newpos uint32, done bool) {
	logical, pos := getlong(pos)
	phys, inROM := memmap.Physical(logical)
	if !inROM {
		addLongComment(pos - 4, logical)
		return fmt.Sprintf("jmp\t$%06X", logical), pos, true
	}
	mklabel(phys, "sub", lpSub)
	labelplaces[pos - 4] = phys
	if phys != (pos - 4) {		// avoid endless recursion on call to self
		var saved *envt

		if *isolateSubs {		// isolate before loading the new pbr
			saved = saveenv()
		}
		oldpbr := env.pbr
		env.pbr = byte((logical >> 16) & 0xFF)
		disassemble(phys)
		env.pbr = oldpbr
		if *isolateSubs {
			restoreenv(saved)
		}
	}
	if f, ok := specialsubs[phys]; ok {		// handle special subroutine action
		pos = f(pos)
	}
	return fmt.Sprintf("jsr\t%%s"), pos, false
}

// rti
// TODO for all of these: touch the stack?
func rti_noarguments(pos uint32) (disassembled string, newpos uint32, done bool) {
	return "rti", pos, true
}

// rtl
func rtl_noarguments(pos uint32) (disassembled string, newpos uint32, done bool) {
	return "rtl", pos, true
}

// rts
func rts_noarguments(pos uint32) (disassembled string, newpos uint32, done bool) {
	return "rts", pos, true
}
