package dexcom

import (
	"math"
)

// Marshaling and unmarshaling of ints and floats in little-endian order.

func marshalUint16(n uint16) []byte {
	return []byte{byte(n & 0xFF), byte(n >> 8)}
}

// nolint
func marshalInt16(n int16) []byte {
	return marshalUint16(uint16(n))
}

func marshalUint32(n uint32) []byte {
	return append(marshalUint16(uint16(n&0xFFFF)), marshalUint16(uint16(n>>16))...)
}

func marshalInt32(n int32) []byte {
	return marshalUint32(uint32(n))
}

func unmarshalUint16(v []byte) uint16 {
	return uint16(v[0]) | uint16(v[1])<<8
}

// nolint
func unmarshalInt16(v []byte) int16 {
	return int16(unmarshalUint16(v))
}

func unmarshalUint32(v []byte) uint32 {
	return uint32(unmarshalUint16(v[0:2])) | uint32(unmarshalUint16(v[2:4]))<<16
}

func unmarshalInt32(v []byte) int32 {
	return int32(unmarshalUint32(v))
}

func unmarshalUint64(v []byte) uint64 {
	return uint64(unmarshalUint32(v[0:4])) | uint64(unmarshalUint32(v[4:8]))<<32
}

func unmarshalFloat64(v []byte) float64 {
	return math.Float64frombits(unmarshalUint64(v))
}
