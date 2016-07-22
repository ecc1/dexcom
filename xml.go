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

func (r *XMLData) Unmarshal(v []byte) error {
	*r = make(map[string]string)
	return xml.Unmarshal(v, r)
}

// ReadFirmwareHeader gets the firmware header from the Dexcom CGM receiver
// and returns it as XMLData.
func (cgm *Cgm) ReadFirmwareHeader() XMLData {
	x := XMLData{}
	p := cgm.Cmd(READ_FIRMWARE_HEADER)
	if cgm.Error() != nil {
		return x
	}
	err := x.Unmarshal(p)
	cgm.SetError(err)
	return x
}

// An XMLRecord contains timestamped XML data.
type XMLRecord struct {
	Timestamp Timestamp
	XML       XMLData
}

func (r *XMLRecord) Unmarshal(v []byte) error {
	p := v[8:]
	i := bytes.IndexByte(p, 0x00)
	if i != -1 {
		p = p[:i]
	}
	r.Timestamp.Unmarshal(v[0:8])
	return r.XML.Unmarshal(p)
}

// ReadXMLRecord gets the given XML record type from the Dexcom CGM receiver.
func (cgm *Cgm) ReadXMLRecord(recordType RecordType) XMLRecord {
	x := XMLRecord{}
	proc := func(v []byte, context RecordContext) (bool, error) {
		// There should only be a single page, containing one record.
		if context.Index != 0 || context.PageNumber != context.StartPage || context.StartPage != context.EndPage {
			return true, fmt.Errorf("unexpected record context %+v", context)
		}
		return false, x.Unmarshal(v)
	}
	cgm.ReadRecords(recordType, proc)
	return x
}
