package dexcom

import (
	"bytes"
	"math"
	"testing"
)

func TestMarshalUint16(t *testing.T) {
	cases := []struct {
		val uint16
		rep []byte
	}{
		{0x1234, []byte{0x34, 0x12}},
		{0, []byte{0, 0}},
		{math.MaxUint16, []byte{0xFF, 0xFF}},
	}
	for _, c := range cases {
		rep := MarshalUint16(c.val)
		if !bytes.Equal(rep, c.rep) {
			t.Errorf("MarshalUint16(%X) == %X, want %X", c.val, rep, c.rep)
		}
		val := UnmarshalUint16(c.rep)
		if val != c.val {
			t.Errorf("UnmarshalUint16(%X) == %X, want %X", c.rep, val, c.val)
		}
	}
}

func TestMarshalUint32(t *testing.T) {
	cases := []struct {
		val uint32
		rep []byte
	}{
		{0x12345678, []byte{0x78, 0x56, 0x34, 0x12}},
		{0, []byte{0, 0, 0, 0}},
		{math.MaxUint32, []byte{0xFF, 0xFF, 0xFF, 0xFF}},
	}
	for _, c := range cases {
		rep := MarshalUint32(c.val)
		if !bytes.Equal(rep, c.rep) {
			t.Errorf("MarshalUint32(%X) == %X, want %X", c.val, rep, c.rep)
		}
		val := UnmarshalUint32(c.rep)
		if val != c.val {
			t.Errorf("UnmarshalUint32(%X) == %X, want %X", c.rep, val, c.val)
		}
	}
}

func TestMarshalInt32(t *testing.T) {
	cases := []struct {
		val int32
		rep []byte
	}{
		{0x12345678, []byte{0x78, 0x56, 0x34, 0x12}},
		{0, []byte{0, 0, 0, 0}},
		{-1, []byte{0xFF, 0xFF, 0xFF, 0xFF}},
		{0x0000FFFF, []byte{0xFF, 0xFF, 0, 0}},
		{-0x10000, []byte{0, 0, 0xFF, 0xFF}},
		{math.MaxInt32, []byte{0xFF, 0xFF, 0xFF, 0x7F}},
		{math.MinInt32, []byte{0, 0, 0, 0x80}},
	}
	for _, c := range cases {
		rep := MarshalInt32(c.val)
		if !bytes.Equal(rep, c.rep) {
			t.Errorf("MarshalInt32(%X) == %X, want %X", c.val, rep, c.rep)
		}
		val := UnmarshalInt32(c.rep)
		if val != c.val {
			t.Errorf("UnmarshalInt32(%X) == %X, want %X", c.rep, val, c.val)
		}
	}
}
