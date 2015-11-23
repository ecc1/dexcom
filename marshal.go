package dexcom

// MarshalUint16 marshals a uint16 value into 2 bytes in little-endian order.
func MarshalUint16(n uint16) []byte {
	return []byte{byte(n), byte(n >> 8)}
}

// MarshalUint32 marshals a uint32 value into 4 bytes in little-endian order.
func MarshalUint32(n uint32) []byte {
	return append(MarshalUint16(uint16(n)), MarshalUint16(uint16(n>>16))...)
}

// MarshalInt32 marshals a int32 value into 4 bytes in little-endian order.
func MarshalInt32(n int32) []byte {
	return MarshalUint32(uint32(n))
}

// UnmarshalUint16 unmarshals 2 bytes in little-endian order into a uint16 value.
func UnmarshalUint16(v []byte) uint16 {
	return uint16(v[0]) + uint16(v[1])<<8
}

// UnmarshalUint32 unmarshals 4 bytes in little-endian order into a uint32 value.
func UnmarshalUint32(v []byte) uint32 {
	return uint32(UnmarshalUint16(v[0:2])) + uint32(UnmarshalUint16(v[2:4]))<<16
}

// UnmarshalInt32 unmarshals 4 bytes in little-endian order into a int32 value.
func UnmarshalInt32(v []byte) int32 {
	return int32(UnmarshalUint32(v))
}
