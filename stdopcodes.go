// 25 may 2013
// based on huc6280disasm
package main

import (
	"fmt"
)

// xxx hhll
func op_absolute(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		switch m {
		case "adc", "and", "eor", "lda", "ora", "sbc":
			makeAUnknown()
		}
		w, pos := getword(pos)
		addDBRComment(pos - 3, w)
		return fmt.Sprintf("%s\t$%04X", m, w), pos, false
	}
}

// xxx (hhll,x) - jump instructions only

// xxx hhll,x
func op_absolutex(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		switch m {
		case "adc", "and", "eor", "lda", "ora", "sbc":
			makeAUnknown()
		}
		w, pos := getword(pos)
		addDBRComment(pos - 3, w)
		return fmt.Sprintf("%s\t$%04X,x", m, w), pos, false
	}
}

// xxx hhll,y
func op_absolutey(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		switch m {
		case "adc", "and", "eor", "lda", "ora", "sbc":
			makeAUnknown()
		}
		w, pos := getword(pos)
		addDBRComment(pos - 3, w)
		return fmt.Sprintf("%s\t$%04X,y", m, w), pos, false
	}
}

// xxx (hhll) - jump instructions only

// xxx hhllmm,x
func op_absolutelongx(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		switch m {
		case "adc", "and", "eor", "lda", "ora", "sbc":
			makeAUnknown()
		}
		l, pos := getlong(pos)
		addLongComment(pos - 4, l)
		return fmt.Sprintf("%s\t$%06X,x", m, l), pos, false
	}
}

// xxx hhllmm
func op_absolutelong(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		switch m {
		case "adc", "and", "eor", "lda", "ora", "sbc":
			makeAUnknown()
		}
		l, pos := getlong(pos)
		addLongComment(pos - 4, l)
		return fmt.Sprintf("%s\t$%06X", m, l), pos, false
	}
}

// xxx a
func op_accumulator(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		// TODO : actually update A if known
		makeAUnknown()
		return fmt.Sprintf("%s\ta", m), pos, false
	}
}

// xxx bank,bank
func op_transfer(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		// MVN and MVP always end with A = 0xFFFF
		env.a.known = true
		env.a.value = 0xFFFF
		
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
		switch m {
		case "adc", "and", "eor", "lda", "ora", "sbc":
			makeAUnknown()
		}
		b, pos := getbyte(pos)
		addDirectComment(pos - 2, b)
		return fmt.Sprintf("%s\t($%02X,x)", m, b), pos, false
	}
}

// xxx nn,x
func op_directx(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		switch m {
		case "adc", "and", "eor", "lda", "ora", "sbc":
			makeAUnknown()
		}
		b, pos := getbyte(pos)
		addDirectComment(pos - 2, b)
		return fmt.Sprintf("%s\t$%02X,x", m, b), pos, false
	}
}

// xxx nn,y
func op_directy(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		b, pos := getbyte(pos)
		addDirectComment(pos - 2, b)
		return fmt.Sprintf("%s\t$%02X,y", m, b), pos, false
	}
}

// xxx (nn),y
func op_indirecty(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		switch m {
		case "adc", "and", "eor", "lda", "ora", "sbc":
			makeAUnknown()
		}
		b, pos := getbyte(pos)
		addDirectComment(pos - 2, b)
		return fmt.Sprintf("%s\t($%02X),y", m, b), pos, false
	}
}

// xxx [nn],y
func op_indirectlongy(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		switch m {
		case "adc", "and", "eor", "lda", "ora", "sbc":
			makeAUnknown()
		}
		b, pos := getbyte(pos)
		addDirectComment(pos - 2, b)
		return fmt.Sprintf("%s\t[$%02X],y", m, b), pos, false
	}
}

// xxx [nn]
func op_indirectlong(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		switch m {
		case "adc", "and", "eor", "lda", "ora", "sbc":
			makeAUnknown()
		}
		b, pos := getbyte(pos)
		addDirectComment(pos - 2, b)
		return fmt.Sprintf("%s\t[$%02X]", m, b), pos, false
	}
}

// xxx (nn)
func op_indirect(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		switch m {
		case "adc", "and", "eor", "lda", "ora", "sbc":
			makeAUnknown()
		}
		b, pos := getbyte(pos)
		addDirectComment(pos - 2, b)
		return fmt.Sprintf("%s\t($%02X)", m, b), pos, false
	}
}

// xxx nn
func op_direct(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		switch m {
		case "adc", "and", "eor", "lda", "ora", "sbc":
			makeAUnknown()
		}
		b, pos := getbyte(pos)
		addDirectComment(pos - 2, b)		// base bank is always zero
		return fmt.Sprintf("%s\t$%02X", m, b), pos, false
	}
}

// xxx #nn
func op_immediate(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		// TODO : update the value of A if it is already known (or being loaded)
		switch m {
		case "adc", "and", "eor", "lda", "ora", "sbc":
			makeAUnknown()
		}
		
		if !env.e.known || !env.m.known {
			addcomment(pos - 1, "(!) cannot disassemble opcode with immediate operand because size unknown")
			return fmt.Sprintf("%s\t???", m), pos, true
		} else {
			if env.e.value | env.m.value == 0 {
				w, pos := getword(pos)
				return fmt.Sprintf("%s\t#$%04X", m, w), pos, false
			} else {
				b, pos := getbyte(pos)
				return fmt.Sprintf("%s\t#$%02X", m, b), pos, false
			}
		}
	}
}

// xxx
func op_noarguments(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		switch m {
		case "tda":
			env.a.value = env.direct.value
			env.a.known = env.direct.known
		case "tsa", "txa", "tya":
			makeAUnknown()
		}
		return fmt.Sprintf("%s", m), pos, false
	}
}

// xxx nn,s
func op_stack(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		switch m {
		// TODO : try to update A based on the stack contents if known
		case "adc", "and", "eor", "lda", "ora", "sbc":
			makeAUnknown()
		}
		b, pos := getbyte(pos)
		return fmt.Sprintf("%s\t$%02X,s", m, b), pos, false
	}
}

// xxx (nn,s),y
func op_indirectstack(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		switch m {
		case "adc", "and", "eor", "lda", "ora", "sbc":
			makeAUnknown()
		}
		b, pos := getbyte(pos)
		addDBRReminderComment(pos - 2)
		return fmt.Sprintf("%s\t($%02X,s),y", m, b), pos, false
	}
}

// xxx #nn
func op_immediateindex(m string) opcode {
	return func(pos uint32) (disassembled string, newpos uint32, done bool) {
		if !env.e.known || !env.x.known {
			addcomment(pos - 1, "(!) cannot disassemble index register opcode with immediate operand because size unknown")
			return fmt.Sprintf("%s\t???", m), pos, true
		} else {
			if env.e.value | env.x.value == 0 {
				w, pos := getword(pos)
				return fmt.Sprintf("%s\t#$%04X", m, w), pos, false
			} else {
				b, pos := getbyte(pos)
				return fmt.Sprintf("%s\t#$%02X", m, b), pos, false
			}
		}
	}
}
