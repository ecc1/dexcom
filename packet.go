package dexcom

import (
	"bytes"
	"fmt"
)

const (
	startOfMessage = 1
	minPacket      = 6
	maxPacket      = 1590
)

func marshalPacket(cmd Command, data []byte) []byte {
	buf := bytes.Buffer{}
	buf.WriteByte(startOfMessage)
	length := uint16(4 + len(data) + 2) // header, data, CRC-16
	buf.Write(marshalUint16(length))
	buf.WriteByte(byte(cmd))
	buf.Write(data)
	body := buf.Bytes()
	buf.Write(marshalUint16(crc16(body))) // append CRC-16
	return buf.Bytes()
}

func (cgm *CGM) sendPacket(pkt []byte) {
	err := cgm.Send(pkt)
	cgm.SetError(err)
}

func (cgm *CGM) receivePacket() []byte {
	header := make([]byte, 4)
	err := cgm.Receive(header)
	if err != nil {
		cgm.SetError(err)
		return nil
	}
	if header[0] != startOfMessage {
		cgm.SetError(fmt.Errorf("unexpected message header % X", header))
		return nil
	}
	rc := Command(header[3])
	if rc != Ack {
		cgm.SetError(fmt.Errorf("unexpected response code %v in header % X", rc, header))
		return nil
	}
	length := unmarshalUint16(header[1:3])
	if length < minPacket || length > maxPacket {
		cgm.SetError(fmt.Errorf("invalid packet length %d in header % X", length, header))
		return nil
	}
	n := length - minPacket
	data := make([]byte, n+2)
	err = cgm.Receive(data)
	if err != nil {
		cgm.SetError(err)
		return nil
	}
	crc := unmarshalUint16(data[n:])
	data = data[:n]
	body := append(header, data...)
	calc := crc16(body)
	if crc != calc {
		cgm.SetError(CRCError{
			Kind:     "packet",
			Received: crc,
			Computed: calc,
			PageType: InvalidPage,
			Data:     body,
		})
	}
	return data
}

// Cmd creates a Dexcom packet with the given command and parameters,
// sends it to the device, and returns the response.
func (cgm *CGM) Cmd(cmd Command, params ...byte) []byte {
	if cgm.Error() != nil {
		return nil
	}
	pkt := marshalPacket(cmd, params)
	cgm.sendPacket(pkt)
	if cgm.Error() != nil {
		return nil
	}
	return cgm.receivePacket()
}
