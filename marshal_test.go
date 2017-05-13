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
		rep := marshalUint16(c.val)
		if !bytes.Equal(rep, c.rep) {
			t.Errorf("marshalUint16(%04X) == % X, want % X", c.val, rep, c.rep)
		}
		val := unmarshalUint16(c.rep)
		if val != c.val {
			t.Errorf("unmarshalUint16(% X) == %04X, want %04X", c.rep, val, c.val)
		}
	}
}

func TestMarshalInt16(t *testing.T) {
	cases := []struct {
		val int16
		rep []byte
	}{
		{0x1234, []byte{0x34, 0x12}},
		{0, []byte{0, 0}},
		{256, []byte{0, 1}},
		{-1, []byte{0xFF, 0xFF}},
		{-256, []byte{0x00, 0xFF}},
		{math.MaxInt16, []byte{0xFF, 0x7F}},
		{math.MinInt16, []byte{0x00, 0x80}},
	}
	for _, c := range cases {
		rep := marshalInt16(c.val)
		if !bytes.Equal(rep, c.rep) {
			t.Errorf("marshalInt16(%d) == % X, want % X", c.val, rep, c.rep)
		}
		val := unmarshalInt16(c.rep)
		if val != c.val {
			t.Errorf("unmarshalInt16(% X) == %d, want %d", c.rep, val, c.val)
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
		rep := marshalUint32(c.val)
		if !bytes.Equal(rep, c.rep) {
			t.Errorf("marshalUint32(%08X) == % X, want % X", c.val, rep, c.rep)
		}
		val := unmarshalUint32(c.rep)
		if val != c.val {
			t.Errorf("unmarshalUint32(% X) == %08X, want %08X", c.rep, val, c.val)
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
		rep := marshalInt32(c.val)
		if !bytes.Equal(rep, c.rep) {
			t.Errorf("marshalInt32(%d) == % X, want % X", c.val, rep, c.rep)
		}
		val := unmarshalInt32(c.rep)
		if val != c.val {
			t.Errorf("unmarshalInt32(% X) == %d, want %d", c.rep, val, c.val)
		}
	}
}

func TestMarshalUint64(t *testing.T) {
	cases := []struct {
		val uint64
		rep []byte
	}{
		{0x0123456789ABCDEF, []byte{0xEF, 0xCD, 0xAB, 0x89, 0x67, 0x45, 0x23, 0x01}},
		{0, []byte{0, 0, 0, 0, 0, 0, 0, 0}},
		{math.MaxUint64, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
	}
	for _, c := range cases {
		val := unmarshalUint64(c.rep)
		if val != c.val {
			t.Errorf("unmarshalUint64(% X) == %016X, want %016X", c.rep, val, c.val)
		}
	}
}

func TestUnmarshalFloat64(t *testing.T) {
	cases := []struct {
		rep []byte
		val float64
	}{
		{[]byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, math.SmallestNonzeroFloat64},
		{[]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xEF, 0x7F}, math.MaxFloat64},
		{[]byte{0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0xD5, 0x3F}, 0.3333333333333333},
		{[]byte{0x18, 0x2D, 0x44, 0x54, 0xFB, 0x21, 0x09, 0x40}, 3.141592653589793},
	}
	for _, c := range cases {
		val := unmarshalFloat64(c.rep)
		if val != c.val {
			t.Errorf("unmarshalFloat64(% X) == %v, want %v", c.rep, val, c.val)
		}
	}
}
