package dexcom

//go:generate ../crcgen/crcgen

func crc16(msg []byte) []byte {
	res := uint16(0)
	for _, b := range msg {
		res = res<<8 ^ crc16Table[byte(res>>8)^b]
	}
	return MarshalUint16(res)
}
