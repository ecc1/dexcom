package dexcom

import (
	"bytes"
	"testing"
)

func TestCrc16(t *testing.T) {
	cases := []struct{ msg, sum []byte }{
		{[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, []byte{0x78, 0x23}},
		{[]byte("0123456789"), []byte{0x58, 0x9C}},
		{[]byte{0xFF, 0xFF, 0xFF, 0xFF}, []byte{0xCF, 0x99}},
		{[]byte{0x01, 0x07, 0x00, 0x10, 0x04}, []byte{0x8B, 0xB8}},
	}
	for _, c := range cases {
		sum := crc16(c.msg)
		if !bytes.Equal(sum, c.sum) {
			t.Errorf("crc16(%X) == %X, want %X", c.msg, sum, c.sum)
		}
	}
}
