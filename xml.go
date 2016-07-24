package dexcom

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

// XmlInfo maps attribute names to values.
// The Dexcom CGM receiver represents its system data as single XML nodes
// with multiple attributes, so a tree structure is not required.
type XmlInfo map[string]string

// UnmarshalXML is called by xml.Unmarshal.
func (ptr *XmlInfo) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	m := *ptr
	for _, attr := range start.Attr {
		m[attr.Name.Local] = attr.Value
	}
	return d.Skip()
}

func unmarshalXmlBytes(v []byte) *XmlInfo {
	m := make(XmlInfo)
	err := xml.Unmarshal(v, &m)
	if err != nil {
		m["InvalidXML"] = string(v)
	}
	return &m
}

func unmarshalXmlInfo(r *Record, v []byte) {
	v = v[8:]
	i := bytes.IndexByte(v, 0x00)
	if i != -1 {
		v = v[:i]
	}
	r.Xml = unmarshalXmlBytes(v)
}

// ReadFirmwareHeader gets the firmware header from the Dexcom CGM receiver
// and returns it as XmlInfo.
func (cgm *Cgm) ReadFirmwareHeader() *XmlInfo {
	v := cgm.Cmd(READ_FIRMWARE_HEADER)
	if cgm.Error() != nil {
		return nil
	}
	return unmarshalXmlBytes(v)
}

// ReadXMLRecord gets the given XML record type from the Dexcom CGM receiver.
func (cgm *Cgm) ReadXMLRecord(pageType PageType) Record {
	x := Record{}
	seen := false
	proc := func(r Record) (bool, error) {
		// There should only be a single page, containing one record.
		if seen {
			return true, fmt.Errorf("unexpected XML record in %v page", pageType)
		}
		x = r
		seen = true
		return false, nil
	}
	cgm.ReadRecords(pageType, proc)
	return x
}
