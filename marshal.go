package dexcom

import (
	"math"
)

// MarshalUint16 marshals a uint16 value into 2 bytes in little-endian order.
func MarshalUint16(n uint16) []byte {
	return []byte{byte(n & 0xFF), byte(n >> 8)}
}

// MarshalInt16 marshals an int16 value into 2 bytes in little-endian order.
func MarshalInt16(n int16) []byte {
	return MarshalUint16(uint16(n))
}

// MarshalUint32 marshals a uint32 value into 4 bytes in little-endian order.
func MarshalUint32(n uint32) []byte {
	return append(MarshalUint16(uint16(n&0xFFFF)), MarshalUint16(uint16(n>>16))...)
}

// MarshalInt32 marshals an int32 value into 4 bytes in little-endian order.
func MarshalInt32(n int32) []byte {
	return MarshalUint32(uint32(n))
}

// UnmarshalUint16 unmarshals 2 bytes in little-endian order into a uint16 value.
func UnmarshalUint16(v []byte) uint16 {
	return uint16(v[0]) | uint16(v[1])<<8
}

// UnmarshalInt16 unmarshals 2 bytes in little-endian order into an int16 value.
func UnmarshalInt16(v []byte) int16 {
	return int16(UnmarshalUint16(v))
}

// UnmarshalUint32 unmarshals 4 bytes in little-endian order into a uint32 value.
func UnmarshalUint32(v []byte) uint32 {
	return uint32(UnmarshalUint16(v[0:2])) | uint32(UnmarshalUint16(v[2:4]))<<16
}

// UnmarshalInt32 unmarshals 4 bytes in little-endian order into an int32 value.
func UnmarshalInt32(v []byte) int32 {
	return int32(UnmarshalUint32(v))
}

// UnmarshalUint64 unmarshals 8 bytes in little-endian order into a uint64 value.
func UnmarshalUint64(v []byte) uint64 {
	return uint64(UnmarshalUint32(v[0:4])) | uint64(UnmarshalUint32(v[4:8]))<<32
}

// UnmarshalFloat64 unmarshals 8 bytes in little-endian order into a float64 value.
func UnmarshalFloat64(v []byte) float64 {
	return math.Float64frombits(UnmarshalUint64(v))
}
