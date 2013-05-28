// 26 may 2013
package main

var opcodes [256]xyz

func init() {
	opcodes = [256]xyz{
	// adc: add with carry
	0x69:	op_immediate("adc")		// adc #nn
	0x6D:	op_absolute("adc")			// adc hhll
	0x6F:	op_absolutelong("adc")		// adc hhllmm
	0x65:	op_direct("adc")			// adc nn
	0x72:	op_indirect("adc")			// adc (nn)
	0x67:	op_indirectlong("adc")		// adc [nn]
	0x7D:	op_absolutex("adc")			// adc hhll,x
	0x7F:	op_absolutelongx("adc")		// adc hhllmm,x
	0x79:	op_absolutey("adc")			// adc hhll,y
	0x75:	op_directx("adc")			// adc nn,x
	0x61:	op_indirectx("adc")			// adc (nn,x)
	0x71:	op_indirecty("adc")			// adc (nn),y
	0x77:	op_indirectlongy("adc")		// adc [nn],y
	0x63:	op_stack("adc")			// adc nn,s
	0x73:	op_indirectstack("adc")		// adc (nn,s),y

	// and: bitwise and
	0x29:	op_immediate("and")		// and #nn
	0x2D:	op_absolute("and")			// and hhll
	0x2F:	op_absolutelong("and")		// and hhllmm
	0x25:	op_direct("and")			// and nn
	0x32:	op_indirect("and")			// and (nn)
	0x27:	op_indirectlong("and")		// and [nn]
	0x3D:	op_absolutex("and")			// and hhll,x
	0x3F:	op_absolutelongx("and")		// and hhllmm,x
	0x39:	op_absolutey("and")			// and hhll,y
	0x35:	op_directx("and")			// and nn,x
	0x21:	op_indirectx("and")			// and (nn,x)
	0x31:	op_indirecty("and")			// and (nn),y
	0x37:	op_indirectlongy("and")		// and [nn],y
	0x23:	op_stack("and")			// and nn,s
	0x33:	op_indirectstack("and")		// and (nn,s),y

	// asl: arithmetic shift left
	0x0A:	op_accumulator("asl")		// asl a
	0x0E:	op_absolute("asl")			// asl hhll
	0x06:	op_direct("asl")			// asl nn
	0x1E:	op_absolutex("asl")			// asl hhll,x
	0x16:	op_directx("asl")			// asl nn,x

	// bcc/blt: branch on carry clear/less than
	// I will use blt because I can never remember which (cc/cs) is which (lt/ge)
//	0x90:	op_branch("bcc")			// bcc addr
	0x90:	op_branch("blt")			// blt addr

	// bcs/bge: branch on carry set/greater than or equal to
	// I will use bge because I can never remember which (cc/cs) is which (lt/ge)
//	0xB0:	op_branch("bcs")			// bcs addr
	0xB0:	op_branch("bge")			// bge addr

	// beq: branch on equal
	0xF0:	op_branch("beq")			// beq addr

	// bit: test if a & operand == 0
	0x89:	op_immediate("bit")			// bit #nn
	0x2C:	op_absolute("bit")			// bit hhll
	0x24:	op_direct("bit")				// bit nn
	0x3C:	op_absolutex("bit")			// bit hhll,x
	0x34:	op_directx("bit")			// bit nn,x

	// bmi: branch on minus
	0x30:	op_branch("bmi")			// bmi addr

	// bne: branch on not equal
	0xD0:	op_branch("bne")			// bne addr

	// bpl: branch on plus
	0x10:	op_branch("bpl")			// bpl addr

	// bra: branch
	0x80:	op_branch("bra")			// bra addr

	// brk: trigger software interrupt
	// TODO handle signature byte?
	0x00:	op_noarguments("brk")		// brk

	// brl: branch to 16-bit offset
	0x82:	op_brl					// brl addr

	// bvc: branch on overflow clear
	0x50:	op_branch("bvc")			// bvc addr

	// bvs: branch on overflow set
	0x70:	op_branch("bvs")			// bvs addr

	// clc: clear carry flag
	0x18:	op_noarguments("clc")		// clc

	// cld: clear decimal flag
	0xD8:	op_noarguments("cld")		// cld

	// cli: enable interrupts
	0x58:	op_noarguments("cli")		// cli

	// clv: clear overflow flag
	0xB8:	op_noarguments("clv")		// clv

	// cmp: compare to a
	0xC9:	op_immediate("cmp")		// cmp #nn
	0xCD:	op_absolute("cmp")			// cmp hhll
	0xCF:	op_absolutelong("cmp")		// cmp hhllmm
	0xC5:	op_direct("cmp")			// cmp nn
	0xD2:	op_indirect("cmp")			// cmp (nn)
	0xC7:	op_indirectlong("cmp")		// cmp [nn]
	0xDD:	op_absolutex("cmp")		// cmp hhll,x
	0xDF:	op_absolutelongx("cmp")		// cmp hhllmm,x
	0xD9:	op_absolutey("cmp")		// cmp hhll,y
	0xD5:	op_directx("cmp")			// cmp nn,x
	0xC1:	op_indirectx("cmp")			// cmp (nn,x)
	0xD1:	op_indirecty("cmp")			// cmp (nn),y
	0xD7:	op_indirectlongy("cmp")		// cmp [nn],y
	0xC3:	op_stack("cmp")			// cmp nn,s
	0xD3:	op_indirectstack("cmp")		// cmp (nn,s),y

	// cop: call coprocessor
	0x02:	op_immediate("cop")		// cop #nn

	// cpx: compare to x
	0xE0:	op_immediate("cpx")		// cpx #nn
	0xEC:	op_absolute("cpx")			// cpx hhll
	0xE4:	op_direct("cpx")			// cpx nn

	// cpy: compare to y
	0xC0:	op_immediate("cpy")		// cpy #nn
	0xCC:	op_absolute("cpy")			// cpy hhll
	0xC4:	op_direct("cpy")			// cpy nn

	// dec: decrement
	0x3A:	op_accumulator("dec")		// dec a
	0xCE:	op_absolute("dec")			// dec hhll
	0xC6:	op_direct("dec")			// dec nn
	0xDE:	op_absolutex("dec")			// dec hhll,x
	0xD6:	op_directx("dec")			// dec nn,x

	// dex: decrement x
	0xCA:	op_noarguments("dex")		// dex

	// dey: decrement y
	0x88:	op_noarguments("dey")		// dey

	// eor: bitwise xor
	0x49:	op_immediate("eor")		// eor #nn
	0x4D:	op_absolute("eor")			// eor hhll
	0x4F:	op_absolutelong("eor")		// eor hhllmm
	0x45:	op_direct("eor")			// eor nn
	0x52:	op_indirect("eor")			// eor (nn)
	0x47:	op_indirectlong("eor")		// eor [nn]
	0x5D:	op_absolutex("eor")			// eor hhll,x
	0x5F:	op_absolutelongx("eor")		// eor hhllmm,x
	0x59:	op_absolutey("eor")			// eor hhll,y
	0x55:	op_directx("eor")			// eor nn,x
	0x41:	op_indirectx("eor")			// eor (nn,x)
	0x51:	op_indirecty("eor")			// eor (nn),y
	0x57:	op_indirectlongy("eor")		// eor [nn],y
	0x43:	op_stack("eor")			// eor nn,s
	0x53:	op_indirectstack("eor")		// eor (nn,s),y
