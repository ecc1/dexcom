package main

import (
	"fmt"
)

// Generate lookup table for CRC-16 calculation
func gen_crc16(poly uint16) {
	fmt.Printf("// Lookup table for CRC-16 calculation with polynomial 0x%04X\n", poly)
	fmt.Printf("var crc16Table = []uint16{\n")
	for i := 0; i < 256; i++ {
		res := uint16(0)
		b := uint16(i << 8)
		for n := 0; n < 8; n++ {
			c := (res ^ b) & (1 << 15)
			res <<= 1
			b <<= 1
			if c != 0 {
				res ^= poly
			}
		}
		if i%8 == 0 {
			fmt.Printf("\t")
		} else {
			fmt.Printf(" ")
		}
		fmt.Printf("0x%04X,", res)
		if (i+1)%8 == 0 {
			fmt.Printf("\n")
		}
	}
	fmt.Printf("}\n")
}

func main() {
	// CCITT polynomial
	gen_crc16(0x1021)
}
