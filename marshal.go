package dexcom

func MarshalUint16(n uint16) []byte {
	return []byte{byte(n), byte(n >> 8)}
}

func MarshalUint32(n uint32) []byte {
	return append(MarshalUint16(uint16(n)), MarshalUint16(uint16(n>>16))...)
}

func MarshalInt32(n int32) []byte {
	return MarshalUint32(uint32(n))
}

func UnmarshalUint16(v []byte) uint16 {
	return uint16(v[0]) + uint16(v[1])<<8
}

func UnmarshalUint32(v []byte) uint32 {
	return uint32(UnmarshalUint16(v[0:2])) + uint32(UnmarshalUint16(v[2:4]))<<16
}

func UnmarshalInt32(v []byte) int32 {
	return int32(UnmarshalUint32(v))
}
