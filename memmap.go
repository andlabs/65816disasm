// 25 may 2013
package main

import (
	"fmt"
)

type MemoryMap interface {
	Physical(logical uint32) (physical uint32, inROM bool)
	BankComment(bank byte) (bankComment string)
}

var memmap MemoryMap

type lowrom struct{}

func (lowrom) Physical(logical uint32) (physical uint32, inROM bool) {
	bank := (logical >> 16) & 0xFF		// keep it as uint32; we'll use it later
	base := uint16(logical & 0xFFFF)
	if (bank & 0x7F) <= 0x5F {
		if (base & 0x8000) == 0x8000 {		// ROM
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
		ROMstart := uint32(bank) * 32768			// banks are 32KB each
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

var LowROM lowrom
