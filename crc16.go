package dexcom

//go:generate ../crcgen/crcgen -size 16 -poly 0x1021

// Compute CRC-16 using CCITT polynomial.
func crc16(msg []byte) []byte {
	res := uint16(0)
	for _, b := range msg {
		res = res<<8 ^ crc16Table[byte(res>>8)^b]
	}
	return MarshalUint16(res)
}
