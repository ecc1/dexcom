package dexcom

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

// XMLData maps attribute names to values.
// The Dexcom CGM receiver represents its system data as single XML nodes
// with multiple attributes, so a tree structure is not required.
type XMLData map[string]string

// UnmarshalXML is called by xml.Unmarshal.
func (ptr *XMLData) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	m := *ptr
	for _, attr := range start.Attr {
		m[attr.Name.Local] = attr.Value
	}
	return d.Skip()
}

// UnmashalXMLData unmarshals a byte array into an XMLData map.
func UnmarshalXMLData(v []byte) (XMLData, error) {
	m := XMLData(make(map[string]string))
	err := xml.Unmarshal(v, &m)
	return m, err
}

// ReadFirmwareHeader gets the firmware header from the Dexcom CGM receiver
// and returns it as XMLData.
func (dev Device) ReadFirmwareHeader() (XMLData, error) {
	p, err := dev.Cmd(READ_FIRMWARE_HEADER)
	if err != nil {
		return nil, err
	}
	return UnmarshalXMLData(p)
}

// An XMLRecord contains timestamped XML data.
type XMLRecord struct {
	Timestamp Timestamp
	XML       XMLData
}

// UnmarshalXMLRecord unmarshals a byte array into an XMLRecord.
func UnmarshalXMLRecord(v []byte) (XMLRecord, error) {
	p := v[8:]
	i := bytes.IndexByte(p, 0x00)
	if i != -1 {
		p = p[:i]
	}
	xml, err := UnmarshalXMLData(p)
	return XMLRecord{
		Timestamp: UnmarshalTimestamp(v[0:8]),
		XML:       xml,
	}, err
}

// ReadXMLRecord gets the given XML record type from the Dexcom CGM receiver.
func (dev Device) ReadXMLRecord(recordType RecordType) (XMLRecord, error) {
	var xml XMLRecord
	proc := func(v []byte, context *RecordContext) error {
		// There should only be a single page, containing one record.
		if context.Index != 0 || context.PageNumber != context.StartPage || context.StartPage != context.EndPage {
			return fmt.Errorf("unexpected record context %v", *context)
		}
		var err error
		xml, err = UnmarshalXMLRecord(v)
		return err
	}
	err := dev.ReadRecords(recordType, proc)
	return xml, err
}
