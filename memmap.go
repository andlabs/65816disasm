// 25 may 2013
package main

import (
	"fmt"
)

type MemoryMap interface {
	Physical(logical uint32) (physical uint32, inROM bool)
	BankComment(bank byte) (bankComment string)
	BankSize() uint32
}

var memmap MemoryMap

// TODO LowROM or LoROM?
type lowrom struct{}

func (lowrom) Physical(logical uint32) (physical uint32, inROM bool) {
	bank := (logical >> 16) & 0xFF		// keep it as uint32; we'll use it later
	base := uint16(logical & 0xFFFF)
	if (bank & 0x7F) <= 0x5F {
		if (base & 0x8000) == 0x8000 {		// ROM
			bank &= 0x7F			// handle mirrors properly
			bank *= 32768			// banks are 32KB each
			base &= 0x7FFF		// take the offset into the bank
			bank |= uint32(base)		// add it in
			return bank, true		// and there's our ROM address
		}
	}
	// TODO convert mirrors to their canonical form?
	return logical, false		// otherwise take the logical address as it is
}

func (lowrom) BankComment(bank byte) (bankComment string) {
	if (bank & 0x7F) <= 0x5F {
		ROMstart := uint32(bank & 0x7F) * 32768		// banks are 32KB each
		return fmt.Sprintf("bank $%02X -> ROM $%06X", bank, ROMstart)
	}
	if (bank >= 0x70) && (bank <= 0x77) {
		SRAMstart := uint32(bank) * 32768			// banks are 32KB each
		return fmt.Sprintf("bank $%02X -> SRAM $%06X", bank, SRAMstart)
	}
	if (bank == 0x7E) || (bank == 0x7F) {
		return fmt.Sprintf("bank $%02X -> RAM", bank)
	}
	return fmt.Sprintf("bank $%02X -> reserved", bank)
}

func (lowrom) BankSize() uint32 {
	return 0x8000
}

// TODO HighROM or HiROM?
type highrom struct{}

func (highrom) Physical(logical uint32) (physical uint32, inROM bool) {
	bank := (logical >> 16) & 0xFF		// keep it as uint32; we'll use it later
	base := uint16(logical & 0xFFFF)
	switch {
	case bank <= 0x3F:
		if (base & 0x8000) != 0x8000 {			// not in ROM
			break
		}
		bank |= 0xC0						// mirrored
		fallthrough
	case bank >= 0xC0:						// fixed to ROM
		bank &^= 0xC0						// get raw bank #
		fallthrough
	case (bank >= 0x40) && (bank <= 0x5F):		// fixed to ROM after above
		bank *= 65536						// banks are now 64KB
		bank |= uint32(base)					// add the base in
		return bank, true					// and there's our ROM address
	}
	// TODO convert mirrors to their canonical form?
	return logical, false		// otherwise take the logical address as it is
}

func (highrom) BankComment(bank byte) (bankComment string) {
	switch {
	case bank <= 0x2F:
		ROMstart := uint32(bank) * 32768			// banks are 32KB each here
		return fmt.Sprintf("bank $%02X -> ROM $%06X", bank, ROMstart)
	case (bank >= 0x30) && (bank <= 0x3F):
		m21SRAMstart := uint32(bank) * 0x2000		// SRAM banks are 8KB each here
		ROMstart := uint32(bank) * 32768			// ROM banks are 32KB each here
		return fmt.Sprintf("bank $%02X -> ROM $%06X; mode 21 SRAM $%06X",
			bank, ROMstart, m21SRAMstart)
	case bank >= 0xC0:
		bank &^= 0xC0
		ROMstart := uint32(bank) * 65536			// banks are 64KB each here
		return fmt.Sprintf("bank $%02X -> ROM $%06X", bank, ROMstart)
	case (bank >= 0x40) && (bank <= 0x5F):
		SRAMstart := uint32(bank) * 65536			// banks are 64KB each here
		return fmt.Sprintf("bank $%02X -> SRAM $%06X", bank, SRAMstart)
	case (bank >= 0x70) && (bank <= 0x77):
		m20SRAMstart := uint32(bank) * 32768		// banks are 32KB each here
		return fmt.Sprintf("bank $%02X -> mode 20 SRAM $%06X", bank, m20SRAMstart)
	case (bank == 0x7E) || (bank == 0x7F):
		return fmt.Sprintf("bank $%02X -> RAM", bank)
	}
	return fmt.Sprintf("bank $%02X -> reserved", bank)
}

func (highrom) BankSize() uint32 {
	return 0x10000
}

var memmaps = map[string]MemoryMap{
	"lowrom":		lowrom{},
	"highrom":	highrom{},
}
