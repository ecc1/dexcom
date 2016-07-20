package dexcom

import (
	"testing"
)

func TestCrc16(t *testing.T) {
	cases := []struct {
		msg []byte
		sum uint16
	}{
		{[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, 0x2378},
		{[]byte("0123456789"), 0x9C58},
		{[]byte{0xFF, 0xFF, 0xFF, 0xFF}, 0x99CF},
		{[]byte{0x01, 0x07, 0x00, 0x10, 0x04}, 0xB88B},
	}
	for _, c := range cases {
		sum := crc16(c.msg)
		if sum != c.sum {
			t.Errorf("crc16(% X) == %X, want %X", c.msg, sum, c.sum)
		}
	}
}
