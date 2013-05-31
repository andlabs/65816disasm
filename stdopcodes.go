// 25 may 2013
// based on huc6280disasm
package main

import (
	"fmt"
)

// xxx hhll
func op_absolute(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		makeAUnknown()
		w, pos := getword(pos)
		addDBRComment(pos - 3, w)
		return fmt.Sprintf("%s\t$%04X", m, w), pos, false
	}
}

// xxx (hhll,x) - jump instructions only

// xxx hhll,x
func op_absolutex(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		makeAUnknown()
		w, pos := getword(pos)
		addDBRComment(pos - 3, w)
		return fmt.Sprintf("%s\t$%04X,x", m, w), pos, false
	}
}

// xxx hhll,y
func op_absolutey(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		makeAUnknown()
		w, pos := getword(pos)
		addDBRComment(pos - 3, w)
		return fmt.Sprintf("%s\t$%04X,y", m, w), pos, false
	}
}

// xxx (hhll) - jump instructions only

// xxx hhllmm,x
func op_absolutelongx(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		makeAUnknown()
		l, pos := getlong(pos)
		addLongComment(pos - 4, l)
		return fmt.Sprintf("%s\t$%06X,x", m, l), pos, false
	}
}

// xxx hhllmm
func op_absolutelong(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		makeAUnknown()
		l, pos := getlong(pos)
		addLongComment(pos - 4, l)
		return fmt.Sprintf("%s\t$%06X", m, l), pos, false
	}
}

// xxx a
func op_accumulator(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		makeAUnknown()
		return fmt.Sprintf("%s\ta", m), pos, false
	}
}

// xxx bank,bank
func op_transfer(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		makeAUnknown()
		dest, pos := getbyte(pos)				// binary dest,src; assembler src,dest
		src, pos := getbyte(pos)
		addBankComment(pos - 3, src, "src")
		addBankComment(pos - 3, dest, "dest")
		return fmt.Sprintf("%s\t#$%02X,#$%02X", m, src, dest), pos, false
	}
}

// xxx (nn,x)
func op_indirectx(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		makeAUnknown()
		b, pos := getbyte(pos)
		addDirectComment(pos - 2, b)
		return fmt.Sprintf("%s\t($%02X,x)", m, b), pos, false
	}
}

// xxx nn,x
func op_directx(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		makeAUnknown()
		b, pos := getbyte(pos)
		addDirectComment(pos - 2, b)
		return fmt.Sprintf("%s\t$%02X,x", m, b), pos, false
	}
}

// xxx nn,y
func op_directy(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		makeAUnknown()
		b, pos := getbyte(pos)
		addDirectComment(pos - 2, b)
		return fmt.Sprintf("%s\t$%02X,y", m, b), pos, false
	}
}

// xxx (nn),y
func op_indirecty(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		makeAUnknown()
		b, pos := getbyte(pos)
		addDirectComment(pos - 2, b)
		return fmt.Sprintf("%s\t($%02X),y", m, b), pos, false
	}
}

// xxx [nn],y
func op_indirectlongy(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		makeAUnknown()
		b, pos := getbyte(pos)
		addDirectLongComment(pos - 2, b)
		return fmt.Sprintf("%s\t[$%02X],y", m, b), pos, false
	}
}

// xxx [nn]
func op_indirectlong(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		makeAUnknown()
		b, pos := getbyte(pos)
		addDirectLongComment(pos - 2, b)
		return fmt.Sprintf("%s\t[$%02X]", m, b), pos, false
	}
}

// xxx (nn)
func op_indirect(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		makeAUnknown()
		b, pos := getbyte(pos)
		addDirectComment(pos - 2, b)
		return fmt.Sprintf("%s\t($%02X)", m, b), pos, false
	}
}

// xxx nn
func op_direct(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		makeAUnknown()
		b, pos := getbyte(pos)
		addDirectLongComment(pos - 2, b)		// base bank is always zero
		return fmt.Sprintf("%s\t$%02X", m, b), pos, false
	}
}

// xxx #nn
func op_immediate(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		makeAUnknown()
		b, pos := getbyte(pos)
		return fmt.Sprintf("%s\t#$%02X", m, b), pos, false
	}
}

// xxx
func op_noarguments(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		makeAUnknown()
		return fmt.Sprintf("%s", m), pos, false
	}
}

// xxx nn,s
func op_stack(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		makeAUnknown()
		b, pos := getbyte(pos)
		return fmt.Sprintf("%s\t$%02X,s", m, b), pos, false
	}
}

// xxx (nn,s),y
func op_indirectstack(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		makeAUnknown()
		b, pos := getbyte(pos)
		addDBRReminderComment(pos - 2)
		return fmt.Sprintf("%s\t($%02X,s),y", m, b), pos, false
	}
}
